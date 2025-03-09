package utils

import (
	"regexp"
	"strings"

	"github.com/PharmaKart/product-svc/internal/models"
	"github.com/PharmaKart/product-svc/pkg/errors"
)

func ValidateProductInput(product *models.Product) error {
	validationErrors := make(map[string]string)
	if strings.TrimSpace(product.Name) == "" {
		validationErrors["name"] = "Name is required"
	}

	if strings.TrimSpace(*product.Description) == "" {
		validationErrors["description"] = "Description is required"
	}

	if product.Price <= 0 {
		validationErrors["price"] = "Price must be greater than 0"
	}

	if product.Stock < 0 {
		validationErrors["stock"] = "Stock must be greater than or equal to 0"
	}

	if product.ImageURL != nil {
		trimmedURL := strings.TrimSpace(*product.ImageURL)
		s3Pattern := `^https://[^.]+\.s3\.[^.]+\.amazonaws\.com/`
		matched, err := regexp.MatchString(s3Pattern, trimmedURL)
		if err != nil || !matched {
			validationErrors["imageURL"] = "Invalid S3 image URL"
		}
	}

	return nil
}

func ValidateInventoryInput(inventory *models.InventoryLog) error {

	if inventory.ChangeType != "order_placed" && inventory.ChangeType != "order_cancelled" && inventory.ChangeType != "stock_added" {
		return errors.NewValidationError("changeType", "Invalid change type")
	}

	return nil
}
