package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pynenborg.com/go-auth/pkg/domain"
)

type PingController struct {
	userService domain.UserService
}

func NewPingController(service domain.UserService) domain.Controller {
	return &PingController{userService: service}
}

func (p PingController) Setup(r *gin.Engine) {
	r.GET("/ping", p.Ping)
	r.GET("/ping/:name", p.PingName)
}

func (p PingController) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"value": "pong"})
}

func (p PingController) PingName(c *gin.Context) {
	nameParam := c.Param("name")
	c.JSON(http.StatusOK, gin.H{"value": nameParam})
}
