package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"golang-tracing-services/internal/entity"
	"golang-tracing-services/internal/http/httpservices/customer"
	customerDto "golang-tracing-services/internal/http/httpservices/customer/dto"
	"golang-tracing-services/internal/http/httpservices/product"
	otel "golang-tracing-services/tracing"
	"sync"
)

type orderService struct {
	customerHttpService customer.CustomerHttpService
	productHttpService  product.ProductHttpService
}

// NewOrderService is function to create new instance orderService. it implements from interface OrderService
func NewOrderService(
	customerHttpService customer.CustomerHttpService,
	productHttpService product.ProductHttpService) entity.OrderService {
	return &orderService{
		customerHttpService: customerHttpService,
		productHttpService:  productHttpService,
	}
}

func (o *orderService) CreateOrder(ctx context.Context, request entity.CreateOrderRequestDTO) error {
	ctx, span := otel.Start(ctx)
	defer span.End()

	logger := logrus.WithContext(ctx)

	var (
		wg               = &sync.WaitGroup{}
		chanErr          = make(chan error, 2)
		existingCustomer *customerDto.DetailCustomerDTO
	)

	// spawn goroutine : hit customer
	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		var span trace.Span
		ctx, span = otel.Start(ctx, "", "hit customer")
		defer span.End()

		var err error
		existingCustomer, err = o.customerHttpService.GetCustomerByEmail(ctx, request.Email)
		if err != nil {
			logger.Error(err)
			span.RecordError(err)
			chanErr <- err
		}
	}(ctx, wg)

	// spawn goroutine : hit product
	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		var span trace.Span
		var err error
		ctx, span = otel.Start(ctx, "", "hit product")
		defer span.End()

		_, err = o.productHttpService.GetProductByID(ctx, 8)
		if err != nil {
			logger.Error(err)
			span.RecordError(err)
			chanErr <- err
			return
		}
	}(ctx, wg)

	wg.Wait()
	close(chanErr)
	for err := range chanErr {
		if err != nil {
			logger.Error(err)
			span.RecordError(err)
			return err
		}
	}

	logger.Infof("success create order : %v", existingCustomer)
	return nil
}
