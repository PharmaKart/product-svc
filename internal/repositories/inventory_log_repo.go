package repositories

import (
	"github.com/PharmaKart/product-svc/internal/models"
	"gorm.io/gorm"
)

type InventoryLogRepository interface {
	LogChange(log *models.InventoryLog) error
	GetLogsByProductID(productID string) ([]models.InventoryLog, error)
}

type inventoryLogRepository struct {
	db *gorm.DB
}

func NewInventoryLogRepository(db *gorm.DB) InventoryLogRepository {
	return &inventoryLogRepository{db}
}

func (r *inventoryLogRepository) LogChange(log *models.InventoryLog) error {
	return r.db.Create(log).Error
}

func (r *inventoryLogRepository) GetLogsByProductID(productID string) ([]models.InventoryLog, error) {
	var logs []models.InventoryLog
	err := r.db.Where("product_id = ?", productID).Find(&logs).Error
	return logs, err
}
