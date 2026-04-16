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
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, logger *slog.Logger) *gin.Engine {
	r := gin.Default()

	// 1. CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 2. Security Headers
	r.Use(secure.New(secure.Config{
		STSSeconds:           31536000,
		STSIncludeSubdomains: true,
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
	}))

	// 3. Session Store
	sessionOptions := sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	store := cookie.NewStore([]byte(os.Getenv("JWT_SECRET")))
	store.Options(sessionOptions)
	r.Use(sessions.Sessions("toko_session", store))

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
	categoryService := service.NewCategoryService(categoryRepo, productRepo, logger, redisClient)
	productService := service.NewProductService(productRepo, categoryRepo, logger, redisClient)
	orderService := service.NewOrderService(db, orderRepo, productRepo, logger, redisClient)
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

		// PUBLIC ROUTES
		api.GET("/products", productHandler.FindAll)
		api.POST("/checkout", orderHandler.Checkout)
		api.GET("/products/:id", productHandler.FindByID)
		api.GET("/categories", categoryHandler.FindAll)


		api.POST("/admin/login", userHandler.Login)

		admin := api.Group("/")
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
			admin.PATCH("/orders/:id/cancel", orderHandler.Cancel)
		}
	}

	return r
}
