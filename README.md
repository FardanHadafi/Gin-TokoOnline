# Toko-Online Backend 🛒

A production-ready e-commerce backend built with **Go (Gin)**, featuring high-performance caching, stateful sessions, and seamless payment integration.

## 🚀 Features

- **High Performance**: Optimized with **Redis** Cache-Aside pattern for products and categories.
- **Stateful Sessions**: Secure multi-device login management using Redis.
- **Image Management**: Professional cloud storage integration with **Cloudinary**.
- **Payment Integration**: Automated payment handling via **Midtrans SNAP** & Webhooks.
- **WhatsApp Integration**: Instant checkout confirmation link for customers.
- **Containerized**: Fully orchestrated with **Docker Compose** & **Nginx** reverse proxy.
- **Test-Driven**: Comprehensive unit test suite with 90% logic coverage target.

## 🛠️ Technology Stack

- **Languange**: Go 1.25+
- **Framework**: Gin Gonic
- **Database**: Supabase
- **ORM**: GORM
- **Cache & Session**: Redis
- **Proxy**: Nginx
- **Payment**: Midtrans
- **Image**: Cloudinary SDK v2
- **Testing**: Testify

## ⚙️ Getting Started

### 1. Environment Setup

Create a `.env` file in the root directory:

```env
DATABASE_URL=postgres://user:pass@db:5432/toko_online
REDIS_URL=redis:6379
JWT_SECRET=your_secret_key
ADMIN_WHATSAPP=628123456789
CLOUDINARY_URL=cloudinary://api_key:api_secret@cloud_name
MIDTRANS_SERVER_KEY=your_server_key
```

### 2. Run with Docker

```bash
docker-compose up --build
```

The API will be available at `http://localhost/api`.

### 3. Run Tests

```bash
go test -v ./test/... -coverpkg=Toko-Online/service
```

## 📂 Project Structure

- `/config`: Database and client initializations.
- `/dto`: Data Transfer Objects for API requests/responses.
- `/handler`: Controller layer (HTTP logic).
- `/middleware`: Auth and security filters.
- `/model`: Database schemas (GORM).
- `/repository`: Data access layer.
- `/service`: Business logic & Cache management.
- `/test`: Unit tests and Mocks.
- `/utils`: Common helper functions.

---
