package seeders

import (
	"log"

	"go-learning/models"

	"gorm.io/gorm"
)

// SeedAll menjalankan seluruh proses seeding data
func SeedAll(db *gorm.DB) {
	log.Println("Memulai proses seeding database...")
	permissions := SeedPermissions(db)
	roles := SeedRoles(db, permissions)
	adminUser := SeedUsers(db, roles)
	SeedProducts(db, adminUser)
	log.Println("Proses seeding selesai!")
}

// SeedPermissions membuat permission default untuk produk, role, dan user
func SeedPermissions(db *gorm.DB) map[string]models.Permission {
	permMap := make(map[string]models.Permission)

	permissions := []models.Permission{
		// Permissions Produk
		{Name: "read-product", Description: "Melihat daftar dan detail produk"},
		{Name: "create-product", Description: "Menambahkan produk baru"},
		{Name: "update-product", Description: "Memperbarui informasi produk"},
		{Name: "delete-product", Description: "Menghapus produk"},

		// Permissions Role Management
		{Name: "read-role", Description: "Melihat daftar dan detail role"},
		{Name: "create-role", Description: "Membuat role baru"},
		{Name: "update-role", Description: "Memperbarui role dan permissions-nya"},
		{Name: "delete-role", Description: "Menghapus role"},

		// Permissions User Management
		{Name: "read-user", Description: "Melihat daftar dan detail pengguna"},
		{Name: "create-user", Description: "Membuat akun pengguna baru"},
		{Name: "update-user", Description: "Memperbarui profil atau role pengguna"},
		{Name: "delete-user", Description: "Menghapus akun pengguna"},
	}

	for _, p := range permissions {
		var existingPerms []models.Permission
		db.Where("name = ?", p.Name).Limit(1).Find(&existingPerms)

		if len(existingPerms) == 0 {
			if err := db.Create(&p).Error; err == nil {
				permMap[p.Name] = p
				log.Printf("[Seeder] Permission dibuat: %s (%s)", p.Name, p.ID)
			}
		} else {
			permMap[p.Name] = existingPerms[0]
		}
	}

	return permMap
}

// SeedRoles membuat role admin, staff, dan user dengan permissions masing-masing
func SeedRoles(db *gorm.DB, perms map[string]models.Permission) map[string]models.Role {
	roleMap := make(map[string]models.Role)

	pReadProd := perms["read-product"]
	pCreateProd := perms["create-product"]
	pUpdateProd := perms["update-product"]
	pDeleteProd := perms["delete-product"]

	pReadRole := perms["read-role"]
	pCreateRole := perms["create-role"]
	pUpdateRole := perms["update-role"]
	pDeleteRole := perms["delete-role"]

	pReadUser := perms["read-user"]
	pCreateUser := perms["create-user"]
	pUpdateUser := perms["update-user"]
	pDeleteUser := perms["delete-user"]

	// Role Admin: Memiliki SEMUA permission
	adminRole := models.Role{
		Name:        "admin",
		Description: "Akses penuh sistem",
		Permissions: []models.Permission{
			pReadProd, pCreateProd, pUpdateProd, pDeleteProd,
			pReadRole, pCreateRole, pUpdateRole, pDeleteRole,
			pReadUser, pCreateUser, pUpdateUser, pDeleteUser,
		},
	}

	// Role Staff: read, create, update produk
	staffRole := models.Role{
		Name:        "staff",
		Description: "Staf pengelola produk",
		Permissions: []models.Permission{pReadProd, pCreateProd, pUpdateProd},
	}

	// Role User: hanya read-product
	userRole := models.Role{
		Name:        "user",
		Description: "Pengguna biasa",
		Permissions: []models.Permission{pReadProd},
	}

	roles := []models.Role{adminRole, staffRole, userRole}

	for _, r := range roles {
		var existingRoles []models.Role
		db.Where("name = ?", r.Name).Limit(1).Find(&existingRoles)

		if len(existingRoles) == 0 {
			if err := db.Create(&r).Error; err == nil {
				roleMap[r.Name] = r
				log.Printf("[Seeder] Role '%s' (%s) dibuat.", r.Name, r.ID)
			}
		} else {
			existing := existingRoles[0]
			db.Model(&existing).Association("Permissions").Replace(r.Permissions)
			roleMap[r.Name] = existing
		}
	}

	return roleMap
}

