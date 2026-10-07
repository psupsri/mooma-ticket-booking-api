package handlers

import (
	"mooma-api/model"
	"mooma-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ShowtimeHandler struct {
	service *service.ShowtimeService
}

func NewShowtimeHandler(service *service.ShowtimeService) *ShowtimeHandler {
	return &ShowtimeHandler{
		service: service,
	}
}

func (h *ShowtimeHandler) CreateShowtime(ctx *gin.Context) {
	var request model.CreateShowtimeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	showtime, err := h.service.CreateShowtime(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, model.SuccessResponse[model.ShowtimeResponse]{
		Status: "success",
		Data:   showtime,
	})
}

func (h *ShowtimeHandler) GetAllShowtime(ctx *gin.Context) {
	showtimes, err := h.service.GetAllShowtime()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, model.SuccessResponse[[]model.ShowtimeResponse]{
		Status: "success",
		Data:   showtimes,
	})
}
