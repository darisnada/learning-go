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
| `admin@example.com` | `admin123` | **admin** | Memiliki semua permission (`read-product`, `create-product`, `update-product`, `delete-product`, `read-role`, `create-role`, `update-role`, `delete-role`) |
| `staff@example.com` | `staff123` | **staff** | Memiliki permission produk (`read-product`, `create-product`, `update-product`) |
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
├── 3. Product (Manajemen Produk)
│   ├── Product Get (Daftar Semua Produk)
│   ├── Product Get by ID (Detail Produk)
│   ├── Product Create (Tambah Produk)
│   ├── Product Update (Ubah Data Produk)
│   └── Product Delete (Hapus Produk)
│
├── 4. Role & Permission (Manajemen Hak Akses RBAC)
│   ├── Permission Get (Daftar Semua Permission di Sistem)
│   ├── Role Get (Daftar Semua Role & Permission-nya)
│   ├── Role Get by ID (Detail Role & Permission-nya)
│   ├── Role Create (Tambah Role Baru)
│   ├── Role Update (Ubah Nama / Deskripsi Role)
│   ├── Role Update Permissions (Update / Assign Permissions ke Role)
│   └── Role Delete (Hapus Role)
│
└── 5. User (Manajemen Pengguna)
    ├── User Get (Daftar Pengguna & Filter Search)
    ├── User Get by ID (Detail Pengguna)
    ├── User Create (Buat Pengguna Baru oleh Admin)
    ├── User Update (Ubah Profil Pengguna)
    ├── User Update Role (Ubah Role Pengguna)
    └── User Delete (Hapus Pengguna)
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
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "c7c7d18e-c1ac-45e0-a69f-a2861e454404",
      "name": "Admin Toko",
      "email": "admin@example.com",
      "role": {
        "id": "567e58bf-09da-4e99-9200-d59799ec7554",
        "name": "admin",
        "description": "Akses penuh sistem"
      }
    }
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
      "id": "c7c7d18e-c1ac-45e0-a69f-a2861e454404",
      "name": "Admin Toko",
      "email": "admin@example.com",
      "role": {
        "id": "567e58bf-09da-4e99-9200-d59799ec7554",
        "name": "admin",
        "permissions": [
          { "id": "3830afbd-65a8-4539-bec9-d3c4da141eda", "name": "read-product" },
          { "id": "3d1526ed-b52a-4ee3-8f75-2766ef2e71c2", "name": "create-product" }
        ]
      }
    }
  }
  ```

---

## 3. Product (Manajemen Produk)

Semua endpoint produk dilindungi autentikasi dan pengecekan permission granular.

---

### • Product Get (Daftar Semua Produk)
- **Method**: `GET`
- **URL**: `/api/products`
- **Query Parameter (Opsional)**: `?search=asus`
- **Permission**: `read-product` (Admin, Staff, User)

---

### • Product Get by ID (Detail Produk)
- **Method**: `GET`
- **URL**: `/api/products/:id`
- **Permission**: `read-product` (Admin, Staff, User)

---

### • Product Create (Tambah Produk)
- **Method**: `POST`
- **URL**: `/api/products`
- **Permission**: `create-product` (Admin, Staff)
- **Request Body**:
  ```json
  {
    "name": "Mechanical Keyboard Custom",
    "description": "Keyboard mekanik 75% hot-swappable dengan gasket mount.",
    "price": 1850000,
    "stock": 12
  }
  ```

---

### • Product Update (Ubah Data Produk)
- **Method**: `PUT`
- **URL**: `/api/products/:id`
- **Permission**: `update-product` (Admin, Staff)
- **Request Body**:
  ```json
  {
    "name": "Laptop Asus ROG Zephyrus G16 (2026)",
    "description": "Versi upgrade dengan prosesor Intel Core Ultra 9.",
    "price": 28999000,
    "stock": 10
  }
  ```

---

### • Product Delete (Hapus Produk)
- **Method**: `DELETE`
- **URL**: `/api/products/:id`
- **Permission**: `delete-product` (**Hanya Admin**)

---

## 4. Role & Permission (Manajemen Hak Akses RBAC)

Semua endpoint di bawah ini memerlukan login dan hak akses admin (`read-role`, `create-role`, `update-role`, `delete-role`).

---

### • Permission Get (Daftar Semua Permission)
Melihat seluruh permission yang tersedia di dalam sistem agar bisa dipilih untuk dihubungkan ke role.
- **Method**: `GET`
- **URL**: `/api/permissions`
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `read-role`
- **Response `200 OK`**:
  ```json
  {
    "message": "Berhasil mengambil daftar permission",
    "data": [
      {
        "id": "3830afbd-65a8-4539-bec9-d3c4da141eda",
        "name": "read-product",
        "description": "Melihat daftar dan detail produk"
      },
      {
        "id": "3d1526ed-b52a-4ee3-8f75-2766ef2e71c2",
        "name": "create-product",
        "description": "Menambahkan produk baru"
      },
      {
        "id": "0e02310f-9734-4c52-a3e3-f41316111cdb",
        "name": "update-product",
        "description": "Memperbarui informasi produk"
      },
      {
        "id": "a99764aa-e977-416b-aa43-f39f858b5d51",
        "name": "delete-product",
        "description": "Menghapus produk"
      },
      {
        "id": "f10134ee-a0c1-48b2-b34c-31ec8fa181d2",
        "name": "read-role",
        "description": "Melihat daftar dan detail role"
      },
      {
        "id": "fe819678-b678-4e58-869e-76c8dbb0208c",
        "name": "create-role",
        "description": "Membuat role baru"
      },
      {
        "id": "f5addc9a-eba9-410d-a965-099a4a329fc4",
        "name": "update-role",
        "description": "Memperbarui role dan permissions-nya"
      },
      {
        "id": "fd9d2515-4a37-4982-ae17-41e7124238a4",
        "name": "delete-role",
        "description": "Menghapus role"
      }
    ]
  }
  ```

---

### • Role Get (Daftar Semua Role)
Mengambil semua role beserta daftar permission yang terkait pada masing-masing role.
- **Method**: `GET`
- **URL**: `/api/roles`
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `read-role`
- **Response `200 OK`**:
  ```json
  {
    "message": "Berhasil mengambil data role",
    "data": [
      {
        "id": "567e58bf-09da-4e99-9200-d59799ec7554",
        "name": "admin",
        "description": "Akses penuh sistem",
        "permissions": [
          { "id": "3830afbd-65a8-4539-bec9-d3c4da141eda", "name": "read-product" },
          { "id": "3d1526ed-b52a-4ee3-8f75-2766ef2e71c2", "name": "create-product" }
        ]
      },
      {
        "id": "a6a6c676-d10c-414d-a36e-a9604e76aa69",
        "name": "staff",
        "description": "Staf pengelola produk",
        "permissions": [
          { "id": "3830afbd-65a8-4539-bec9-d3c4da141eda", "name": "read-product" },
          { "id": "3d1526ed-b52a-4ee3-8f75-2766ef2e71c2", "name": "create-product" },
          { "id": "0e02310f-9734-4c52-a3e3-f41316111cdb", "name": "update-product" }
        ]
      }
    ]
  }
  ```

---

### • Role Get by ID (Detail Role)
Mengambil detail satu role berdasarkan ID UUID beserta permissions-nya.
- **Method**: `GET`
- **URL**: `/api/roles/:id`
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `read-role`
- **Response `200 OK`**:
  ```json
  {
    "message": "Detail role ditemukan",
    "data": {
      "id": "a6a6c676-d10c-414d-a36e-a9604e76aa69",
      "name": "staff",
      "description": "Staf pengelola produk",
      "permissions": [
        { "id": "3830afbd-65a8-4539-bec9-d3c4da141eda", "name": "read-product" },
        { "id": "3d1526ed-b52a-4ee3-8f75-2766ef2e71c2", "name": "create-product" }
      ]
    }
  }
  ```

---

### • Role Create (Tambah Role Baru)
Membuat role baru, dapat langsung melampirkan array UUID permissions yang diinginkan.
- **Method**: `POST`
- **URL**: `/api/roles`
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `create-role`
- **Request Body**:
  ```json
  {
    "name": "supervisor",
    "description": "Supervisor operasional cabang",
    "permission_ids": [
      "3830afbd-65a8-4539-bec9-d3c4da141eda",
      "3d1526ed-b52a-4ee3-8f75-2766ef2e71c2",
      "0e02310f-9734-4c52-a3e3-f41316111cdb"
    ]
  }
  ```
- **Response `201 Created`**:
  ```json
  {
    "message": "Role berhasil dibuat",
    "data": {
      "id": "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
      "name": "supervisor",
      "description": "Supervisor operasional cabang",
      "permissions": [ ... ]
    }
  }
  ```

---

### • Role Update (Ubah Info Role)
Mengubah nama atau deskripsi dari role.
- **Method**: `PUT`
- **URL**: `/api/roles/:id`
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `update-role`
- **Request Body**:
  ```json
  {
    "name": "senior-staff",
    "description": "Staf senior gudang & produk"
  }
  ```
- **Response `200 OK`**:
  ```json
  {
    "message": "Role berhasil diperbarui",
    "data": {
      "id": "a6a6c676-d10c-414d-a36e-a9604e76aa69",
      "name": "senior-staff",
      "description": "Staf senior gudang & produk"
    }
  }
  ```

---

### • Role Update Permissions (Assign Permissions ke Role)
Memperbarui/mengganti kumpulan permissions yang dimiliki oleh suatu role.
- **Method**: `PUT`
- **URL**: `/api/roles/:id/permissions`
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `update-role`
- **Request Body**:
  ```json
  {
    "permission_ids": [
      "3830afbd-65a8-4539-bec9-d3c4da141eda",
      "3d1526ed-b52a-4ee3-8f75-2766ef2e71c2",
      "0e02310f-9734-4c52-a3e3-f41316111cdb",
      "a99764aa-e977-416b-aa43-f39f858b5d51"
    ]
  }
  ```
- **Response `200 OK`**:
  ```json
  {
    "message": "Permissions untuk role berhasil diperbarui",
    "data": {
      "id": "a6a6c676-d10c-414d-a36e-a9604e76aa69",
      "name": "staff",
      "description": "Staf pengelola produk",
      "permissions": [
        { "id": "3830afbd-65a8-4539-bec9-d3c4da141eda", "name": "read-product" },
        { "id": "3d1526ed-b52a-4ee3-8f75-2766ef2e71c2", "name": "create-product" },
        { "id": "0e02310f-9734-4c52-a3e3-f41316111cdb", "name": "update-product" },
        { "id": "a99764aa-e977-416b-aa43-f39f858b5d51", "name": "delete-product" }
      ]
    }
  }
  ```

---

### • Role Delete (Hapus Role)
Menghapus role dari sistem. Terdapat proteksi keamanan: role bawaan sistem (`admin` / `user`) atau role yang masih digunakan oleh user aktif **tidak dapat dihapus**.
- **Method**: `DELETE`
- **URL**: `/api/roles/:id`
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `delete-role`
- **Response `200 OK`**:
  ```json
  {
    "message": "Role berhasil dihapus"
  }
  ```
- **Response `400 Bad Request`** (Jika mencoba hapus role admin/user):
  ```json
  {
    "error": "Role sistem bawaan ('admin') tidak dapat dihapus"
  }
  ```
- **Response `400 Bad Request`** (Jika role masih dipakai user):
  ```json
  {
    "error": "Tidak dapat menghapus role yang masih digunakan oleh pengguna aktif"
  }
  ```

---

## 5. User (Manajemen Pengguna)

Semua endpoint di bawah ini memerlukan login dan hak akses pengelolaan pengguna (`read-user`, `create-user`, `update-user`, `delete-user`).

---

### • User Get (Daftar Pengguna)
Mengambil seluruh daftar pengguna. Mendukung query parameter `?search=` untuk memfilter berdasarkan nama atau email.
- **Method**: `GET`
- **URL**: `/api/users`
- **Query Parameter (Opsional)**: `?search=budi`
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `read-user`
- **Response `200 OK`**:
  ```json
  {
    "message": "Berhasil mengambil data pengguna",
    "data": [
      {
        "id": "e6deafb6-2a97-4d62-bf6b-827da8edebd0",
        "name": "Admin Toko",
        "email": "admin@example.com",
        "role_id": "3ceba60e-23b3-4c4a-97eb-e8e8604c6371",
        "role": {
          "id": "3ceba60e-23b3-4c4a-97eb-e8e8604c6371",
          "name": "admin",
          "description": "Akses penuh sistem"
        },
        "created_at": "2026-10-09T13:03:12+07:00"
      }
    ]
  }
  ```

---

### • User Get by ID (Detail Pengguna)
Mengambil detail satu pengguna berdasarkan ID UUID.
- **Method**: `GET`
- **URL**: `/api/users/:id`
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `read-user`
- **Response `200 OK`**:
  ```json
  {
    "message": "Detail pengguna ditemukan",
    "data": {
      "id": "5c3a46b8-096d-4d16-a3ad-e48386bdbfb5",
      "name": "Staff Gudang",
      "email": "staff@example.com",
      "role_id": "5e82ebef-0edb-42d2-9b29-da807d647954",
      "role": {
        "id": "5e82ebef-0edb-42d2-9b29-da807d647954",
        "name": "staff"
      }
    }
  }
  ```

---

### • User Create (Tambah Pengguna oleh Admin)
Membuat akun pengguna baru dan langsung menentukan role-nya.
- **Method**: `POST`
- **URL**: `/api/users`
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `create-user`
- **Request Body**:
  ```json
  {
    "name": "Dian Sastro",
    "email": "dian@example.com",
    "password": "password123",
    "role_id": "5e82ebef-0edb-42d2-9b29-da807d647954"
  }
  ```
- **Response `201 Created`**:
  ```json
  {
    "message": "Pengguna baru berhasil dibuat",
    "data": {
      "id": "a9182736-45ef-4b2a-8910-bcde12345678",
      "name": "Dian Sastro",
      "email": "dian@example.com",
      "role_id": "5e82ebef-0edb-42d2-9b29-da807d647954",
      "role": {
        "id": "5e82ebef-0edb-42d2-9b29-da807d647954",
        "name": "staff"
      }
    }
  }
  ```

---

### • User Update (Ubah Data Profil Pengguna)
Memperbarui informasi nama, email, atau reset kata sandi pengguna.
- **Method**: `PUT`
- **URL**: `/api/users/:id`
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `update-user`
- **Request Body**:
  ```json
  {
    "name": "Dian Sastrowardoyo",
    "email": "dian.new@example.com",
    "password": "newpassword123"
  }
  ```
- **Response `200 OK`**:
  ```json
  {
    "message": "Data pengguna berhasil diperbarui",
    "data": {
      "id": "a9182736-45ef-4b2a-8910-bcde12345678",
      "name": "Dian Sastrowardoyo",
      "email": "dian.new@example.com"
    }
  }
  ```

---

### • User Update Role (Ubah Role Pengguna)
Khusus mengubah / memindahkan role dari seorang pengguna.
- **Method**: `PUT`
- **URL**: `/api/users/:id/role`
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `update-user`
- **Request Body**:
  ```json
  {
    "role_id": "3ceba60e-23b3-4c4a-97eb-e8e8604c6371"
  }
  ```
- **Response `200 OK`**:
  ```json
  {
    "message": "Role pengguna berhasil diperbarui",
    "data": {
      "id": "a9182736-45ef-4b2a-8910-bcde12345678",
      "name": "Dian Sastrowardoyo",
      "role_id": "3ceba60e-23b3-4c4a-97eb-e8e8604c6371",
      "role": {
        "id": "3ceba60e-23b3-4c4a-97eb-e8e8604c6371",
        "name": "admin"
      }
    }
  }
  ```

---

### • User Delete (Hapus Pengguna)
Menghapus pengguna dari sistem. Terdapat proteksi keamanan: user tidak dapat menghapus akunnya sendiri yang sedang aktif digunakan untuk login.
- **Method**: `DELETE`
- **URL**: `/api/users/:id`
- **Autentikasi**: Wajib Login (`Bearer Token`)
- **Permission Diperlukan**: `delete-user`
- **Response `200 OK`**:
  ```json
  {
    "message": "Pengguna berhasil dihapus"
  }
  ```
- **Response `400 Bad Request`** (Jika mencoba hapus akun sendiri):
  ```json
  {
    "error": "Anda tidak dapat menghapus akun Anda sendiri yang sedang aktif"
  }
  ```
