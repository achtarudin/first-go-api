package repository

import "gorm.io/gorm"

// Repositories holds all repository instances
type Repositories struct {
	User        UserRepository
	Merchant    MerchantRepository
	Food        FoodRepository
	Courier     CourierRepository
	Transaction TransactionRepository
}

// NewRepositories creates a new instance of all repositories
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:        NewUserRepository(db),
		Merchant:    NewMerchantRepository(db),
		Food:        NewFoodRepository(db),
		Courier:     NewCourierRepository(db),
		Transaction: NewTransactionRepository(db),
	}
}
