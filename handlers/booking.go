package handlers

import (
	"errors"
	"mooma-api/model"
	"mooma-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	service *service.BookingService
}

func NewBookingHandler(service *service.BookingService) *BookingHandler {
	return &BookingHandler{
		service: service,
	}
}

// currentCustomerID reads the ID set by AuthMiddleware, answering 401 if it is missing.
func currentCustomerID(ctx *gin.Context) (uint, bool) {
	value, exists := ctx.Get("customer_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return 0, false
	}

	id, ok := value.(uint)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid user ID",
		})
		return 0, false
	}
	return id, true
}

func (h *BookingHandler) GetMyBookings(ctx *gin.Context) {
	customerID, ok := currentCustomerID(ctx)
	if !ok {
		return
	}
	bookings, err := h.service.GetCustomerBookings(customerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, model.SuccessResponse[[]model.GetBookingResponse]{
		Status: "success",
		Data:   bookings,
	})
}

func (h *BookingHandler) CancelBooking(ctx *gin.Context) {
	customerID, ok := currentCustomerID(ctx)
	if !ok {
		return
	}
	id := ctx.Param("id")
	if err := h.service.CancelBooking(id, customerID); err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, service.ErrBookingNotFound):
			status = http.StatusNotFound
		case errors.Is(err, service.ErrCancellationClosed):
			status = http.StatusConflict
		}
		ctx.JSON(status, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "booking cancelled",
	})
}

func (h *BookingHandler) CreateBooking(ctx *gin.Context) {
	customerID, ok := currentCustomerID(ctx)
	if !ok {
		return
	}
	var request model.CreateBookingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	response, err := h.service.CreateBooking(request, customerID)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, service.ErrNoSeats):
			status = http.StatusBadRequest
		case errors.Is(err, service.ErrSeatAlreadyBooked), errors.Is(err, service.ErrBookingClosed):
			status = http.StatusConflict
		case errors.Is(err, service.ErrShowtimeNotFound):
			status = http.StatusNotFound
		}
		ctx.JSON(status, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": response,
	})
}
