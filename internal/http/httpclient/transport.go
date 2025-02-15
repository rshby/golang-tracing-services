package httpclient

import (
	"bytes"
	"github.com/sirupsen/logrus"
	ioOtel "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"golang-tracing-services/internal/utils/helper"
	otel "golang-tracing-services/tracing"
	"io"
	"moul.io/http2curl"
	"net/http"
)

type transport struct {
	base http.RoundTripper
}

func NewTransport() http.RoundTripper {
	return &transport{
		base: http.DefaultTransport,
	}
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	ctx, span := otel.Start(r.Context(), "", "HTTP")
	defer span.End()

	logger := logrus.WithContext(ctx)

	// inject traceID to header
	ioOtel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))

	curlCommand, err := http2curl.GetCurlCommand(r)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	span.SetAttributes(attribute.String("http.request.curl", curlCommand.String()))

	if r.Body != nil {
		// get buffer from Pool. don't forget to put back to Pool
		buf := helper.BytesBufferPool.Get().(*bytes.Buffer)
		defer func() {
			buf.Reset()
			helper.BytesBufferPool.Put(buf)
		}()

		// copy from request body to buffer. request body will be empty after this
		if _, err = io.Copy(buf, r.Body); err != nil {
			logger.Error(err)
			return nil, err
		}

		// add span attribute
		span.SetAttributes(attribute.String("http.request.body", buf.String()))

		// reassign to request body
		_ = r.Body.Close()
		r.Body = io.NopCloser(bytes.NewReader(nil))
	}

	// execute HTTP call
	r = r.WithContext(ctx)
	res, err := t.base.RoundTrip(r)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if res != nil {
		// get buffer from Pool. don't forget to put back to Pool
		buf := helper.BytesBufferPool.Get().(*bytes.Buffer)
		defer func() {
			buf.Reset()
			helper.BytesBufferPool.Put(buf)
		}()

		// copy from response body to buffer. response body will be empty after this
		if _, err = io.Copy(buf, res.Body); err != nil {
			logger.Error(err)
			return nil, err
		}

		// reassign to response body
		_ = res.Body.Close()
		res.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
		span.SetAttributes(attribute.String("http.response.body", buf.String()))
	}

	return res, nil
}
