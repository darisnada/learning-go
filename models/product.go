package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          string         `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(150);not null" json:"name" binding:"required"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:decimal(12,2);not null" json:"price" binding:"required,gt=0"`
	Stock       int            `gorm:"type:int;not null;default:0" json:"stock" binding:"gte=0"`
	UserID      string         `gorm:"type:char(36);index" json:"user_id"`
	User        *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	return
}
