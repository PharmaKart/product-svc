package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID                   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name                 string    `gorm:"not null"`
	Description          *string
	Price                float64 `gorm:"not null"`
	Stock                int     `gorm:"not null;check:stock >= 0"`
	RequiresPrescription bool    `gorm:"default:false"`
	ImageURL             *string
	CreatedAt            time.Time `gorm:"default:now()"`
	UpdatedAt            time.Time `gorm:"default:now()"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
