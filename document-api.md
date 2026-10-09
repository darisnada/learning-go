# Dokumentasi REST API

Dokumentasi lengkap seluruh endpoint API yang tersedia, dikelompokkan berdasarkan fitur, lengkap dengan metode HTTP, header, format payload (request body), permission yang dibutuhkan, dan contoh respon. Semua entitas menggunakan format **UUID v4**.

---

## 🌐 Informasi Umum
- **Base URL**: `http://localhost:8080`
- **Format Pertukaran Data**: `application/json`
- **Tipe ID**: `UUID v4` (String 36 karakter, contoh: `8d78882e-42cb-4974-801b-cad0b8ffce38`)
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
      "id": "c1f7b8a2-3e4b-4a5d-b891-98721cba1024",
      "name": "Budi Santoso",
      "email": "budi@example.com",
      "role": {
        "id": "73f5c7c8-c0da-43e1-ba5f-9e252d22a537",
        "name": "user",
        "description": "Pengguna biasa"
      },
      "created_at": "2026-10-09T11:11:39+07:00"
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
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZThlZDU3NGYtOWVjYS00ODIxLTk1YTAtNjkwZDk4YzE4MmJjIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsImV4cCI6MTc4OTA1MDY3NH0...",
    "user": {
      "id": "e8ed574f-9eca-4821-95a0-690d98c182bc",
      "name": "Admin Toko",
      "email": "admin@example.com",
      "role": {
        "id": "4a95116a-9dad-4ebe-9aaf-835c11167f47",
        "name": "admin",
        "description": "Akses penuh sistem",
        "permissions": [
          { "id": "10134ee1-a0c1-48b2-b34c-31ec8fa181d2", "name": "read-product", "description": "Melihat daftar dan detail produk" },
          { "id": "e8196784-b678-4e58-869e-76c8dbb0208c", "name": "create-product", "description": "Menambahkan produk baru" },
          { "id": "5addc9aa-eba9-410d-a965-099a4a329fc4", "name": "update-product", "description": "Memperbarui informasi produk" },
          { "id": "d9d2515e-4a37-4982-ae17-41e7124238a4", "name": "delete-product", "description": "Menghapus produk" }
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
      "id": "e8ed574f-9eca-4821-95a0-690d98c182bc",
      "name": "Admin Toko",
      "email": "admin@example.com",
      "role": {
        "id": "4a95116a-9dad-4ebe-9aaf-835c11167f47",
        "name": "admin",
        "permissions": [
          { "id": "10134ee1-a0c1-48b2-b34c-31ec8fa181d2", "name": "read-product" },
          { "id": "e8196784-b678-4e58-869e-76c8dbb0208c", "name": "create-product" },
          { "id": "5addc9aa-eba9-410d-a965-099a4a329fc4", "name": "update-product" },
          { "id": "d9d2515e-4a37-4982-ae17-41e7124238a4", "name": "delete-product" }
        ]
      },
      "created_at": "2026-10-09T11:11:39+07:00"
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
        "id": "8d78882e-42cb-4974-801b-cad0b8ffce38",
        "name": "Laptop Asus ROG Zephyrus",
        "description": "Laptop gaming bertenaga tinggi dengan layar 165Hz dan RTX 4070.",
        "price": 24999000,
        "stock": 15,
        "user_id": "e8ed574f-9eca-4821-95a0-690d98c182bc",
        "created_at": "2026-10-09T11:11:39+07:00",
        "updated_at": "2026-10-09T11:11:39+07:00"
      }
    ]
  }
  ```

---

### • Product Get by ID (Detail Produk)
Mengambil detail satu produk berdasarkan UUID.
- **Method**: `GET`
- **URL**: `/api/products/:id` (Contoh: `/api/products/8d78882e-42cb-4974-801b-cad0b8ffce38`)
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
      "id": "8d78882e-42cb-4974-801b-cad0b8ffce38",
      "name": "Laptop Asus ROG Zephyrus",
      "description": "Laptop gaming bertenaga tinggi dengan layar 165Hz dan RTX 4070.",
      "price": 24999000,
      "stock": 15,
      "user_id": "e8ed574f-9eca-4821-95a0-690d98c182bc"
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
Menambahkan produk baru ke dalam sistem. ID produk akan otomatis di-generate berupa UUID v4.
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
      "id": "e45a1910-b77e-49b0-9db0-f8ca779eb672",
      "name": "Mechanical Keyboard Custom",
      "description": "Keyboard mekanik 75% hot-swappable dengan gasket mount.",
      "price": 1850000,
      "stock": 12,
      "user_id": "e8ed574f-9eca-4821-95a0-690d98c182bc",
      "created_at": "2026-10-09T11:15:00+07:00",
      "updated_at": "2026-10-09T11:15:00+07:00"
    }
  }
  ```
- **Response `403 Forbidden`** (Jika role tidak punya permission `create-product`):
  ```json
  {
    "error": "Akses ditolak: Anda tidak memiliki izin untuk tindakan ini",
    "permission": "create-product"
  }
  ```

---

### • Product Update (Ubah Data Produk)
Memperbarui informasi produk berdasarkan ID UUID.
- **Method**: `PUT`
- **URL**: `/api/products/:id` (Contoh: `/api/products/8d78882e-42cb-4974-801b-cad0b8ffce38`)
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
      "id": "8d78882e-42cb-4974-801b-cad0b8ffce38",
      "name": "Laptop Asus ROG Zephyrus G16 (2026)",
      "description": "Versi upgrade dengan prosesor Intel Core Ultra 9 dan 32GB RAM.",
      "price": 28999000,
      "stock": 10,
      "user_id": "e8ed574f-9eca-4821-95a0-690d98c182bc",
      "updated_at": "2026-10-09T11:16:00+07:00"
    }
  }
  ```
- **Response `403 Forbidden`**:
  ```json
  {
    "error": "Akses ditolak: Anda tidak memiliki izin untuk tindakan ini",
    "permission": "update-product"
  }
  ```

---

### • Product Delete (Hapus Produk)
Menghapus produk dari sistem berdasarkan ID UUID.
- **Method**: `DELETE`
- **URL**: `/api/products/:id` (Contoh: `/api/products/8d78882e-42cb-4974-801b-cad0b8ffce38`)
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
- **Response `403 Forbidden`**:
  ```json
  {
    "error": "Akses ditolak: Anda tidak memiliki izin untuk tindakan ini",
    "permission": "delete-product"
  }
  ```
