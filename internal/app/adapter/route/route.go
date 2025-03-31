package route

import (
	"go-universal-clipboard/internal/app/adapter/controller"
	ginHandler "go-universal-clipboard/internal/app/infrastructure/gin"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewRoutes),
)

type Routes struct {
	gin          ginHandler.RequestHandler
	wsController *controller.WebSocketController
}

func NewRoutes(gin ginHandler.RequestHandler, wsController *controller.WebSocketController) *Routes {
	return &Routes{
		gin:          gin,
		wsController: wsController,
	}
}

func (r *Routes) Setup() {
	// Serve static files
	r.gin.Gin.Static("/static", "./internal/app/infrastructure/gin/static")

	// Root endpoint - redirect to a random room
	r.gin.Gin.GET("/", func(c *gin.Context) {
		roomID := generateRoomID()
		c.Redirect(http.StatusTemporaryRedirect, "/"+roomID)
	})

	// Room endpoint - serve the clipboard page
	r.gin.Gin.GET("/:roomId", func(c *gin.Context) {
		c.File("./internal/app/infrastructure/gin/static/index.html")
	})

	// WebSocket endpoint
	r.gin.Gin.GET("/ws", r.wsController.HandleWebSocket)
}

func generateRoomID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
