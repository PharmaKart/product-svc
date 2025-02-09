package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryLog struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProductID  uuid.UUID `gorm:"not null"`
	ChangeType string    `gorm:"type:varchar(50);not null;check:change_type IN ('order_placed', 'order_cancelled', 'stock_added')"`
	Quantity   int       `gorm:"not null"`
	CreatedAt  time.Time `gorm:"default:now()"`
}

func (il *InventoryLog) BeforeCreate(tx *gorm.DB) (err error) {
	il.ID = uuid.New()
	return
}
