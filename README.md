# Product Recommendation API

This API provides product recommendations based on user interactions. It allows users to search for products and categories, as well as create and manage products and categories.

## Installation

1. Install dependencies:
```sh
go mod tidy
```
3. create DB sample_app

2. Start the server:
```sh
go run cmd/main.go
```

The server should be running at http://localhost:8181.

## Usage

The following routes are available:

### Authentication

- POST /login
  - Log in with email and password.
- POST /auth/register
  - Register a new user.
  - Requires authentication.

### Management

- POST /manage/product
  - Create a new product.
  - Requires authentication.
- DELETE /manage/product/:id
  - Delete a product by ID.
  - Requires authentication.
- POST /manage/category
  - Create a new category.
  - Requires authentication.

### Searching

- GET /search/product/recommendation
  - Get product recommendations for the authenticated user.
  - Requires authentication.
- GET /search/product
  - Get all products.
- GET /search/product/:id
  - Get a product by ID.
- GET /search/category
  - Get all categories.
- GET /search/category/:id
  - Get a category by ID.

## Database Structure

The database structure is defined using DBML. Here is a summary:
https://dbdiagram.io/d/645e5141dca9fb07c4fbdba0

## Technologies Used

- Golang
- Fiber
- GORM
- PostgreSQL
