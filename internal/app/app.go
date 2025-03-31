package app

import (
	"context"
	"go-universal-clipboard/internal/app/adapter/controller"
	"go-universal-clipboard/internal/app/adapter/route"
	"go-universal-clipboard/internal/app/infrastructure/gin"
	"go-universal-clipboard/internal/cfg"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"go.uber.org/fx"
)

type App struct {
	Cfg    *cfg.Config
	ctx    context.Context
	module fx.Option
}

func NewApp(cfg *cfg.Config) *App {
	return &App{
		Cfg: cfg,
		ctx: context.Background(),
		module: fx.Options(
			controller.Module,
			route.Module,
			gin.Module,
		),
	}
}

func (a *App) Run() {
	app := fx.New(a.module, fx.Options(
		fx.Invoke(func(
			routes route.Routes,
			requestHandler gin.RequestHandler,
		) {
			requestHandler.Gin.Use(cors.New(cors.Config{
				AllowOrigins:     []string{"*"},
				AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "api-key"},
				ExposeHeaders:    []string{"Content-Length"},
				AllowCredentials: false,
				MaxAge:           12 * time.Hour,
			}))

			routes.Setup()
			if err := requestHandler.Gin.Run(":" + strconv.Itoa(a.Cfg.App.Port)); err != nil {
				panic(err)
			}
		}),
	))
	if err := app.Start(a.ctx); err != nil {
		panic(err)
	}
}
