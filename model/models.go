package model

// This file provides a central place to import all models
// and can be used for database migrations and model registration

// AllModels returns a slice of all model structs for auto-migration
func AllModels() []interface{} {
	return []interface{}{
		&User{},
		&Merchant{},
		&Courier{},
		&Food{},
		&Transaction{},
	}
}
