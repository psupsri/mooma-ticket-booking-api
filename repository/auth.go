package repository

import (
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
