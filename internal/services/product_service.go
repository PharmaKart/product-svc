package services

import (
	"github.com/PharmaKart/product-svc/internal/models"
	"github.com/PharmaKart/product-svc/internal/repositories"
	"github.com/PharmaKart/product-svc/pkg/utils"
)

type ProductService interface {
	CreateProduct(product *models.Product) (string, error)
	GetProduct(id string) (*models.Product, error)
	ListProducts(search string, filters models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.Product, int32, error)
	UpdateProduct(id string, name string, description string, price float64, requiresPrescription bool, imageURL string) error
	DeleteProduct(id string) error
	UpdateStock(log *models.InventoryLog) error
	GetInventoryLogs(productID string, filters models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.InventoryLog, int32, error)
}

type productService struct {
	ProductRepository      repositories.ProductRepository
	InventoryLogRepository repositories.InventoryLogRepository
}

func NewProductService(productRepository repositories.ProductRepository, inventoryLogRepository repositories.InventoryLogRepository) ProductService {
	return &productService{
		ProductRepository:      productRepository,
		InventoryLogRepository: inventoryLogRepository,
	}
}

func (s *productService) CreateProduct(product *models.Product) (string, error) {
	// Validate the product input
	if err := utils.ValidateProductInput(product); err != nil {
		return "", err
	}

	// Add the product to the database
	productID, err := s.ProductRepository.CreateProduct(product)
	if err != nil {
		return "", err
	}
	return productID, nil
}

func (s *productService) GetProduct(id string) (*models.Product, error) {
	product, err := s.ProductRepository.GetProduct(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) ListProducts(search string, filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.Product, int32, error) {
	products, total, err := s.ProductRepository.ListProducts(search, filter, sortBy, sortOrder, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (s *productService) UpdateProduct(id string, name string, description string, price float64, requiresPrescription bool, imageURL string) error {
	// Get the product from the database
	product, err := s.ProductRepository.GetProduct(id)
	if err != nil {
		return err
	}

	// Update the product fields
	product.Name = name
	product.Description = &description
	product.Price = price
	product.RequiresPrescription = requiresPrescription
	if imageURL != "" {
		product.ImageURL = &imageURL
	}

	// Validate the product input
	if err := utils.ValidateProductInput(product); err != nil {
		return err
	}

	// Update the product in the database
	if err := s.ProductRepository.UpdateProduct(product); err != nil {
		return err
	}
	return nil
}

func (s *productService) DeleteProduct(id string) error {
	// Delete the product from the database
	if err := s.ProductRepository.DeleteProduct(id); err != nil {
		return err
	}
	return nil
}

func (s *productService) UpdateStock(log *models.InventoryLog) error {
	// Validate the inventory input
	if err := utils.ValidateInventoryInput(log); err != nil {
		return err
	}

	// Update the stock in the database
	if err := s.ProductRepository.UpdateStock(log.ProductID, log.QuantityChange); err != nil {
		return err
	}

	// Log the inventory change
	if err := s.InventoryLogRepository.LogChange(log); err != nil {
		return err
	}

	return nil
}

func (s *productService) GetInventoryLogs(productID string, filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.InventoryLog, int32, error) {
	logs, total, err := s.InventoryLogRepository.GetLogsByProductID(productID, filter, sortBy, sortOrder, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
