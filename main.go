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
	"gorm.io/gorm"
)

func newServer(lc fx.Lifecycle, db *gorm.DB) *gin.Engine {
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

			if err := utils.AutoMigrate(db); err != nil {
				log.Println("Error migrating database:", err)
				return nil
			}

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
				log.Println("Error shutting down server:", err)
				return err
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
		fx.Provide(utils.NewDB),
		fx.Provide(apis.UsersAPI),
		fx.Provide(newServer),
		fx.Invoke(
			apis.UserRoutes),
	)

	app.Run()
}
