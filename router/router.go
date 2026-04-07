package router

import (
	"Toko-Online/config"
	"Toko-Online/handler"
	"Toko-Online/middleware"
	"Toko-Online/repository"
	"Toko-Online/service"
	"log/slog"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	redis_sessions "github.com/gin-contrib/sessions/redis"
	"github.com/utrack/gin-csrf"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, logger *slog.Logger) *gin.Engine {
	r := gin.Default()

	// 1. CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 2. Security Headers (XSS Protection)
	r.Use(secure.New(secure.Config{
		STSSeconds:           31536000,
		STSIncludeSubdomains: true,
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
	}))

	// 3. Session Store
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "redis" // Default to service name in docker
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}
	redisAddr := redisHost + ":" + redisPort

	store, err := redis_sessions.NewStore(10, "tcp", redisAddr, os.Getenv("REDIS_PASSWORD"), os.Getenv("JWT_SECRET"))
	if err != nil {
		logger.Error("Failed to initialize Redis session store, falling back to CookieStore", "error", err)
		fallbackStore := cookie.NewStore([]byte(os.Getenv("JWT_SECRET")))
		r.Use(sessions.Sessions("toko_session", fallbackStore))
	} else {
		logger.Info("Successfully initialized Redis session store")
		r.Use(sessions.Sessions("toko_session", store))
	}

	// 4. CSRF Middleware
	csrfMiddleware := csrf.Middleware(csrf.Options{
		Secret: os.Getenv("JWT_SECRET"),
		ErrorFunc: func(c *gin.Context) {
			c.JSON(400, gin.H{"error": "CSRF token mismatch"})
			c.Abort()
		},
	})

	// 0. Initialize External Clients
	redisClient := config.NewRedisClient(logger)
	cloudinaryClient := config.NewCloudinaryClient(logger)

	// 1. Initialize Repositories
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	userRepo := repository.NewUserRepository(db)
	settingRepo := repository.NewSettingRepository(db)

	// 2. Initialize Services
	categoryService := service.NewCategoryService(categoryRepo, logger, redisClient)
	productService := service.NewProductService(productRepo, categoryRepo, logger, redisClient)
	orderService := service.NewOrderService(db, orderRepo, productRepo, logger)
	userService := service.NewUserService(userRepo, logger, redisClient)
	settingService := service.NewSettingService(settingRepo, logger)
	uploadService := service.NewUploadService(cloudinaryClient, logger)

	// 3. Initialize Handlers
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)
	userHandler := handler.NewUserHandler(userService)
	settingHandler := handler.NewSettingHandler(settingService)
	uploadHandler := handler.NewUploadHandler(uploadService)

	// 4. Setup Routes
	api := r.Group("/api")
	{
		api.POST("/webhook/midtrans", orderHandler.Webhook)

		// ROUTES WITH CSRF
		apiWithCsrf := api.Group("/")
		apiWithCsrf.Use(csrfMiddleware)
		{
			// Endpoint to get CSRF token
			apiWithCsrf.GET("/csrf-token", func(c *gin.Context) {
				c.JSON(200, gin.H{"csrf_token": csrf.GetToken(c)})
			})

			// PUBLIC ROUTES
			apiWithCsrf.GET("/products", productHandler.FindAll)
			apiWithCsrf.POST("/checkout", orderHandler.Checkout)
			apiWithCsrf.GET("/products/:id", productHandler.FindByID)

			// HIDDEN ADMIN LOGIN ROUTE
			apiWithCsrf.POST("/admin/login", userHandler.Login)

			admin := apiWithCsrf.Group("/")
			admin.Use(middleware.AdminAuthMiddleware(redisClient))
			{
				// Session Management
				admin.POST("/admin/logout", userHandler.Logout)

				// Image Upload
				admin.POST("/upload", uploadHandler.Upload)

				// Products
				admin.POST("/products", productHandler.Create)
				admin.PATCH("/products/:id", productHandler.Update)
				admin.DELETE("/products/:id", productHandler.Delete)
				
				// Categories
				admin.GET("/categories", categoryHandler.FindAll)
				admin.POST("/categories", categoryHandler.Create)
				admin.PATCH("/categories/:id", categoryHandler.Update)
				admin.DELETE("/categories/:id", categoryHandler.Delete)

				// Settings
				admin.GET("/settings", settingHandler.Get)
				admin.PATCH("/settings", settingHandler.Update)

				// Profile
				admin.GET("/profile/:id", userHandler.GetProfile)
				admin.PATCH("/profile/:id", userHandler.UpdateProfile)
				
				// Orders (Dashboard)
				admin.GET("/orders", orderHandler.FindAll)
				admin.GET("/orders/:id", orderHandler.FindByID)
			}
		}
	}

	return r
}
