package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name" binding:"required"`
	Email     string         `gorm:"type:varchar(191);uniqueIndex;not null" json:"email" binding:"required,email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"password,omitempty" binding:"required,min=6"`
	RoleID    uint           `gorm:"not null;default:3" json:"role_id"`
	Role      *Role          `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// HashPassword hashes user's plain text password using bcrypt
func (u *User) HashPassword(password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedBytes)
	return nil
}

// CheckPassword checks if input plain password matches hashed password
func (u *User) CheckPassword(providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(providedPassword))
}

// HasPermission mengecek apakah user memiliki permission tertentu melalui Role-nya
func (u *User) HasPermission(permissionName string) bool {
	if u.Role == nil {
		return false
	}
	// Role admin memiliki akses ke semua permission
	if u.Role.Name == "admin" {
		return true
	}
	for _, perm := range u.Role.Permissions {
		if perm.Name == permissionName {
			return true
		}
	}
	return false
}
