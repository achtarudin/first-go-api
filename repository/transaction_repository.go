package repository

import (
	"cutbray/first_api/repository/model"

	"gorm.io/gorm"
)

// TransactionRepository interface defines methods for transaction operations
type TransactionRepository interface {
	Create(transaction *model.Transaction) error
	GetByID(id uint) (*model.Transaction, error)
	GetByUserID(userID uint) ([]model.Transaction, error)
	GetByMerchantID(merchantID uint) ([]model.Transaction, error)
	GetByCourierID(courierID uint) ([]model.Transaction, error)
	GetByStatus(status string) ([]model.Transaction, error)
	Update(transaction *model.Transaction) error
	Delete(id uint) error
	GetAll(limit, offset int) ([]model.Transaction, error)
}

// transactionRepository implements TransactionRepository interface
type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transaction *model.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) GetByID(id uint) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.db.Preload("User").Preload("Merchant").Preload("Courier").Preload("Food").First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) GetByUserID(userID uint) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Preload("Merchant").Preload("Courier").Preload("Food").Where("user_id = ?", userID).Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) GetByMerchantID(merchantID uint) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Preload("User").Preload("Courier").Preload("Food").Where("merchant_id = ?", merchantID).Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) GetByCourierID(courierID uint) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Preload("User").Preload("Merchant").Preload("Food").Where("courier_id = ?", courierID).Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) GetByStatus(status string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Preload("User").Preload("Merchant").Preload("Courier").Preload("Food").Where("status = ?", status).Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) Update(transaction *model.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *transactionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Transaction{}, id).Error
}

func (r *transactionRepository) GetAll(limit, offset int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Preload("User").Preload("Merchant").Preload("Courier").Preload("Food").Limit(limit).Offset(offset).Find(&transactions).Error
	return transactions, err
}
