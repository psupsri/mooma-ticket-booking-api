package model

type Customer struct {
	ID           uint `gorm:"primaryKey"`
	FirstName    string
	LastName     string
	Tel          string
	Email        string `gorm:"unique"`
	PasswordHash string
}

type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Tel       string `json:"tel"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}

type CustomerResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Tel       string `json:"tel"`
	Email     string `json:"email"`
}
