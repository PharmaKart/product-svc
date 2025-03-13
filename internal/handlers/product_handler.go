package handlers

import (
	"context"
	"fmt"

	"github.com/PharmaKart/product-svc/internal/models"
	"github.com/PharmaKart/product-svc/internal/proto"
	"github.com/PharmaKart/product-svc/internal/repositories"
	"github.com/PharmaKart/product-svc/internal/services"
	"github.com/PharmaKart/product-svc/pkg/errors"
	"github.com/PharmaKart/product-svc/pkg/utils"
	"github.com/google/uuid"
)

type ProductHandler interface {
	CreateProduct(ctx context.Context, req *proto.CreateProductRequest) (*proto.CreateProductResponse, error)
	GetProduct(ctx context.Context, req *proto.GetProductRequest) (*proto.GetProductResponse, error)
	ListProducts(ctx context.Context, req *proto.ListProductsRequest) (*proto.ListProductsResponse, error)
	UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest) (*proto.UpdateProductResponse, error)
	DeleteProduct(ctx context.Context, req *proto.DeleteProductRequest) (*proto.DeleteProductResponse, error)
	UpdateStock(ctx context.Context, req *proto.UpdateStockRequest) (*proto.UpdateStockResponse, error)
	GetInventoryLogs(ctx context.Context, req *proto.GetInventoryLogsRequest) (*proto.GetInventoryLogsResponse, error)
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
		if appErr, ok := errors.IsAppError(err); ok {
			return &proto.CreateProductResponse{
				Success: false,
				Error: &proto.Error{
					Type:    string(appErr.Type),
					Message: appErr.Message,
					Details: utils.ConvertMapToKeyValuePairs(appErr.Details),
				},
			}, nil
		}
		return &proto.CreateProductResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.InternalError),
				Message: "An unexpected error occurred",
			},
		}, nil
	}

	return &proto.CreateProductResponse{
		Success:              true,
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
		if appErr, ok := errors.IsAppError(err); ok {
			return &proto.GetProductResponse{
				Success: false,
				Error: &proto.Error{
					Type:    string(appErr.Type),
					Message: appErr.Message,
					Details: utils.ConvertMapToKeyValuePairs(appErr.Details),
				},
			}, nil
		}
		return &proto.GetProductResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.InternalError),
				Message: "An unexpected error occurred",
			},
		}, nil
	}

	return &proto.GetProductResponse{
		Success: true,
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
	var filter models.Filter
	if req.Filter != nil {
		filter = models.Filter{
			Column:   req.Filter.Column,
			Operator: req.Filter.Operator,
			Value:    req.Filter.Value,
		}
	}
	products, total, err := h.ProductService.ListProducts(filter, req.SortBy, req.SortOrder, req.Page, req.Limit)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			return &proto.ListProductsResponse{
				Success: false,
				Error: &proto.Error{
					Type:    string(appErr.Type),
					Message: appErr.Message,
					Details: utils.ConvertMapToKeyValuePairs(appErr.Details),
				},
			}, nil
		}
		return &proto.ListProductsResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.InternalError),
				Message: "An unexpected error occurred",
			},
		}, nil
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

	return &proto.ListProductsResponse{
		Success:  true,
		Products: pbProducts,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}

func (h *productHandler) UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest) (*proto.UpdateProductResponse, error) {
	err := h.ProductService.UpdateProduct(req.ProductId, req.Product.Name, req.Product.Description, req.Product.Price, req.Product.ImageUrl)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			return &proto.UpdateProductResponse{
				Success: false,
				Error: &proto.Error{
					Type:    string(appErr.Type),
					Message: appErr.Message,
					Details: utils.ConvertMapToKeyValuePairs(appErr.Details),
				},
			}, nil
		}
		return &proto.UpdateProductResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.InternalError),
				Message: "An unexpected error occurred",
			},
		}, nil
	}

	return &proto.UpdateProductResponse{
		Success: true,
		Message: "Product updated successfully",
	}, nil
}

func (h *productHandler) DeleteProduct(ctx context.Context, req *proto.DeleteProductRequest) (*proto.DeleteProductResponse, error) {
	err := h.ProductService.DeleteProduct(req.ProductId)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			return &proto.DeleteProductResponse{
				Success: false,
				Error: &proto.Error{
					Type:    string(appErr.Type),
					Message: appErr.Message,
					Details: utils.ConvertMapToKeyValuePairs(appErr.Details),
				},
			}, nil
		}
		return &proto.DeleteProductResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.InternalError),
				Message: "An unexpected error occurred",
			},
		}, nil
	}

	return &proto.DeleteProductResponse{
		Success: true,
		Message: "Product deleted successfully",
	}, nil
}

func (h *productHandler) UpdateStock(ctx context.Context, req *proto.UpdateStockRequest) (*proto.UpdateStockResponse, error) {
	productId, err := uuid.Parse(req.ProductId)
	if err != nil {
		return &proto.UpdateStockResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.ValidationError),
				Message: "Invalid product ID",
				Details: utils.ConvertMapToKeyValuePairs(map[string]string{"productId": fmt.Sprintf("Invalid UUID: %s", req.ProductId)}),
			},
		}, nil
	}

	log := &models.InventoryLog{
		ProductID:      productId,
		QuantityChange: int(req.QuantityChange),
		ChangeType:     req.Reason,
	}

	err = h.ProductService.UpdateStock(log)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			return &proto.UpdateStockResponse{
				Success: false,
				Error: &proto.Error{
					Type:    string(appErr.Type),
					Message: appErr.Message,
					Details: utils.ConvertMapToKeyValuePairs(appErr.Details),
				},
			}, nil
		}
		return &proto.UpdateStockResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.InternalError),
				Message: "An unexpected error occurred",
			},
		}, nil
	}

	return &proto.UpdateStockResponse{
		Success: true,
		Message: "Stock updated successfully",
	}, nil
}

func (h *productHandler) GetInventoryLogs(ctx context.Context, req *proto.GetInventoryLogsRequest) (*proto.GetInventoryLogsResponse, error) {
	var filter models.Filter
	if req.Filter != nil {
		filter = models.Filter{
			Column:   req.Filter.Column,
			Operator: req.Filter.Operator,
			Value:    req.Filter.Value,
		}
	}

	logs, total, err := h.ProductService.GetInventoryLogs(req.ProductId, filter, req.SortBy, req.SortOrder, req.Page, req.Limit)
	if err != nil {
		if appErr, ok := errors.IsAppError(err); ok {
			return &proto.GetInventoryLogsResponse{
				Success: false,
				Error: &proto.Error{
					Type:    string(appErr.Type),
					Message: appErr.Message,
					Details: utils.ConvertMapToKeyValuePairs(appErr.Details),
				},
			}, nil
		}
		return &proto.GetInventoryLogsResponse{
			Success: false,
			Error: &proto.Error{
				Type:    string(errors.InternalError),
				Message: "An unexpected error occurred",
			},
		}, nil
	}

	var pbLogs []*proto.InventoryLog
	for _, log := range logs {
		pbLogs = append(pbLogs, &proto.InventoryLog{
			Id:             log.ID.String(),
			ProductId:      log.ProductID.String(),
			QuantityChange: int32(log.QuantityChange),
			ChangeType:     log.ChangeType,
			CreatedAt:      log.CreatedAt.String(),
		})
	}

	return &proto.GetInventoryLogsResponse{
		Success: true,
		Logs:    pbLogs,
		Total:   total,
		Page:    req.Page,
		Limit:   req.Limit,
	}, nil
}
