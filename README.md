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

```
Table Category {
  id uint [pk]
  name string [not null]
}

Table Product {
  id uint [pk]
  name string [not null]
  description string [not null]
  price float [not null]
  category_id uint [not null]
  interactions uint
  weighted_score float
}

Table Interaction {
  id uint [pk]
  user_id uint
  product_id uint
  category_id uint
}

Table User {
  id uint [pk]
  name string [not null]
  email string [not null]
  password string [not null]
  role string [not null]
}

Ref: Product.category_id > Category.id
Ref: Interaction.user_id > User.id
Ref: Interaction.product_id > Product.id
Ref: Interaction.category_id > Category.id
Ref: User.id < interactions.user_id
Ref: Product.id < interactions.product_id
``` 

## Technologies Used

- Golang
- Fiber
- GORM
- PostgreSQL
