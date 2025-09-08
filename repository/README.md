# Repository Package

This package contains GORM models and repository implementations for the food delivery API.

## Structure

```
repository/
├── model/                    # GORM models
│   ├── user.go              # User model
│   ├── merchant.go          # Merchant model
│   ├── courier.go           # Courier model
│   ├── food.go              # Food model
│   ├── transaction.go       # Transaction model
│   └── models.go            # Model utilities
├── database.go              # Database connection and configuration
├── user_repository.go       # User repository implementation
├── merchant_repository.go   # Merchant repository implementation
├── food_repository.go       # Food repository implementation
├── courier_repository.go    # Courier repository implementation
├── transaction_repository.go # Transaction repository implementation
├── repositories.go          # Repository manager
├── example.go              # Usage examples
└── README.md               # This file
```

## Models

### User
- ID, Username, Email, Password
- Soft deletes enabled
- Relationships: Merchants, Transactions

### Merchant
- ID, Name, Address, UserID
- Soft deletes enabled
- Relationships: User, Foods, Transactions

### Courier
- ID, Name, Phone, Latitude, Longitude
- Soft deletes enabled
- Relationships: Transactions

### Food
- ID, Name, Price, MerchantID, Description
- Soft deletes enabled
- Relationships: Merchant, Transactions

### Transaction
- ID, UserID, MerchantID, CourierID, FoodID
- Quantity, TotalPrice, Status
- Soft deletes enabled
- Relationships: User, Merchant, Courier, Food

## Usage

### 1. Database Setup

```go
config := repository.DatabaseConfig{
    Host:     "localhost",
    Port:     3306,
    User:     "your_username",
    Password: "your_password",
    DBName:   "your_database",
}

db, err := repository.NewDatabase(config)
if err != nil {
    log.Fatal("Failed to connect to database:", err)
}
defer db.Close()

// Run auto-migration
if err := db.AutoMigrate(); err != nil {
    log.Fatal("Failed to migrate database:", err)
}
```

### 2. Repository Usage

```go
// Initialize repositories
repos := repository.NewRepositories(db.DB)

// Create a user
user := &model.User{
    Username: "john_doe",
    Email:    "john@example.com",
    Password: "hashed_password",
}
err := repos.User.Create(user)

// Get user by ID
user, err := repos.User.GetByID(1)

// Get user by email
user, err := repos.User.GetByEmail("john@example.com")

// Update user
user.Username = "john_updated"
err := repos.User.Update(user)

// Delete user (soft delete)
err := repos.User.Delete(1)
```

## Repository Methods

Each repository provides standard CRUD operations:

- `Create(entity)` - Create new entity
- `GetByID(id)` - Get entity by ID
- `Update(entity)` - Update existing entity
- `Delete(id)` - Soft delete entity
- `GetAll(limit, offset)` - Get paginated list

### Additional Methods

**UserRepository:**
- `GetByEmail(email)` - Get user by email
- `GetByUsername(username)` - Get user by username

**MerchantRepository:**
- `GetByUserID(userID)` - Get merchants by user ID

**FoodRepository:**
- `GetByMerchantID(merchantID)` - Get foods by merchant ID

**CourierRepository:**
- `GetNearby(lat, lng, radius)` - Get couriers near location

**TransactionRepository:**
- `GetByUserID(userID)` - Get transactions by user ID
- `GetByMerchantID(merchantID)` - Get transactions by merchant ID
- `GetByCourierID(courierID)` - Get transactions by courier ID
- `GetByStatus(status)` - Get transactions by status

## Database Configuration

Make sure to set up your MySQL database and update the configuration in your main application:

```sql
CREATE DATABASE your_database_name;
```

The models will automatically create the necessary tables when `AutoMigrate()` is called.

## Notes

- All models use soft deletes (DeletedAt field)
- Timestamps are automatically managed by GORM
- Foreign key relationships are properly defined
- Repository pattern is implemented for clean separation of concerns
- All price fields use integer type (store in smallest currency unit)
