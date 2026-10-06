package service

import (
	"mooma-api/config"
	"mooma-api/model"
	"mooma-api/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (s *AuthService) Login(request model.LoginRequest) (string, error) {
	var accessToken string
	customer, err := s.authRepo.FindByEmail(request.Email)

	if err != nil {
		return accessToken, err
	}

	if err := comparePassword(customer.PasswordHash, request.Password); err != nil {
		return accessToken, err
	}

	accessToken, err = generateToken(customer)
	if err != nil {
		return accessToken, err
	}

	return accessToken, nil
}

func comparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func hashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passwordHash), err
}

func generateToken(customer model.Customer) (string, error) {
	key := config.Load().JWTSecret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"customer_id": customer.ID,
			"exp":         time.Now().Add(24 * time.Hour).Unix(),
		})
	return token.SignedString([]byte(key))
}
