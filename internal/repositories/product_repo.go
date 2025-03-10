package repositories

import (
	"fmt"
	"strings"

	"github.com/PharmaKart/product-svc/internal/models"
	"github.com/PharmaKart/product-svc/pkg/errors"
	"github.com/PharmaKart/product-svc/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *models.Product) (string, error)
	GetProduct(id string) (*models.Product, error)
	GetProductByName(name string) (*models.Product, error)
	ListProducts(filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.Product, int32, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id string) error
	UpdateStock(id uuid.UUID, quantity int) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) CreateProduct(product *models.Product) (string, error) {
	existingProduct, err := r.GetProductByName(product.Name)
	if err != nil {
		if err.Error() != fmt.Sprintf("Product with name '%s' not found", product.Name) {
			return "", errors.NewInternalError(err)
		}
	}

	if existingProduct != nil && existingProduct.ID != uuid.Nil {
		return "", errors.NewConflictError(fmt.Sprintf("Product with name '%s' already exists", product.Name))
	}

	if err := r.db.Create(product).Error; err != nil {
		return "", errors.NewInternalError(err)
	}
	return product.ID.String(), nil
}

func (r *productRepository) GetProduct(id string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("id = ?", id).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError(fmt.Sprintf("Product with ID '%s' not found", id))
		}
		return nil, errors.NewInternalError(err)
	}
	return &product, nil
}

func (r *productRepository) GetProductByName(name string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("name = ?", name).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError(fmt.Sprintf("Product with name '%s' not found", name))
		}
		return nil, errors.NewInternalError(err)
	}
	return &product, nil
}

func (r *productRepository) ListProducts(filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.Product, int32, error) {
	var products []models.Product
	var total int64

	allowedColumns := utils.GetModelColumns(&models.Product{})

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

	query := r.db.Model(&models.Product{})

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

	err = query.Find(&products).Error
	if err != nil {
		return nil, 0, errors.NewInternalError(err)
	}

	return products, int32(total), nil
}

func (r *productRepository) UpdateProduct(product *models.Product) error {
	_, err := r.GetProduct(product.ID.String())
	if err != nil {
		return err
	}

	existingProduct, err := r.GetProductByName(product.Name)
	if err != nil {
		if _, ok := errors.IsAppError(err); !ok && existingProduct != nil && existingProduct.ID != product.ID {
			return errors.NewConflictError(fmt.Sprintf("Product with name '%s' already exists", product.Name))
		}
	}

	if err := r.db.Save(product).Error; err != nil {
		return errors.NewInternalError(err)
	}

	return nil
}

func (r *productRepository) DeleteProduct(id string) error {
	_, err := r.GetProduct(id)
	if err != nil {
		return err
	}

	if err := r.db.Where("id = ?", id).Delete(&models.Product{}).Error; err != nil {
		return errors.NewInternalError(err)
	}

	return nil
}

func (r *productRepository) UpdateStock(id uuid.UUID, quantity int) error {
	_, err := r.GetProduct(id.String())
	if err != nil {
		return err
	}

	result := r.db.Model(&models.Product{}).Where("id = ?", id).Update("stock", gorm.Expr("stock + ?", quantity))
	if result.Error != nil {
		return errors.NewInternalError(result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.NewNotFoundError(fmt.Sprintf("Product with ID '%s' not found", id))
	}

	return nil
}
