package repository

import (
	"errors"
	"mooma-api/model"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Create(customer *model.Customer) error {
	return r.db.Create(customer).Error
}

func (r *AuthRepository) FindByEmail(email string) (model.Customer, error) {
	var customer model.Customer
	result := r.db.Where("email = ?", email).First(&customer)
	if result.RowsAffected == 0 {
		return customer, errors.New("customer not found")
	}
	if result.Error != nil {
		return customer, result.Error
	}
	return customer, nil
}
