package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID          string         `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Description string         `gorm:"type:varchar(255)" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	return
}

type Role struct {
	ID          string         `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description string         `gorm:"type:varchar(255)" json:"description"`
	Permissions []Permission   `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return
}
