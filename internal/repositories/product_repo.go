package repositories

import (
	"github.com/PharmaKart/product-svc/internal/models"
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
	if err := r.db.Create(product).Error; err != nil {
		return "", err
	}
	return product.ID.String(), nil
}

func (r *productRepository) GetProduct(id string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("id = ?", id).First(&product).Error
	return &product, err
}

func (r *productRepository) GetProductByName(name string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("name = ?", name).First(&product).Error
	return &product, err
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
		return nil, 0, err
	}

	err = query.Model(&models.Product{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return products, int32(total), err

}

func (r *productRepository) UpdateProduct(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) DeleteProduct(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Product{}).Error
}

func (r *productRepository) UpdateStock(id uuid.UUID, quantity int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Update("stock", gorm.Expr("stock - ?", quantity)).Error
}
