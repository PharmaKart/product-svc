package handlers

import (
	"context"
	"fmt"

	"github.com/PharmaKart/product-svc/internal/models"
	"github.com/PharmaKart/product-svc/internal/proto"
	"github.com/PharmaKart/product-svc/internal/repositories"
	"github.com/PharmaKart/product-svc/internal/services"
	"github.com/google/uuid"
)

type ProductHandler interface {
	CreateProduct(ctx context.Context, req *proto.CreateProductRequest) (*proto.CreateProductResponse, error)
	GetProduct(ctx context.Context, req *proto.GetProductRequest) (*proto.GetProductResponse, error)
	ListProducts(ctx context.Context, req *proto.ListProductsRequest) (*proto.ListProductsResponse, error)
	UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest) (*proto.UpdateProductResponse, error)
	DeleteProduct(ctx context.Context, req *proto.DeleteProductRequest) (*proto.DeleteProductResponse, error)
	UpdateStock(ctx context.Context, req *proto.UpdateStockRequest) (*proto.UpdateStockResponse, error)
}

type productHandler struct {
	proto.UnimplementedProductServiceServer
	ProductService services.ProductService
}

func NewProductHandler(productRepo repositories.ProductRepository, inventorylogRepo repositories.InventoryLogRepository) *productHandler {
	return &productHandler{
		ProductService: services.NewProductService(productRepo, inventorylogRepo),
	}
}

func (h *productHandler) CreateProduct(ctx context.Context, req *proto.CreateProductRequest) (*proto.CreateProductResponse, error) {
	product := &models.Product{
		Name:                 req.Product.Name,
		Description:          &req.Product.Description,
		Price:                req.Product.Price,
		Stock:                int(req.Product.Stock),
		RequiresPrescription: req.Product.RequiresPrescription,
		ImageURL:             &req.Product.ImageUrl,
	}

	productID, err := h.ProductService.CreateProduct(product)
	if err != nil {
		return nil, err
	}

	return &proto.CreateProductResponse{
		Id:                   productID,
		Name:                 product.Name,
		Description:          *product.Description,
		Price:                product.Price,
		Stock:                int32(product.Stock),
		RequiresPrescription: product.RequiresPrescription,
		ImageUrl:             *product.ImageURL,
	}, nil
}

func (h *productHandler) GetProduct(ctx context.Context, req *proto.GetProductRequest) (*proto.GetProductResponse, error) {
	product, err := h.ProductService.GetProduct(req.ProductId)
	if err != nil {
		return nil, err
	}

	return &proto.GetProductResponse{
		Product: &proto.Product{
			Id:                   product.ID.String(),
			Name:                 product.Name,
			Description:          *product.Description,
			Price:                product.Price,
			Stock:                int32(product.Stock),
			RequiresPrescription: product.RequiresPrescription,
			ImageUrl:             *product.ImageURL,
		},
	}, nil
}

func (h *productHandler) ListProducts(ctx context.Context, req *proto.ListProductsRequest) (*proto.ListProductsResponse, error) {
	products, err := h.ProductService.ListProducts(req.Page, req.Limit, req.SortBy, req.SortOrder, req.Filter, req.FilterValue)
	if err != nil {
		return nil, err
	}

	var pbProducts []*proto.Product
	for _, product := range products {
		pbProducts = append(pbProducts, &proto.Product{
			Id:                   product.ID.String(),
			Name:                 product.Name,
			Description:          *product.Description,
			Price:                product.Price,
			Stock:                int32(product.Stock),
			RequiresPrescription: product.RequiresPrescription,
			ImageUrl:             *product.ImageURL,
		})
	}

	return &proto.ListProductsResponse{Products: pbProducts}, nil
}

func (h *productHandler) UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest) (*proto.UpdateProductResponse, error) {
	err := h.ProductService.UpdateProduct(req.Product.Id, req.Product.Name, req.Product.Description, req.Product.Price, req.Product.ImageUrl)
	if err != nil {
		return nil, err
	}

	return &proto.UpdateProductResponse{Message: "Product updated successfully"}, nil
}

func (h *productHandler) DeleteProduct(ctx context.Context, req *proto.DeleteProductRequest) (*proto.DeleteProductResponse, error) {
	err := h.ProductService.DeleteProduct(req.ProductId)
	if err != nil {
		return nil, err
	}

	return &proto.DeleteProductResponse{Message: "Product deleted successfully"}, nil
}

func (h *productHandler) UpdateStock(ctx context.Context, req *proto.UpdateStockRequest) (*proto.UpdateStockResponse, error) {
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, fmt.Errorf("Invalid product ID")
	}

	log := &models.InventoryLog{
		ProductID:  productId,
		Quantity:   int(req.Quantity),
		ChangeType: req.Reason,
	}

	err = h.ProductService.UpdateStock(log)
	if err != nil {
		return nil, err
	}

	return &proto.UpdateStockResponse{Message: "Stock updated successfully"}, nil
}
