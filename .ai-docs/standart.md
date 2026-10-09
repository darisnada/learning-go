# Standar Pengembangan REST API (Wajib)
## Framework: Go (Golang), Gin, GORM, MySQL & JWT Authentication

Dokumen ini merupakan pedoman dan standar baku pengembangan REST API yang **wajib dipatuhi** untuk memastikan seluruh layanan dan fitur memiliki pola arsitektur yang seragam, aman, konsisten, dan mudah dipelihara.

---

## 1. Arsitektur & Struktur Direktori

Kode proyek harus dipisahkan berdasarkan tanggung jawab (*Separation of Concerns*):

```text
├── config/             # Manajemen env & konfigurasi koneksi database
├── controllers/        # Layer HTTP Handler (menerima request, validasi input, memanggil logic)
├── middlewares/        # HTTP Middleware (Autentikasi JWT, Otorisasi RBAC Permission, CORS)
├── models/             # Definisi Struct Entity Database (GORM, Hook UUID, Tag validasi)
├── seeders/            # Pengisian data awal (Roles, Permissions, Akun Demo, Data Dummy)
├── utils/              # Helper bersama (Format Response, JWT Generator, Custom Validator)
├── document-api.md     # Dokumentasi endpoint API per fitur
├── postman_collection.json # Koleksi pengujian otomatis Postman
├── .env.example        # Template konfigurasi environment yang aman untuk Git
├── .gitignore          # Daftar file yang wajib diabaikan (termasuk .env)
├── main.go             # Entry point aplikasi & registrasi routing
└── go.mod / go.sum     # Manajemen dependensi Go
```

---

## 2. Standardisasi Format JSON Response

Seluruh endpoint **wajib** mengembalikan struktur JSON yang konsisten, baik saat request berhasil maupun gagal.

### A. Format Respon Sukses (`HTTP 200 / 201`)
Setiap respon berhasil wajib menyertakan flag `"success": true`:
```json
{
  "success": true,
  "message": "Data produk berhasil diambil",
  "data": {
    "id": "8d78882e-42cb-4974-801b-cad0b8ffce38",
    "name": "Laptop Gaming",
    "price": 15000000,
    "stock": 10
  }
}
```

### B. Format Respon Gagal / Error (`HTTP 4xx / 5xx`)
Setiap respon error wajib menyertakan flag `"success": false` dan pesan `"message"` yang jelas bagi client/frontend:
```json
{
  "success": false,
  "message": "Validasi input gagal",
  "errors": {
    "Name": "Field 'Name' wajib diisi",
    "Price": "Field 'Price' harus lebih besar dari 0"
  }
}
```

---

## 3. Standardisasi Penggunaan UUID (Primary & Foreign Key)

1. **Seluruh Entitas Wajib Menggunakan UUID v4**:
   - Dilarang menggunakan integer auto-increment untuk entitas publik agar ID tidak mudah ditebak (*enumeration attack*).
   - Tipe data di MySQL adalah `CHAR(36)`.
2. **Hook Otomatis Sebelum Simpan (`BeforeCreate`)**:
   Setiap struct model GORM wajib mengimplementasikan hook otomatis pembuatan UUID:
   ```go
   func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
       if m.ID == "" {
           m.ID = uuid.NewString()
       }
       return
   }
   ```
3. **Validasi URL Parameter**:
   Setiap handler yang menerima ID dari URL wajib memverifikasi validitas format UUID sebelum melakukan query ke database:
   ```go
   if _, err := uuid.Parse(id); err != nil {
       utils.ErrorResponse(c, http.StatusBadRequest, "Format ID tidak valid (harus berupa UUID)")
       return
   }
   ```

---

## 4. Keamanan & Autentikasi JWT

1. **Header Standar**:
   Setiap request ke route terproteksi wajib menyertakan header:
   ```http
   Authorization: Bearer <TOKEN_JWT>
   ```
