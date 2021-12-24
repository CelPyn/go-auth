package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"pynenborg.com/go-auth/pkg/domain"
	"pynenborg.com/go-auth/pkg/service"
	"pynenborg.com/go-auth/pkg/web"
)

func main() {
	usersPath := flag.String("u", "./resources/users.json", "The users file")
	userService := service.NewUserService(*usersPath)

	r := gin.Default()

	controllers := []domain.Controller{
		web.NewPingController(userService),
		domain.Controller(&web.PongController{}),
	}

	for i := range controllers {
		controller := controllers[i]
		controller.Setup(r)
	}

	// Listen and Serve in 0.0.0.0:8080
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
