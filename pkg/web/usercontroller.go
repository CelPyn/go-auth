package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pynenborg.com/go-auth/pkg/domain"
	"pynenborg.com/go-auth/pkg/service"
)

type UserController struct {
	userService domain.UserService
}

func NewUserController(service domain.UserService) domain.Controller {
	return &UserController{userService: service}
}

func (u UserController) Setup(r *gin.Engine) {
	r.POST("/users", u.createUser)
	r.GET("/users/:name", u.getUser)
}

func (u UserController) createUser(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func (u UserController) getUser(c *gin.Context) {
	token := c.Request.Header.Get("authorization")
	if len(token) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to access this resource"})
		return
	}

	validationErr := service.ValidateBearer(token)
	if validationErr != nil {
		c.JSON(validationErr.Status, gin.H{"message": validationErr.Message})
		return
	}

	name := c.Param("name")
	user, httpError := u.userService.Get(name)

	if httpError != nil {
		c.JSON(httpError.Status, gin.H{"message": httpError.Message, "cause": httpError.Cause})
	}

	c.JSON(http.StatusOK, user)
}
