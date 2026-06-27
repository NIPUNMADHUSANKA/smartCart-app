# SmartCart API

Backend REST API for the SmartCart grocery shopping assistant, built with **Go** and **Fiber**.  
Includes AI-powered shopping list generation using the OpenAI API.

## Features

- User registration and login with JWT access tokens
- Password reset and account deletion
- Category CRUD with priority and status management
- Shopping item CRUD with category filtering and partial updates
- AI shopping list generation from natural language prompts
- Edit, confirm, and regenerate AI-suggested shopping lists
- JWT authentication middleware on all protected routes
- Rate limiting and security headers (Helmet)
- Panic recovery middleware
- CORS configured for web frontend

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| Web Framework | Fiber v2 |
| Database | PostgreSQL (pgx v5 driver) |
| Authentication | JWT (golang-jwt/jwt v5) |
| Validation | go-playground/validator v10 |
| AI | OpenAI API (GPT-4.1 Mini) |
| Environment | godotenv |
| IDs | Google UUID v4 |

## Project Structure

```
smartCart-app/
├── controllers/
│   ├── user_controller.go          # Register, login, profile, reset password
│   ├── category_controller.go      # Category CRUD
│   ├── shopping_item_controller.go # Shopping item CRUD
│   └── ai_model_controller.go      # AI suggestion generation & management
├── database/
│   └── database_connection.go      # PostgreSQL connection pool
├── middleware/
│   ├── auth_middleware.go          # JWT validation
│   └── rate_limiting.go            # Request rate limiting
├── models/
│   ├── user_model.go
│   ├── category_model.go           # Category, ShoppingItem, enums
│   └── ai_model.go                 # AiSuggestion, AICategory, AIItem
├── routes/
│   ├── protected_routes.go         # JWT-guarded routes
│   └── unprotected_routes.go       # Public auth routes
├── utils/
│   ├── token_util.go               # JWT generation & validation
│   └── generate_ai_util.go         # OpenAI prompt generation & DB save
├── go.mod
├── go.sum
└── main.go
```

## Requirements

- Go 1.21+
- PostgreSQL
- OpenAI API key

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/NIPUNMADHUSANKA/smartCart-app.git
cd smartCart-app
```

### 2. Create a `.env` file

```env
DATABASE_URL=postgres://user:password@localhost:5432/smartcart
MAX_CONNS=10
MIN_CONNS=2
SECRET_KEY=your_jwt_secret
SECRET_REFRESH_KEY=your_refresh_secret
OPENAI_API_KEY=your_openai_api_key
```

### 3. Install dependencies

```bash
go mod download
```

### 4. Run the server

```bash
go run .
```

The server starts on **http://localhost:8080**.

---

## API Reference

Base path: `/api/smart-cart/`

### Auth (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `auth/register` | Register a new user |
| POST | `auth/login` | Login and receive JWT |

### Auth (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `auth/me` | Get current user info |
| GET | `auth/info` | Get user details |
| PATCH | `auth/resetPassword` | Reset password |
| DELETE | `auth/remove` | Delete account |

### Categories (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `category` | Get all categories |
| POST | `category` | Create a category |
| GET | `category/:categoryId` | Get category by ID |
| DELETE | `category/:categoryId` | Delete category |
| PATCH | `category/:categoryId` | Update category |

### Shopping Items (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `shopping-item` | Get all shopping items |
| POST | `shopping-item` | Create a shopping item |
| GET | `shopping-item/:itemId` | Get item by ID |
| GET | `shopping-item/findByCategory/:categoryId` | Get items by category |
| DELETE | `shopping-item/:itemId` | Delete item |
| PATCH | `shopping-item/:itemId` | Update item |

### AI Model (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `ai-model` | Generate AI shopping list from prompt |
| GET | `ai-model` | Get all AI suggestions |
| POST | `ai-model/confirmAIShopping` | Confirm AI category → save as real items |
| POST | `ai-model/regenerateAIShopping/:suggestionId` | Delete old suggestion and regenerate |
| POST | `ai-model/addAIShoppingItem` | Add item to AI suggestion |
| PATCH | `ai-model/updateAIShoppingItem` | Update an AI item |
| DELETE | `ai-model/:categoryId` | Delete an AI category |
| DELETE | `ai-model/deleteAISuggestion/:suggestionId` | Delete full AI suggestion |
| DELETE | `ai-model/deleteAIShoppingItem/:categoryId/:itemId` | Delete a single AI item |

---

## Environment Variables

| Variable | Description |
|---|---|
| `DATABASE_URL` | PostgreSQL connection string |
| `MAX_CONNS` | Max DB connection pool size |
| `MIN_CONNS` | Min DB connection pool size |
| `SECRET_KEY` | JWT access token secret |
| `SECRET_REFRESH_KEY` | JWT refresh token secret |
| `OPENAI_API_KEY` | OpenAI API key for AI features |
