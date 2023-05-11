package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/common/config"
	ghandler "app/common/gstuff/handler"
	"app/web/route"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	cfg = config.GetConfig()
)

// AppSrv : API Server
type AppSrv struct{}

// NewAppSrv : Tạo mới đối tượng API server
func NewAppSrv() *AppSrv {
	return &AppSrv{}
}

// Start ..
func (AppSrv) Start() (err error) {
	// TODO: init adapter: mongo, redis
	// mongo
	cfg.Mongo.Get("core").Init()

	// Echo instance
	e := echo.New()
	e.Validator = ghandler.NewValidator()
	e.HTTPErrorHandler = ghandler.Error

	// Middlewares
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	// Routes => handler
	route.App(e)

	// Start server
	go func() {
		if err := e.Start(":" + cfg.Port["app"]); err != nil {
			log.Println("⇛ shutting down the server")
			log.Println(fmt.Sprintf("%v\n", err.Error()))
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	return nil
}
