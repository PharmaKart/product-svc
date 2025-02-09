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
	ListProducts() ([]models.Product, error)
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

func (r *productRepository) ListProducts() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error
	return products, err
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
