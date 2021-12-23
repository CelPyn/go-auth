package main

import (
	"github.com/gin-gonic/gin"
	"pynenborg.com/go-auth/pkg/service"
	"pynenborg.com/go-auth/pkg/web"
)

func main() {
	r := gin.Default()

	controllers := []service.Controller{web.PingController{}, web.PongController{}}

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
