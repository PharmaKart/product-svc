# Product Service

The **Product Service** is a core component of the Pharmakart platform, responsible for managing product information, inventory, and related operations. It provides secure endpoints for adding, updating, deleting, and retrieving products, as well as managing inventory updates.

---

## Table of Contents
1. [Overview](#overview)
2. [Features](#features)
3. [Prerequisites](#prerequisites)
4. [Setup and Installation](#setup-and-installation)
5. [Running the Service](#running-the-service)
6. [Environment Variables](#environment-variables)
7. [Contributing](#contributing)
8. [License](#license)

---

## Overview

The Product Service handles:
- Product management (add, update, delete, retrieve).
- Inventory management (update stock levels).
- Role-based access control (admin and customer access).

It is built using **gRPC** for communication and **PostgreSQL** for data storage.

---

## Features

- **Product Management**:
  - Add, update, and delete products.
  - Retrieve product details and list products with pagination and filtering.
- **Inventory Management**:
  - Update product inventory (restocking and sales).
- **Role-Based Access Control**:
  - Admin-only access for product modifications.
  - Customer access for product retrieval.

---

## Prerequisites

Before setting up the service, ensure you have the following installed:
- **Docker**
- **Go** (for building and running the service)
- **Protobuf Compiler** (`protoc`) for generating gRPC/protobuf files

---

## Setup and Installation

### 1. Clone the Repository
Clone the repository and navigate to the product service directory:
```bash
git clone https://github.com/PharmaKart/product-svc.git
cd product-svc
```

### 2. Generate Protobuf Files
Generate the protobuf files using the provided `Makefile`:
```bash
make proto
```

### 3. Install Dependencies
Run the following command to ensure all dependencies are installed:
```bash
go mod tidy
```

### 4. Build the Service
To build the service, run:
```bash
make build
```

---

## Running the Service

### Option 1: Run Using Docker
To run the service using Docker, execute:
```bash
docker run -p 50052:50052 pharmakart/product-svc
```

### Option 2: Run Using Makefile
To run the service directly using Go, execute:
```bash
make run
```

The service will be available at:
- **gRPC**: `localhost:50052`

---

## Environment Variables

The service requires the following environment variables. Create a `.env` file in the `product-svc` directory with the following:

```env
PRODUCT_DB_HOST=postgres
PRODUCT_DB_PORT=5432
PRODUCT_DB_USER=postgres
PRODUCT_DB_PASSWORD=yourpassword
PRODUCT_DB_NAME=pharmakartdb
```

---

## Contributing

Contributions are welcome! Please follow these steps:
1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Submit a pull request with a detailed description of your changes.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Support

For any questions or issues, please open an issue in the repository or contact the maintainers.
