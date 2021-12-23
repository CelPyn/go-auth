package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type PongController struct {

}

func (p PongController) Setup(r *gin.Engine) {
	r.GET("/pong", p.Pong)
}

func (p PongController) Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"value": "ping"})
}
