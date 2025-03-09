package repositories

import (
	"fmt"

	"github.com/PharmaKart/product-svc/internal/models"
	"github.com/PharmaKart/product-svc/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *models.Product) (string, error)
	GetProduct(id string) (*models.Product, error)
	GetProductByName(name string) (*models.Product, error)
	ListProducts(page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Product, int32, error)
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
		if err.Error() != "record not found" {
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

func (r *productRepository) ListProducts(page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Product, int32, error) {
	var products []models.Product
	var total int64

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	query := r.db
	if filter != "" && filterValue != "" {
		query = query.Where(filter+" = ?", filterValue)
	}

	if sortBy != "" {
		if sortOrder == "" {
			sortOrder = "asc"
		}
		query = query.Order(sortBy + " " + sortOrder)
	}

	err := query.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&products).Error
	if err != nil {
		return nil, 0, errors.NewInternalError(err)
	}

	err = query.Model(&models.Product{}).Count(&total).Error
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
