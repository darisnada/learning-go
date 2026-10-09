# REST API Go + Gin + GORM + MySQL (RBAC Role & Permission)

Aplikasi REST API lengkap dengan autentikasi JWT, sistem otorisasi **RBAC (Role-Based Access Control)** dengan granular permissions (`read-product`, `create-product`, `update-product`, `delete-product`), migrasi tabel otomatis, dan data seeder.

---

## 🛠️ Tech Stack & Library
- **Bahasa**: Go (Golang)
- **Framework Web**: [Gin Gonic](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/) + Driver MySQL
- **Autentikasi**: [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) & [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- **Otorisasi**: RBAC Middleware (`middlewares.RequirePermission`)
- **Env Loader**: [godotenv](https://github.com/joho/godotenv)

---

## 🔐 Matriks Role & Permissions

| Permission | Keterangan | Admin | Staff | User Biasa |
| :--- | :--- | :---: | :---: | :---: |
| `read-product` | Melihat daftar dan detail produk | ✅ | ✅ | ✅ |
| `create-product` | Menambahkan produk baru | ✅ | ✅ | ❌ |
| `update-product` | Memperbarui data produk | ✅ | ✅ | ❌ |
| `delete-product` | Menghapus produk | ✅ | ❌ | ❌ |

---

## 👥 Akun Bawaan (Hasil Seeder)

| Email | Password | Role | Hak Akses |
| :--- | :--- | :--- | :--- |
| `admin@example.com` | `admin123` | **admin** | Memiliki semua permissions (`read`, `create`, `update`, `delete`) |
| `staff@example.com` | `staff123` | **staff** | Bisa `read-product`, `create-product`, `update-product` (tidak bisa delete) |
| `user@example.com` | `user123` | **user** | Hanya bisa `read-product` (view only) |

---

## 📁 Struktur Folder
```text
c:\project\go-learning\
├── config/
│   └── database.go           # Koneksi MySQL & Auto Migration
├── controllers/
│   ├── auth_controller.go    # Handler Register, Login, Profile
│   └── product_controller.go # Handler CRUD Produk
├── middlewares/
│   ├── auth_middleware.go    # Verifikasi Bearer JWT
│   └── permission_middleware.go # Verifikasi granular Permission (RBAC)
├── models/
│   ├── rbac.go               # Model Role dan Permission
│   ├── user.go               # Model User & method HasPermission
│   └── product.go            # Model Product
├── seeders/
│   └── seeder.go             # Seeder Permissions, Roles, Users, Products
├── utils/
│   └── jwt.go                # Fungsi Generate & Validate JWT
├── .env                      # Konfigurasi database & port
├── .env.example              # Template variabel lingkungan
├── main.go                   # Routing, Middleware, dan Server entrypoint
└── go.mod
```

---

## 🚀 Cara Menjalankan

1. **Konfigurasi Database** di file [`.env`](file:///c:/project/go-learning/.env):
   ```ini
   DB_USER=root
   DB_PASSWORD=password_mariadb_anda
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_NAME=go_restapi
   JWT_SECRET=supersecretjwtkey12345
   PORT=8080
   ```

2. **Jalankan Migrasi & Seeder**:
   ```powershell
   # Reset database & isi data awal (Role, Permission, Akun Demo, Produk)
   go run main.go -fresh-seed

   # Atau cukup seeder saja jika tabel sudah ada:
   go run main.go -seed
   ```

3. **Jalankan Server API**:
   ```powershell
   go run main.go
   ```

---

## 📡 Daftar Endpoint API

### 1. Autentikasi (`/api/auth`)

#### A. Registrasi Pengguna
- **Method & URL**: `POST /api/auth/register`
- **Request Body**:
  ```json
  {
    "name": "Budi Santoso",
    "email": "budi@example.com",
    "password": "password123"
  }
  ```
  *(Default role pengguna baru adalah `user`)*.

#### B. Login
- **Method & URL**: `POST /api/auth/login`
- **Request Body**:
  ```json
  {
    "email": "admin@example.com",
    "password": "admin123"
  }
  ```
- **Response**: Mengembalikan token JWT (`token`) dan detail role serta permission.

#### C. Profil User (Perlu Login)
- **Method & URL**: `GET /api/auth/me`
- **Header**: `Authorization: Bearer <TOKEN_JWT>`

---

### 2. Produk (`/api/products`) — Dilindungi RBAC

Semua request wajib menyertakan header:
```text
Authorization: Bearer <TOKEN_JWT>
```

| Method | Endpoint | Required Permission | Siapa yang Boleh? |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/products` | `read-product` | Admin, Staff, User |
| `GET` | `/api/products/:id` | `read-product` | Admin, Staff, User |
| `POST` | `/api/products` | `create-product` | Admin, Staff |
| `PUT` | `/api/products/:id` | `update-product` | Admin, Staff |
| `DELETE` | `/api/products/:id` | `delete-product` | **Hanya Admin** |

#### Respon jika Role Tidak Memiliki Izin (Contoh User mencoba Delete):
**Status HTTP**: `403 Forbidden`
```json
{
  "error": "Akses ditolak: Anda tidak memiliki izin untuk tindakan ini",
  "permission": "delete-product"
}
```
