package router

import (
	"github.com/gin-gonic/gin"
	"golang-tracing-services/internal/entity"
	"golang-tracing-services/internal/http/controller"
	"golang-tracing-services/internal/http/httpservices/customer"
	"golang-tracing-services/internal/http/httpservices/product"
	"golang-tracing-services/internal/service"
	"net/http"
)

type Router struct {
	app             *gin.RouterGroup
	httpClient      *http.Client
	orderController entity.OrderController
}

func NewRouter(app *gin.RouterGroup, httpClient *http.Client) {
	r := Router{
		app:        app,
		httpClient: httpClient,
	}

	r.Register()
	r.ApiV1(r.app)
}

func (r *Router) ApiV1(app *gin.RouterGroup) {
	v1Group := app.Group("/v1")
	{
		orderV1Group := v1Group.Group("/order")
		{
			orderV1Group.POST("", r.orderController.Create)
		}
	}
}

func (r *Router) Register() {
	customerHttpService := customer.NewCustomerHttpService(r.httpClient)
	productHttpService := product.NewProductHttpService(r.httpClient)

	// create instance service
	orderService := service.NewOrderService(customerHttpService, productHttpService)

	// create instance controller
	orderController := controller.NewOrderController(orderService)

	r.orderController = orderController
}
