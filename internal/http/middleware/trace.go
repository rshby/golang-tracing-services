package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	ioOtel "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	otel "golang-tracing-services/tracing"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := ioOtel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		c.Request = c.Request.WithContext(ctx)

		// span name
		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())

		// create new span
		ctx, span := otel.Start(c.Request.Context(), "", spanName)
		defer span.End()

		traceID := span.SpanContext().TraceID().String()
		span.SetAttributes(attribute.String("traceID", traceID))
		ctx = context.WithValue(ctx, "traceID", traceID)

		w := NewResponseBodyWriter(c.Writer)
		c.Writer = w
		c.Request = c.Request.WithContext(ctx)

		// continue to next handler
		c.Next()

		// get response and set to span information
		span.SetAttributes(
			attribute.Int("http.status.code", w.ginResponseWriter.Status()),
			attribute.String("http.response.body", w.body.String()))

		w.PutBack()
	}
}
