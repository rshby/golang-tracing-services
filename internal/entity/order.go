package entity

import (
	"context"
	"github.com/gin-gonic/gin"
)

type OrderService interface {
	CreateOrder(ctx context.Context, request CreateOrderRequestDTO) error
}

type OrderController interface {
	Create(c *gin.Context)
}

type CreateOrderRequestDTO struct {
	Email string `json:"email"`
}
