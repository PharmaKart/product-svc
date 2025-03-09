package repositories

import (
	"fmt"

	"github.com/PharmaKart/product-svc/internal/models"
	"github.com/PharmaKart/product-svc/pkg/errors"
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
	if err := r.db.Create(log).Error; err != nil {
		return errors.NewInternalError(err)
	}
	return nil
}

func (r *inventoryLogRepository) GetLogsByProductID(productID string) ([]models.InventoryLog, error) {
	var logs []models.InventoryLog

	if err := r.db.Where("product_id = ?", productID).Find(&logs).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError(fmt.Sprintf("No inventory logs found for product ID '%s'", productID))
		}
		return nil, errors.NewInternalError(err)
	}

	if len(logs) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("No inventory logs found for product ID '%s'", productID))
	}

	return logs, nil
}
