package handlers

import (
	"mooma-api/model"
	"mooma-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CinemaHandler struct {
	service *service.CinemaService
}

func NewCinemaHandler(service *service.CinemaService) *CinemaHandler {
	return &CinemaHandler{
		service: service,
	}
}

func (h *CinemaHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	cinema, err := h.service.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, model.SuccessResponse[model.CinemaResponse]{
		Status: "success",
		Data:   cinema,
	})
}
