package utils

import (
	"errors"
	"strings"

	"github.com/PharmaKart/product-svc/internal/models"
)

func ValidateProductInput(product *models.Product) error {
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("name is required")
	}

	if strings.TrimSpace(*product.Description) == "" {
		return errors.New("description is required")
	}

	if product.Price <= 0 {
		return errors.New("price must be greater than 0")
	}

	if product.Stock < 0 {
		return errors.New("stock must be greater than or equal to 0")
	}

	// Validate image URL to be a valid S3 URL
	// TODO: Uncomment this after adding S3 upload
	// if product.ImageURL != nil {
	// 	if !strings.HasPrefix(strings.TrimSpace(*product.ImageURL), "https://s3.amazonaws.com/") {
	// 		return errors.New("Invalid image uploaded")
	// 	}
	// }

	return nil
}

func ValidateInventoryInput(inventory *models.InventoryLog) error {

	if inventory.ChangeType != "order_placed" && inventory.ChangeType != "order_cancelled" && inventory.ChangeType != "stock_added" {
		return errors.New("invalid change type")
	}

	return nil
}
