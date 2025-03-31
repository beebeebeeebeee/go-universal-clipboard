package route

import (
	"go-universal-clipboard/internal/app/adapter/controller"
	"go-universal-clipboard/internal/app/infrastructure/gin"
)

type AppRoute struct {
	handler             gin.RequestHandler
	staticController    controller.StaticController
	webSocketController *controller.WebSocketController
}

func NewAppRoute(
	handler gin.RequestHandler,
	staticController controller.StaticController,
	webSocketController *controller.WebSocketController,
) AppRoute {
	return AppRoute{
		handler:             handler,
		staticController:    staticController,
		webSocketController: webSocketController,
	}
}

func (ar *AppRoute) Setup() {
	app := ar.handler.Gin.Group("/")
	{
		ar.staticController.Setup(app.Group(""))
		ar.webSocketController.Setup(app.Group("/ws"))
	}
}
