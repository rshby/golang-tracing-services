package httpclient

import (
	"bytes"
	ioOtel "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
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
	ctx, span := otel.Start(r.Context())
	defer span.End()

	// inject traceID to header
	ioOtel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))
	r = r.WithContext(ctx)

	res, err := t.base.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	res.Body.Close()

	res.Body = io.NopCloser(bytes.NewReader(body))

	curlCommand, err := http2curl.GetCurlCommand(r)
	if err != nil {
		return nil, err
	}

	span.AddEvent("execute API", trace.WithAttributes(
		attribute.String("http.request.curl", curlCommand.String()),
		attribute.String("http.response.body", string(body)),
	))

	return res, nil
}