2. **Claims Minimal & Aman**:
   - Hanya simpan informasi penting (seperti `user_id` dan `email`).
   - **Dilarang keras** menyimpan data sensitif (password hash, kartu kredit, pin) di dalam payload JWT.
3. **Masa Berlaku Token**:
   - Token harus memiliki masa berlaku (*expiration time*) yang jelas (default: 24 jam untuk pengembangan, disarankan 1-2 jam di server produksi).
4. **Kerahasiaan Kunci**:
   - `JWT_SECRET` **wajib dimuat dari `.env`**. Dilarang melakukan *hardcode* kunci rahasia di dalam file kode Go.

---

## 5. Otorisasi Berbasis Hak Akses (RBAC Permission)

1. **Penulisan Permission Standar**:
   Permission wajib ditulis dengan pola `<action>-<resource>` menggunakan huruf kecil dan tanda hubung:
   - `read-product`
   - `create-product`
   - `update-product`
   - `delete-product`
2. **Pengecekan Berlapis**:
   - Lapisan 1: `middlewares.AuthMiddleware()` memastikan token valid dan user terotentikasi.
   - Lapisan 2: `middlewares.RequirePermission("permission-name")` memastikan role user memiliki hak akses terhadap tindakan tersebut.
3. **Respon Hak Akses Ditolak**:
   Wajib mengembalikan HTTP Status Code **`403 Forbidden`**.

---

## 6. Aturan Penanganan Error & Integritas Data

1. **Prinsip *Early Return***:
   Lakukan validasi secepat mungkin. Jika validasi gagal, segera lakukan `return` sebelum menyentuh operasi database.
2. **Pencegahan Error Update**:
   - Operasi update hanya boleh dieksekusi jika data target ditemukan dan seluruh input valid.
   - Jika terjadi error di tengah jalan, database GORM/MySQL secara otomatis melakukan **Rollback** transaksi sehingga data lama tidak berubah sebagian (*atomic*).
3. **Penyembunyian Detail Database Internal**:
   - Jangan mengembalikan pesan error internal SQL mentah (seperti SQL syntax error atau driver crash) langsung ke client untuk mencegah kebocoran informasi arsitektur database.
4. **Global Route Handling**:
   - Wajib menangani URL yang tidak ditemukan (`404`) dan metode HTTP yang salah (`405`) dalam format respons JSON standar (bukan halaman HTML bawaan web server).

---

## 7. Migrasi & Seeder Database

1. **Dukungan CLI Flags**:
   Aplikasi wajib mendukung perintah pemeliharaan database langsung lewat terminal:
   - `go run main.go -seed` : Menjalankan seeder data awal.
   - `go run main.go -fresh` : Menghapus tabel lama & membuat ulang schema (*fresh migration*).
   - `go run main.go -fresh-seed` : Reset tabel dan isi ulang seeder secara lengkap.
2. **Idempoten & Tanpa Log Palsu**:
   Query pengecekan data awal di seeder tidak boleh memicu log peringatan merah (`record not found`). Gunakan metode `Find()` atau konfigurasi `IgnoreRecordNotFoundError: true`.

---

## 8. Checklist Kesiapan Produksi (*Production Readiness*)

Sebelum kode di-merge ke branch `main` atau dirilis ke server produksi, pastikan:

- [ ] File `.env` masuk ke dalam `.gitignore` dan tidak pernah ter-commit ke repositori Git.
- [ ] Tersedia file `.env.example` yang mencantumkan seluruh variabel yang dibutuhkan tanpa memuat kredensial asli.
- [ ] Seluruh endpoint memiliki dokumentasi teknis di `document-api.md`.
- [ ] Tersedia file `postman_collection.json` untuk pengujian otomatis.
- [ ] Kode lulus kompilasi tanpa warning/error (`go build .`).
- [ ] Dependensi bersih dan terdaftar rapi (`go mod tidy`).

