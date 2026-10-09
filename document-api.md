# Dokumentasi REST API

Dokumentasi lengkap seluruh endpoint API yang tersedia, dikelompokkan berdasarkan fitur, lengkap dengan metode HTTP, header, format payload (request body), permission yang dibutuhkan, dan contoh respon.

---

## 🌐 Informasi Umum
- **Base URL**: `http://localhost:8080`
- **Format Pertukaran Data**: `application/json`
- **Header Autentikasi**:
  ```http
  Authorization: Bearer <TOKEN_JWT_ANDA>
  Content-Type: application/json
  ```

---

## 👥 Akun untuk Pengujian (Seeder)

Gunakan akun berikut untuk menguji perbedaan permission/hak akses:

| Email | Password | Role | Hak Akses |
| :--- | :--- | :--- | :--- |
| `admin@example.com` | `admin123` | **admin** | Memiliki semua permission (`read-product`, `create-product`, `update-product`, `delete-product`) |
| `staff@example.com` | `staff123` | **staff** | Memiliki permission `read-product`, `create-product`, `update-product` |
| `user@example.com` | `user123` | **user** | Hanya memiliki permission `read-product` |

---

## 📋 Daftar Kelompok Endpoint

```text
├── 1. System (Health Check)
│   └── Health Check Server
│
├── 2. Auth (Autentikasi & Akun)
│   ├── Auth Register
│   ├── Auth Login
│   └── Auth Profile (Get Current User)
│
└── 3. Product (Manajemen Produk)
    ├── Product Get (Daftar Semua Produk)
    ├── Product Get by ID (Detail Produk)
    ├── Product Create (Tambah Produk)
    ├── Product Update (Ubah Data Produk)
    └── Product Delete (Hapus Produk)
```

---

## 1. System

### • Health Check Server
Mengecek apakah server berjalan dengan normal.
- **Method**: `GET`
- **URL**: `/`
- **Autentikasi**: Publik (Tanpa Token)
- **Response `200 OK`**:
  ```json
  {
    "message": "API Go + Gin + GORM + MySQL (RBAC Active) berjalan normal!",
    "status": "healthy"
  }
  ```

---

## 2. Auth (Autentikasi & Akun)

### • Auth Register
Mendaftarkan akun pengguna baru. Akun baru secara default akan memiliki role `user`.
- **Method**: `POST`
- **URL**: `/api/auth/register`
- **Autentikasi**: Publik (Tanpa Token)
- **Request Body**:
  ```json
  {
    "name": "Budi Santoso",
    "email": "budi@example.com",
    "password": "password123"
  }
  ```
- **Response `201 Created`**:
  ```json
  {
    "message": "Registrasi berhasil",
    "data": {
      "id": 4,
      "name": "Budi Santoso",
      "email": "budi@example.com",
      "role": {
        "id": 3,
        "name": "user",
        "description": "Pengguna biasa"
      },
      "created_at": "2026-10-09T10:44:34.500+07:00"
    }
  }
  ```
- **Response `409 Conflict`** (Email sudah digunakan):
  ```json
  {
    "error": "Email sudah terdaftar"
  }
  ```

---

### • Auth Login
Melakukan login untuk mendapatkan JWT Token yang digunakan pada endpoint yang terproteksi.
- **Method**: `POST`
- **URL**: `/api/auth/login`
- **Autentikasi**: Publik (Tanpa Token)
- **Request Body**:
  ```json
  {
    "email": "admin@example.com",
    "password": "admin123"
  }
  ```
- **Response `200 OK`**:
  ```json
  {
    "message": "Login berhasil",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6ImFkbWluQGV4YW1wbGUuY29tIiwiZXhwIjoxNzg5MDUwNjc0LCJpYXQiOjE3ODg5NjQyNzR9...",
    "user": {
      "id": 1,
      "name": "Admin Toko",
      "email": "admin@example.com",
      "role": {
        "id": 1,
        "name": "admin",
        "description": "Akses penuh sistem",
        "permissions": [
          { "id": 1, "name": "read-product", "description": "Melihat daftar dan detail produk" },
          { "id": 2, "name": "create-product", "description": "Menambahkan produk baru" },
          { "id": 3, "name": "update-product", "description": "Memperbarui informasi produk" },
          { "id": 4, "name": "delete-product", "description": "Menghapus produk" }
        ]
      }
    }
  }
  ```
