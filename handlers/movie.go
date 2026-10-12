package handlers

import (
	"mooma-api/model"
	"mooma-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	service *service.MovieService
}

func NewMovieHandler(service *service.MovieService) *MovieHandler {
	return &MovieHandler{
		service: service,
	}
}

func (h *MovieHandler) CreateMovie(ctx *gin.Context) {
	var request model.CreateMovieRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	movie, err := h.service.CreateMovie(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, model.SuccessResponse[model.MovieResponse]{
		Status: "success",
		Data:   movie,
	})
}

func (h *MovieHandler) GetAllMovie(ctx *gin.Context) {
	// ?now_showing=true: only movies with a showtime still open for booking.
	movies, err := h.service.GetAllMovie(ctx.Query("now_showing") == "true")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, model.SuccessResponse[[]model.MovieResponse]{
		Status: "success",
		Data:   movies,
	})
}

func (h *MovieHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	movie, err := h.service.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, model.SuccessResponse[model.MovieResponse]{
		Status: "success",
		Data:   movie,
	})
}

func (h *MovieHandler) GetByIDWithShowtimes(ctx *gin.Context) {
	id := ctx.Param("id")
	movie, err := h.service.GetByIDWithShowtimes(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, model.SuccessResponse[model.GetMovieShowtimesResponse]{
		Status: "success",
		Data:   movie,
	})
}
