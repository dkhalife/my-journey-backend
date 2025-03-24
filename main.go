package main

import (
	"context"
	"log"
	"net/http"

	"dkhalife.com/journey/internal/apis"
	"dkhalife.com/journey/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func newServer(lc fx.Lifecycle) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(utils.RequestLogger())

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting server")
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Println("Error starting server:", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping server")
			if err := srv.Shutdown(ctx); err != nil {
				log.Fatalf("Error shutting down server: %v", err)
			}
			return nil
		},
	})

	return r
}

func main() {
	app := fx.New(
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.NopLogger
		}),
		fx.Provide(apis.UsersAPI),
		fx.Provide(newServer),
		fx.Invoke(
			apis.UserRoutes,
			func(r *gin.Engine) {}),
	)

	app.Run()
}