- **Response `401 Unauthorized`**:
  ```json
  {
    "error": "Email atau kata sandi salah"
  }
  ```

---

### • Auth Profile (Get Current User)
Melihat data pengguna dan hak akses permission yang sedang login.
- **Method**: `GET`
- **URL**: `/api/auth/me`
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Headers**:
  ```http
  Authorization: Bearer <TOKEN_JWT>
  ```
- **Response `200 OK`**:
  ```json
  {
    "data": {
      "id": 1,
      "name": "Admin Toko",
      "email": "admin@example.com",
      "role": {
        "id": 1,
        "name": "admin",
        "permissions": [
          { "id": 1, "name": "read-product" },
          { "id": 2, "name": "create-product" },
          { "id": 3, "name": "update-product" },
          { "id": 4, "name": "delete-product" }
        ]
      },
      "created_at": "2026-10-09T10:44:34.500+07:00"
    }
  }
  ```

---

## 3. Product (Manajemen Produk)

Semua endpoint produk dilindungi autentikasi dan pengecekan permission granular.

---

### • Product Get (Daftar Semua Produk)
Mengambil seluruh data produk. Mendukung filter pencarian berdasarkan nama produk.
- **Method**: `GET`
- **URL**: `/api/products`
- **Query Parameter (Opsional)**:
  - `search`: Kata kunci pencarian nama produk (Contoh: `/api/products?search=asus`)
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `read-product`
- **Role yang Diizinkan**: `admin`, `staff`, `user`
- **Headers**:
  ```http
  Authorization: Bearer <TOKEN_JWT>
  ```
- **Response `200 OK`**:
  ```json
  {
    "message": "Berhasil mengambil data produk",
    "data": [
      {
        "id": 1,
        "name": "Laptop Asus ROG Zephyrus",
        "description": "Laptop gaming bertenaga tinggi dengan layar 165Hz dan RTX 4070.",
        "price": 24999000,
        "stock": 15,
        "user_id": 1,
        "created_at": "2026-10-09T10:44:34.500+07:00",
        "updated_at": "2026-10-09T10:44:34.500+07:00"
      }
    ]
  }
  ```

---

### • Product Get by ID (Detail Produk)
Mengambil detail satu produk berdasarkan ID.
- **Method**: `GET`
- **URL**: `/api/products/:id` (Contoh: `/api/products/1`)
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `read-product`
- **Role yang Diizinkan**: `admin`, `staff`, `user`
- **Headers**:
  ```http
  Authorization: Bearer <TOKEN_JWT>
  ```
- **Response `200 OK`**:
  ```json
  {
    "message": "Detail produk ditemukan",
    "data": {
      "id": 1,
      "name": "Laptop Asus ROG Zephyrus",
      "description": "Laptop gaming bertenaga tinggi dengan layar 165Hz dan RTX 4070.",
      "price": 24999000,
      "stock": 15,
      "user_id": 1
    }
  }
  ```
- **Response `404 Not Found`**:
  ```json
  {
    "error": "Produk tidak ditemukan"
  }
  ```

---

### • Product Create (Tambah Produk)
Menambahkan produk baru ke dalam sistem.
- **Method**: `POST`
- **URL**: `/api/products`
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `create-product`
- **Role yang Diizinkan**: `admin`, `staff`
- **Headers**:
  ```http
  Authorization: Bearer <TOKEN_JWT>
  Content-Type: application/json
  ```
- **Request Body**:
  ```json
  {
    "name": "Mechanical Keyboard Custom",
    "description": "Keyboard mekanik 75% hot-swappable dengan gasket mount.",
    "price": 1850000,
    "stock": 12
  }
  ```
