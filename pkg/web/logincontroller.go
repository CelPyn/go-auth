package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pynenborg.com/go-auth/pkg/domain"
)

type LoginController struct {
	userService domain.UserService
}

type LoginRequest struct {
	Name string `json:"name"`
	Password string `json:"password"`
}

func NewLoginController(service domain.UserService) domain.Controller {
	return &LoginController{userService: service}
}

func (l LoginController) Setup(r *gin.Engine) {
	r.POST("/auth/login", l.login)
}

func (l LoginController) login(c *gin.Context) {
	var loginRequest LoginRequest

	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
	}

	jwt, httpError := l.userService.Login(loginRequest.Name, loginRequest.Password)

	if httpError != nil {
		c.JSON(httpError.Status, gin.H{"message": httpError.Message, "cause": httpError.Cause})
	} else {
		c.JSON(http.StatusOK, gin.H{"token": jwt})
	}
}
