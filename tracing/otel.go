package tracing

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"runtime"
)

func newOtelGRPCExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	//return otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint("localhost:4317"))
	return otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
}

func NewTraceProvider(ctx context.Context, exp ...sdktrace.SpanExporter) func() {
	if len(exp) == 0 {
		exp = make([]sdktrace.SpanExporter, 1)
		var err error
		exp[0], err = newOtelGRPCExporter(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("go-tracing-service"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp[0]),
		sdktrace.WithResource(r),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, // Propagasi trace-id, span-id, dll
		propagation.Baggage{},      // Propagasi metadata tambahan (key-value pairs)
	))

	return func() {
		_ = tp.Shutdown(ctx)
	}
}

func Start(ctx context.Context, name ...string) (context.Context, trace.Span) {
	pc, _, _, _ := runtime.Caller(1)
	if len(name) < 2 {
		name = make([]string, 2)
		name[0] = "go-tracing-service"
		name[1] = runtime.FuncForPC(pc).Name()
	} else {
		if name[0] == "" {
			name[0] = "go-tracing-service"
		}

		if name[1] == "" {
			name[1] = runtime.FuncForPC(pc).Name()
		}
	}

	c, ok := ctx.(*gin.Context)
	if ok {
		ctx = c.Request.Context()
	}

	return otel.Tracer(name[0]).Start(ctx, name[1])
}
