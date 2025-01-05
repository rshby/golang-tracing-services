package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"golang-tracing-services/internal/config"
	"golang-tracing-services/internal/http/httpservices/customer/dto"
	otel "golang-tracing-services/tracing"
	"io"
	"net/http"
)

// CustomerHttpService is an interface
type CustomerHttpService interface {
	GetCustomerByEmail(ctx context.Context, email string) (*dto.DetailCustomerDTO, error)
}

type customerHttpService struct {
	httpClient *http.Client
}

// NewCustomerHttpService is function to create new instance customerHttpService. it implements from interface CustomerHttpService
func NewCustomerHttpService(httpClient *http.Client) CustomerHttpService {
	return &customerHttpService{
		httpClient: httpClient,
	}
}

func (c *customerHttpService) GetCustomerByEmail(ctx context.Context, email string) (*dto.DetailCustomerDTO, error) {
	ctx, span := otel.Start(ctx)
	defer span.End()

	// create httprequest
	url := spew.Sprintf("%s/v2/membership/customer/%s", config.ApiCustomerUrl(), email)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// add header request
	req.Header.Add("Device-Id", uuid.New().String())
	req.Header.Add("source", config.ApiCustomerSource())
	req.Header.Add("Authorization", config.ApiCustomerBasicAuth())

	// execute
	res, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	defer func() {
		if err = res.Body.Close(); err != nil {
			span.RecordError(err)
			return
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// switch status response
	switch res.StatusCode {
	case http.StatusOK:
		var responseObject dto.ResourceDetailCustomerDTO
		if err = json.Unmarshal(body, &responseObject); err != nil {
			span.RecordError(err)
			return nil, err
		}

		return &responseObject.Data, nil
	default:
		err = fmt.Errorf("get response %d %s [%s]", res.StatusCode, http.StatusText(res.StatusCode), string(body))
		span.RecordError(err)
		return nil, err
	}
}
