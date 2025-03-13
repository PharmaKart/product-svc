package repositories

import (
	"strings"

	"github.com/PharmaKart/product-svc/internal/models"
	"github.com/PharmaKart/product-svc/pkg/errors"
	"github.com/PharmaKart/product-svc/pkg/utils"
	"gorm.io/gorm"
)

type InventoryLogRepository interface {
	LogChange(log *models.InventoryLog) error
	GetLogsByProductID(productID string, filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.InventoryLog, int32, error)
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

func (r *inventoryLogRepository) GetLogsByProductID(productID string, filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.InventoryLog, int32, error) {
	var logs []models.InventoryLog
	var total int64

	allowedColumns := utils.GetModelColumns(&models.InventoryLog{})

	allowedOperators := map[string]string{
		"eq":      "=",           // Equal to
		"neq":     "!=",          // Not equal to
		"gt":      ">",           // Greater than
		"gte":     ">=",          // Greater than or equal to
		"lt":      "<",           // Less than
		"lte":     "<=",          // Less than or equal to
		"like":    "LIKE",        // LIKE for pattern matching
		"ilike":   "ILIKE",       // Case insensitive LIKE (for PostgreSQL)
		"in":      "IN",          // IN for multiple values
		"null":    "IS NULL",     // IS NULL check
		"notnull": "IS NOT NULL", // IS NOT NULL check
	}

	query := r.db.Model(&models.InventoryLog{})
	if filter != (models.Filter{}) {
		if _, allowed := allowedColumns[filter.Column]; !allowed {
			return nil, 0, errors.NewBadRequestError("invalid filter column: " + filter.Column)
		}

		op, allowed := allowedOperators[filter.Operator]
		if !allowed {
			return nil, 0, errors.NewBadRequestError("invalid filter operator: " + filter.Operator)
		}

		switch filter.Operator {
		case "like", "ilike":
			query = query.Where(filter.Column+" "+op+" ?", "%"+filter.Value+"%")
		case "in":
			values := strings.Split(filter.Value, ",")
			query = query.Where(filter.Column+" "+op+" (?)", values)
		case "null", "notnull":
			query = query.Where(filter.Column + " " + op)
		default:
			query = query.Where(filter.Column+" "+op+" ?", filter.Value)
		}
	}

	if sortBy != "" {
		if _, allowed := allowedColumns[sortBy]; !allowed {
			return nil, 0, errors.NewBadRequestError("invalid sort column: " + sortBy)
		}

		sortOrder = strings.ToLower(sortOrder)
		if sortOrder != "asc" && sortOrder != "desc" {
			sortOrder = "asc"
		}

		query = query.Order(sortBy + " " + sortOrder)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, errors.NewInternalError(err)
	}

	if limit > 0 {
		offset := max(int((page-1)*limit), 0)
		query = query.Offset(offset).Limit(int(limit))
	}

	err = query.Find(&logs).Error
	if err != nil {
		return nil, 0, errors.NewInternalError(err)
	}

	return logs, int32(total), nil
}