// SeedUsers menambahkan akun pengguna awal (Admin, Staff, User biasa)
func SeedUsers(db *gorm.DB, roles map[string]models.Role) *models.User {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("[Seeder] Tabel users sudah memiliki data, melewati seed user.")
		var admin models.User
		db.Where("email = ?", "admin@example.com").First(&admin)
		return &admin
	}

	adminRoleID := roles["admin"].ID
	staffRoleID := roles["staff"].ID
	userRoleID := roles["user"].ID

	usersToSeed := []struct {
		user     models.User
		password string
	}{
		{
			user: models.User{
				Name:   "Admin Toko",
				Email:  "admin@example.com",
				RoleID: adminRoleID,
			},
			password: "admin123",
		},
		{
			user: models.User{
				Name:   "Staff Gudang",
				Email:  "staff@example.com",
				RoleID: staffRoleID,
			},
			password: "staff123",
		},
		{
			user: models.User{
				Name:   "Pelanggan Biasa",
				Email:  "user@example.com",
				RoleID: userRoleID,
			},
			password: "user123",
		},
	}

	var firstAdmin *models.User
	for i, item := range usersToSeed {
		u := item.user
		if err := u.HashPassword(item.password); err != nil {
			log.Printf("[Seeder] Gagal hash password: %v", err)
			continue
		}

		if err := db.Create(&u).Error; err != nil {
			log.Printf("[Seeder] Gagal membuat user '%s': %v", u.Email, err)
		} else {
			log.Printf("[Seeder] Akun dibuat: %s | ID: %s | Role: %s | Pass: %s", u.Email, u.ID, u.RoleID, item.password)
			if i == 0 {
				firstAdmin = &u
			}
		}
	}

	return firstAdmin
}

// SeedProducts menambahkan daftar produk awal
func SeedProducts(db *gorm.DB, user *models.User) {
	var count int64
	db.Model(&models.Product{}).Count(&count)
	if count > 0 {
		log.Println("[Seeder] Tabel products sudah memiliki data, melewati seed produk.")
		return
	}

	var userID string
	if user != nil && user.ID != "" {
		userID = user.ID
	}

	sampleProducts := []models.Product{
		{
			Name:        "Laptop Asus ROG Zephyrus",
			Description: "Laptop gaming bertenaga tinggi dengan layar 165Hz dan RTX 4070.",
			Price:       24999000,
			Stock:       15,
			Image:       "/uploads/products/laptop-rog.jpg",
			UserID:      userID,
		},
		{
			Name:        "Mechanical Keyboard Keychron K2",
			Description: "Keyboard mekanik wireless dengan switch Gateron Brown dan lampu RGB.",
			Price:       1250000,
			Stock:       30,
			Image:       "/uploads/products/keyboard-keychron.jpg",
			UserID:      userID,
		},
		{
			Name:        "Logitech MX Master 3S",
			Description: "Mouse ergonomis nirkabel dengan klik senyap dan scroll elektromagnetik.",
			Price:       1550000,
			Stock:       25,
			Image:       "/uploads/products/mouse-mxmaster.jpg",
			UserID:      userID,
		},
		{
			Name:        "Monitor LG UltraGear 27 Inch",
			Description: "Monitor gaming IPS 144Hz 1ms dengan resolusi QHD (2K).",
			Price:       4350000,
			Stock:       10,
			Image:       "/uploads/products/monitor-lg.jpg",
			UserID:      userID,
		},
		{
			Name:        "Sony WH-1000XM5 Headphone",
			Description: "Headphone wireless dengan fitur peredam bising (ANC) terbaik.",
			Price:       4999000,
			Stock:       8,
			Image:       "/uploads/products/headphone-sony.jpg",
			UserID:      userID,
		},
	}

	for _, prod := range sampleProducts {
		if err := db.Create(&prod).Error; err != nil {
			log.Printf("[Seeder] Gagal membuat produk '%s': %v", prod.Name, err)
		} else {
			log.Printf("[Seeder] Produk berhasil dibuat: %s (UUID: %s)", prod.Name, prod.ID)
		}
	}
}
