package handlers

import (
	"mooma-api/model"
	"mooma-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var request model.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "invalid request",
		})
		return
	}
	customer, err := h.authService.Register(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "registration failed",
		})
		return
	}

	ctx.JSON(http.StatusCreated, model.SuccessResponse[model.CustomerResponse]{
		Status: "success",
		Data: model.CustomerResponse{
			FirstName: customer.FirstName,
			LastName:  customer.LastName,
			Tel:       customer.Tel,
			Email:     customer.Email,
		},
	})
}