- **Response `201 Created`**:
  ```json
  {
    "message": "Produk berhasil ditambahkan",
    "data": {
      "id": 6,
      "name": "Mechanical Keyboard Custom",
      "description": "Keyboard mekanik 75% hot-swappable dengan gasket mount.",
      "price": 1850000,
      "stock": 12,
      "user_id": 1,
      "created_at": "2026-10-09T10:48:00+07:00",
      "updated_at": "2026-10-09T10:48:00+07:00"
    }
  }
  ```
- **Response `403 Forbidden`** (Jika role tidak punya permission `create-product`, misal `user` biasa):
  ```json
  {
    "error": "Akses ditolak: Anda tidak memiliki izin untuk tindakan ini",
    "permission": "create-product"
  }
  ```

---

### • Product Update (Ubah Data Produk)
Memperbarui informasi nama, deskripsi, harga, atau stok produk.
- **Method**: `PUT`
- **URL**: `/api/products/:id` (Contoh: `/api/products/1`)
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `update-product`
- **Role yang Diizinkan**: `admin`, `staff`
- **Headers**:
  ```http
  Authorization: Bearer <TOKEN_JWT>
  Content-Type: application/json
  ```
- **Request Body**:
  ```json
  {
    "name": "Laptop Asus ROG Zephyrus G16 (2026)",
    "description": "Versi upgrade dengan prosesor Intel Core Ultra 9 dan 32GB RAM.",
    "price": 28999000,
    "stock": 10
  }
  ```
- **Response `200 OK`**:
  ```json
  {
    "message": "Produk berhasil diperbarui",
    "data": {
      "id": 1,
      "name": "Laptop Asus ROG Zephyrus G16 (2026)",
      "description": "Versi upgrade dengan prosesor Intel Core Ultra 9 dan 32GB RAM.",
      "price": 28999000,
      "stock": 10,
      "user_id": 1,
      "updated_at": "2026-10-09T10:49:15+07:00"
    }
  }
  ```
- **Response `403 Forbidden`** (Jika role tidak punya permission `update-product`, misal `user` biasa):
  ```json
  {
    "error": "Akses ditolak: Anda tidak memiliki izin untuk tindakan ini",
    "permission": "update-product"
  }
  ```

---

### • Product Delete (Hapus Produk)
Menghapus produk dari sistem (menggunakan *soft delete* GORM).
- **Method**: `DELETE`
- **URL**: `/api/products/:id` (Contoh: `/api/products/1`)
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `delete-product`
- **Role yang Diizinkan**: **Hanya `admin`**
- **Headers**:
  ```http
  Authorization: Bearer <TOKEN_JWT>
  ```
- **Response `200 OK`**:
  ```json
  {
    "message": "Produk berhasil dihapus"
  }
  ```
- **Response `403 Forbidden`** (Jika dicoba oleh role `staff` atau `user` biasa):
  ```json
  {
    "error": "Akses ditolak: Anda tidak memiliki izin untuk tindakan ini",
    "permission": "delete-product"
  }
  ```
- **Response `404 Not Found`**:
  ```json
  {
    "error": "Produk tidak ditemukan"
  }
  ```

---

## ⚠️ Respon Kode HTTP Umum

| Status Code | Makna | Penyebab |
| :---: | :--- | :--- |
| **`200 OK`** | Berhasil | Request berhasil diproses. |
| **`201 Created`** | Berhasil Dibuat | Data baru (User atau Produk) berhasil dibuat. |
| **`400 Bad Request`** | Data Tidak Valid | Format JSON salah atau ada field wajib yang kosong / tidak memenuhi validasi. |
| **`401 Unauthorized`** | Tidak Terotentikasi | Token JWT tidak dikirim, format bukan `Bearer <token>`, atau token kedaluwarsa. |
| **`403 Forbidden`** | Akses Ditolak | Token valid, namun role pengguna tidak memiliki izin (permission) yang sesuai. |
| **`404 Not Found`** | Tidak Ditemukan | Data yang dicari berdasarkan ID tidak ada di database. |
| **`409 Conflict`** | Duplikasi Data | Email yang didaftarkan sudah ada di sistem. |
| **`500 Internal Error`** | Kesalahan Server | Terjadi error pada database atau server internal. |

