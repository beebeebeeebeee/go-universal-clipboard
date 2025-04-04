package controller

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"
	_ "time/tzdata"
)

type StaticController struct {
	staticPath string
}

func NewStaticController() StaticController {
	return StaticController{
		staticPath: filepath.Join("internal", "app", "infrastructure", "gin", "static"),
	}
}

func (rc *StaticController) Setup(g *gin.RouterGroup) {
	rg := g.Group("")
	{
		rg.Static("/static", rc.staticPath)
		rg.GET("", rc.generateRoomID)
		rg.GET("/:roomId", rc.getRoom)
	}
}

func (rc *StaticController) generateRoomID(c *gin.Context) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	c.Redirect(http.StatusTemporaryRedirect, "/"+string(b))
}

func (rc *StaticController) getRoom(c *gin.Context) {
	roomId := c.Param("roomId")
	if roomId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room parameter is required"})
		return
	}

	c.File(filepath.Join(rc.staticPath, "index.html"))
}
