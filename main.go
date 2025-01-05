package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/codes"
	"golang-tracing-services/internal/config"
	"golang-tracing-services/internal/http/httpclient"
	"golang-tracing-services/internal/http/middleware"
	"golang-tracing-services/internal/router"
	otel "golang-tracing-services/tracing"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	shutdownTraceProvider := otel.NewTraceProvider(context.Background())
	defer shutdownTraceProvider()

	app := gin.Default()
	app.Use(middleware.TraceMiddleware())

	router.NewRouter(&app.RouterGroup, httpclient.NewHttpClient())

	app.GET("/ping", controller)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port()),
		Handler: app,
	}

	chanSignal := make(chan os.Signal, 1)
	chanErr := make(chan error, 1)
	chanQuit := make(chan struct{}, 1)

	signal.Notify(chanSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-chanSignal:
				logrus.Info("receive signal interrupt ⚠️")
				gracefullShutdown(srv)
				chanQuit <- struct{}{}
				return
			case e := <-chanErr:
				logrus.Infof("receive signal error ⚠️ : %s", e.Error())
				gracefullShutdown(srv)
				chanQuit <- struct{}{}
				return
			}
		}
	}()

	go func() {
		logrus.Infof("running http server listening on port %d ⏳", config.Port())
		if err := srv.ListenAndServe(); err != nil {
			chanErr <- err
			return
		}
	}()

	<-chanQuit
	close(chanSignal)
	close(chanErr)
	close(chanQuit)

	logrus.Info("server exit ‼️")
}

func controller(c *gin.Context) {
	ctx, span := otel.Start(c.Request.Context())
	defer span.End()

	// call service
	err := service(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func service(ctx context.Context) error {
	ctx, span := otel.Start(ctx)
	defer span.End()

	// call repository
	if err := repository(ctx); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func repository(ctx context.Context) error {
	ctx, span := otel.Start(ctx)
	defer span.End()

	var err error
	err = errors.New("error database")
	if err != nil {
		span.RecordError(err)
	}

	return err
}

func gracefullShutdown(srv *http.Server) {
	if srv == nil {
		return
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Warn("force close server ⚠️")
		_ = srv.Close()
		return
	}

	_ = srv.Close()
	logrus.Info("success gracefull shutdown http server ❎")
}
