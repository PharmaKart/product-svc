# Variables
PROJECT_NAME = product-svc
GATEWAY_NAME = gateway-svc
ORDER_SERVICE_NAME = order-svc
GO = go
PROTO_DIR = internal/proto
PROTO_FILE = $(PROTO_DIR)/product.proto
PROTO_OUT = $(PROTO_DIR)
PORT = 50052

# Targets
.PHONY: build run proto clean

# Build the service
build:
	@echo "Building $(PROJECT_NAME)..."
	$(GO) build -o bin/$(PROJECT_NAME) ./cmd/main.go

# Run the service
run: build
	@echo "Running $(PROJECT_NAME) on port $(PORT)..."
	./bin/$(PROJECT_NAME)

# Generate Go code from .proto file
proto:
	@echo "Generating Go code from $(PROTO_FILE)..."
	protoc --go_out=$(PROTO_OUT) --go-grpc_out=$(PROTO_OUT) $(PROTO_FILE)
	cp $(PROTO_DIR)/product.pb.go ../$(GATEWAY_NAME)/internal/proto/product.pb.go
	cp $(PROTO_DIR)/product_grpc.pb.go ../$(GATEWAY_NAME)/internal/proto/product_grpc.pb.go
	cp $(PROTO_DIR)/product.pb.go ../$(ORDER_SERVICE_NAME)/internal/proto/product.pb.go
	cp $(PROTO_DIR)/product_grpc.pb.go ../$(ORDER_SERVICE_NAME)/internal/proto/product_grpc.pb.go

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin/$(PROJECT_NAME)