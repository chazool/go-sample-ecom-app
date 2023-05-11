
1. POST `/register` - This endpoint is used to register a new user.
   - Request: 
      - URL: `http://localhost:3000/register`
      - Body: `{ "name": "John", "email": "john@example.com", "password": "password" }`
   - Response:
      - Status Code: 200
      - Body: `{ "message": "user registered successfully", "token": "<jwt token>" }`

2. POST `/login` - This endpoint is used to authenticate a user and generate a JWT token.
   - Request:
      - URL: `http://localhost:3000/login`
      - Body: `{ "email": "john@example.com", "password": "password" }`
   - Response:
      - Status Code: 200
      - Body: `{ "token": "<jwt token>" }`

3. GET `/product` - This endpoint retrieves all products.
   - Request:
      - URL: `http://localhost:3000/product`
      - Headers: `{ "Authorization": "Bearer <jwt token>" }`
   - Response:
      - Status Code: 200
      - Body: `[ { "id": 1, "name": "Product 1", "description": "Description 1", "price": 10.5, "created_at": "2023-05-11T15:00:00Z", "updated_at": "2023-05-11T15:00:00Z", "category_id": 1, "category": { "id": 1, "name": "Category 1" } }, { "id": 2, "name": "Product 2", "description": "Description 2", "price": 20.5, "created_at": "2023-05-11T15:00:00Z", "updated_at": "2023-05-11T15:00:00Z", "category_id": 2, "category": { "id": 2, "name": "Category 2" } } ]`

4. POST `/product` - This endpoint creates a new product.
   - Request:
      - URL: `http://localhost:3000/product`
      - Headers: `{ "Authorization": "Bearer <jwt token>" }`
      - Body: `{ "name": "Product 3", "description": "Description 3", "price": 30.5, "category_id": 1 }`
   - Response:
      - Status Code: 200
      - Body: `{ "id": 3, "name": "Product 3", "description": "Description 3", "price": 30.5, "created_at": "2023-05-11T15:00:00Z", "updated_at": "2023-05-11T15:00:00Z", "category_id": 1, "category": { "id": 1, "name": "Category 1" } }`

5. GET `/product/:id` - This endpoint retrieves a single product by ID.
   - Request:
      - URL: `http://localhost:3000/product/1`
      - Headers: `{ "Authorization": "Bearer <jwt token>" }`
   - Response:
      - Status Code: 200
      - Body: `{ "id": 1, "name": "Product 1", "description": "Description 1", "price": 10.5, "created_at": "2023-