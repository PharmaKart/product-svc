package main

import (
	"net"

	"github.com/PharmaKart/product-svc/internal/handlers"
	pb "github.com/PharmaKart/product-svc/internal/proto"
	"github.com/PharmaKart/product-svc/internal/repositories"
	"github.com/PharmaKart/product-svc/pkg/config"
	"github.com/PharmaKart/product-svc/pkg/utils"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func main() {
	// Initialize logger
	utils.InitLogger()

	// Load configuration
	cfg := config.LoadConfig()

	//Initialize repositories
	productrepo := repositories.NewProductRepository(&gorm.DB{})
	inventorylogrepo := repositories.NewInventoryLogRepository(&gorm.DB{})

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productrepo, inventorylogrepo)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", ":"+cfg.Port)

	if err != nil {
		utils.Logger.Fatal("Failed to listen", map[string]interface{}{
			"error": err,
		})
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, productHandler)

	utils.Info("Starting product service", map[string]interface{}{
		"port": cfg.Port,
	})

	if err := grpcServer.Serve(lis); err != nil {
		utils.Logger.Fatal("Failed to serve", map[string]interface{}{
			"error": err,
		})
	}
}
