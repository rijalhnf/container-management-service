# Container Management Service

REST API berbasis **Go + Gin + PostgreSQL** yang sudah dilengkapi dengan:

- ✅ **JWT Authentication** — register, login, dan validasi token
- ✅ **Auth Middleware** — semua route produk dilindungi, wajib pakai token
- ✅ **CORS Middleware** — mendukung akses dari frontend berbeda origin
- ✅ **CRUD Product** — contoh entity lengkap sebagai referensi
- ✅ **Docker Compose** — satu perintah untuk menjalankan API + database

---

## Struktur Folder

```
container-management-service/
├── docker-compose.yml           # Orkestrasi Docker (API + PostgreSQL)
├── postman_collection.json      # Koleksi Postman untuk testing
└── core-project/                # Source code utama
    ├── main.go
    ├── go.mod
    ├── Dockerfile
    ├── config/
    │   └── database.go          # Koneksi ke PostgreSQL via GORM
    ├── models/
    │   ├── product.go           # Entity Product (contoh)
    │   └── user.go              # Entity User (untuk auth)
    ├── controllers/
    │   ├── product_controller.go  # Handler CRUD Product
    │   └── auth_controller.go     # Handler Register / Login / Me
    ├── middleware/
    │   ├── auth.go              # JWT middleware (guard route)
    │   └── cors.go              # CORS middleware
    ├── utils/
    │   └── jwt.go               # Generate & validasi JWT token
    ├── services/                # (opsional) business logic terpisah
    └── routes/
        └── routes.go            # Definisi semua route
```

---

## Cara Menjalankan

> Pastikan **Docker** sudah berjalan.

```bash
cd container-management-service

# Pertama kali / setelah ada perubahan kode
docker compose down -v && docker compose up --build

# Selanjutnya (tanpa rebuild)
docker compose up
```

API berjalan di: `http://localhost:8888`

---

## Environment Variables

Didefinisikan di `docker-compose.yml`:

| Variable | Default | Keterangan |
|---|---|---|
| `APP_PORT` | `8888` | Port server |
| `DB_HOST` | `postgres` | Host database |
| `DB_PORT` | `5432` | Port database |
| `DB_USER` | `postgres` | Username database |
| `DB_PASSWORD` | `postgres` | Password database |
| `DB_NAME` | `crud_db` | Nama database |
| `JWT_SECRET` | *(lihat compose)* | Secret untuk signing JWT — **ganti di production!** |
| `CORS_ORIGIN` | `*` | Origin yang diizinkan. Ganti ke domain spesifik di production |

---

## Daftar Endpoint

### Public (tanpa token)

| Method | Endpoint | Keterangan |
|---|---|---|
| `GET` | `/health` | Health check |
| `POST` | `/api/auth/register` | Daftar akun baru |
| `POST` | `/api/auth/login` | Login, mendapatkan JWT token |

### Protected 🔒 (wajib `Authorization: Bearer <token>`)

| Method | Endpoint | Keterangan |
|---|---|---|
| `GET` | `/api/auth/me` | Info user dari token |
| `GET` | `/api/products` | Ambil semua produk |
| `GET` | `/api/products/:id` | Ambil produk by ID |
| `POST` | `/api/products` | Buat produk (single / batch array) |
| `PUT` | `/api/products/:id` | Update produk by ID |
| `DELETE` | `/api/products/:id` | Hapus produk by ID |

---

## Contoh Penggunaan (curl)

```bash
# 1. Register
curl -X POST http://localhost:8888/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}'

# 2. Login → salin token dari response
curl -X POST http://localhost:8888/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}'

# 3. Akses route yang dilindungi
curl http://localhost:8888/api/products \
  -H "Authorization: Bearer <token>"
```

Atau impor `postman_collection.json` ke Postman — token akan **otomatis tersimpan** setelah Login.

---

## Contoh Entity: Product

### Model (`models/product.go`)

```go
type Product struct {
    ID        uint      `gorm:"primaryKey;autoIncrement"`
    Name      string    `gorm:"type:varchar(255)"`
    Price     int       `gorm:"type:int"`
    Stock     int       `gorm:"type:int"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Controller (`controllers/product_controller.go`)

Berisi 5 fungsi handler: `GetProducts`, `GetProductByID`, `CreateProduct`, `UpdateProduct`, `DeleteProduct`.
Setiap fungsi menerima `*gin.Context` dan berinteraksi langsung ke database via `config.DB`.

### Route (`routes/routes.go`)

```go
api := router.Group("/api", middleware.AuthMiddleware())
{
    api.GET("products", controllers.GetProducts)
    api.GET("products/:id", controllers.GetProductByID)
    api.POST("products", controllers.CreateProduct)
    api.PUT("products/:id", controllers.UpdateProduct)
    api.DELETE("products/:id", controllers.DeleteProduct)
}
```

---

## Cara Menambahkan Entity Baru

Misal ingin menambahkan entity **Category**:

### 1. Buat Model — `models/category.go`

```go
package models

import "time"

type Category struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name      string    `gorm:"type:varchar(255);not null" json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 2. AutoMigrate di `main.go`

```go
config.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Category{})
```

### 3. Buat Controller — `controllers/category_controller.go`

```go
package controllers

import (
    "net/http"
    "go-gin-postgre-crud/config"
    "go-gin-postgre-crud/models"
    "github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
    var categories []models.Category
    config.DB.Find(&categories)
    c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
    var category models.Category
    if err := c.ShouldBindJSON(&category); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    config.DB.Create(&category)
    c.JSON(http.StatusCreated, category)
}

// ... tambahkan GetByID, Update, Delete sesuai kebutuhan
```

### 4. (Opsional) Buat Service — `services/category_service.go`

Gunakan layer service jika business logic-nya kompleks (validasi, kalkulasi, pemanggilan service lain, dsb).

```go
package services

import (
    "go-gin-postgre-crud/config"
    "go-gin-postgre-crud/models"
)

func GetAllCategories() ([]models.Category, error) {
    var categories []models.Category
    result := config.DB.Find(&categories)
    return categories, result.Error
}
```

### 5. Daftarkan Route di `routes/routes.go`

```go
// Di dalam blok protected (api group)
api.GET("categories", controllers.GetCategories)
api.POST("categories", controllers.CreateCategory)
```

### Ringkasan Checklist

| Langkah | File | Wajib? |
|---|---|---|
| Buat struct entity | `models/category.go` | ✅ |
| AutoMigrate | `main.go` | ✅ |
| Buat handler CRUD | `controllers/category_controller.go` | ✅ |
| Daftarkan route | `routes/routes.go` | ✅ |
| Buat service | `services/category_service.go` | ⬜ Opsional |
