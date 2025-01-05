package product

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-tracing-services/internal/http/httpservices/product/dto"
	otel "golang-tracing-services/tracing"
	"io"
	"net/http"
)

type ProductHttpService interface {
	GetProductByID(ctx context.Context, id int) (string, error)
}

type productHttpService struct {
	httpClient *http.Client
}

// NewProductHttpService is function to create new instance productHttpService. it implements from interface ProductHttpService
func NewProductHttpService(httpClient *http.Client) ProductHttpService {
	return &productHttpService{
		httpClient: httpClient,
	}
}

// GetProductByID is method to get product by ID
func (p *productHttpService) GetProductByID(ctx context.Context, id int) (string, error) {
	ctx, span := otel.Start(ctx)
	defer span.End()

	// create http request
	url := fmt.Sprintf("http://localhost:4002/v1/product/%d", id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	// execute api
	req = req.WithContext(ctx)
	res, err := p.httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	switch res.StatusCode {
	case http.StatusOK:
		var response dto.GetProductResponseDTO
		if err = json.Unmarshal(body, &response); err != nil {
			span.RecordError(err)
			return "", err
		}

		return response.Data, nil
	default:
		err = fmt.Errorf("get response %d %s : [%s]", res.StatusCode, http.StatusText(res.StatusCode), string(body))
		span.RecordError(err)
		return "", err
	}
}
