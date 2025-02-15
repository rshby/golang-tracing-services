package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/codes"
	"golang-tracing-services/internal/entity"
	otel "golang-tracing-services/tracing"
	"net/http"
)

type orderController struct {
	orderService entity.OrderService
}

// NewOrderController is function to create new instance orderController. it implements from interface OrderController
func NewOrderController(orderService entity.OrderService) entity.OrderController {
	return &orderController{
		orderService: orderService,
	}
}

func (o *orderController) Create(c *gin.Context) {
	ctx, span := otel.Start(c)
	defer span.End()

	logger := logrus.WithContext(ctx)

	var request entity.CreateOrderRequestDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"traceID": ctx.Value("traceID").(string),
		})
		return
	}

	// call method in service
	if err := o.orderService.CreateOrder(ctx, request); err != nil {
		logger.Error(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"traceID": ctx.Value("traceID").(string),
		})
		return
	}

	// success
	c.JSON(http.StatusOK, gin.H{
		"message": "success create order",
		"traceID": ctx.Value("traceID").(string),
	})
}
