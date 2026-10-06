package service

import (
	"mooma-api/model"
	"mooma-api/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo *repository.AuthRepository
}

func NewAuthService(authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

func (s *AuthService) Register(request model.RegisterRequest) (model.Customer, error) {
	var customer model.Customer
	hashedPassword, err := hashPassword(request.Password)
	if err != nil {
		return customer, err
	}
	customer = model.Customer{
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		Tel:          request.Tel,
		Email:        request.Email,
		PasswordHash: hashedPassword,
	}

	if err := s.authRepo.Create(&customer); err != nil {
		return customer, err
	}

	return customer, nil
}

func hashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passwordHash), err
}
