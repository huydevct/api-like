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

	ghandler "app/common/gstuff/handler"
	"app/web/route"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// AutofarmerNetAPI8 : AutofarmerNetAPI8 Server
type AutofarmerNetAPI8 struct{}

// NewAutofarmerNetAPI8 : Tạo mới đối tượng API server
func NewAutofarmerNetAPI8() *AutofarmerNetAPI8 {
	return &AutofarmerNetAPI8{}
}

// Start ..
func (AutofarmerNetAPI8) Start() (err error) {
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
	route.AutofarmerNetAPI8(e)

	// Start server
	go func() {
		if err := e.Start(":" + cfg.Port["api8_autofarmer_net"]); err != nil {
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
