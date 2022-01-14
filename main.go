package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"pynenborg.com/go-auth/pkg/domain"
	"pynenborg.com/go-auth/pkg/service"
	"pynenborg.com/go-auth/pkg/web"
	"time"
)

func main() {
	mongoClient := setupMongo()
	userService := service.NewUserService(mongoClient)

	controllers := []domain.Controller{
		web.NewLoginController(userService),
		web.NewUserController(userService),
	}

	r := setupRest(controllers)

	// Listen and Serve in 0.0.0.0:8080
	err := r.Run(":8080")
	if err != nil {
		log.Fatalln(err)
	}
}

func setupRest(controllers []domain.Controller) *gin.Engine {
	r := gin.Default()

	for i := range controllers {
		controller := controllers[i]
		controller.Setup(r)
	}

	return r
}

func setupMongo() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:example@localhost:27017/"))

	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
