package service

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"pynenborg.com/go-auth/pkg/domain"
	"time"
)

type DefaultUserService struct {
	mongoClient *mongo.Client
}

type customClaims struct {
	Username string `json:"username"`
	Grants []string `json:"grants"`
	jwt.StandardClaims
}

func NewUserService(mongoClient *mongo.Client) DefaultUserService {
	return DefaultUserService{mongoClient: mongoClient}
}

func (us DefaultUserService) Create(name string, password string) {

}

func (us DefaultUserService) Get(name string) (domain.User, *domain.HttpError){
	user, httpError := us.findUser(name)
	return *user, httpError
}

func (us DefaultUserService) Login(name string, password string) (string, *domain.HttpError) {

	user, httpError := us.findUser(name)

	if httpError != nil {
		return "", httpError
	}

	passwordMatches := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil

	if !passwordMatches {
		return "", &domain.HttpError{
			Message: fmt.Sprintf("Login attempt for user [%s] failed: password did not match", name),
			Status: http.StatusUnauthorized,
		}
	}

	return GenerateJwt(*user)
}

func (us DefaultUserService) findUser(name string) (*domain.User, *domain.HttpError) {
	var user domain.User
	err := us.mongoClient.Database("admin").Collection("user").FindOne(context.TODO(), bson.D{{"name", name}}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &domain.HttpError{
				Cause:   err,
				Message: fmt.Sprintf("User [%s] could not be found", name),
				Status:  404,
			}
		}
		return nil, &domain.HttpError{
			Cause:   err,
			Message: fmt.Sprintf("Error when looking for user [%s]", name),
			Status:  500,
		}
	}

	return &user, nil
}

func createClaims(user domain.User) customClaims {
	expiresAt := time.Now().Unix() + 1

	return customClaims{
		Username: user.Name,
		Grants: user.Grants,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expiresAt, Issuer: "go-auth.pynenborg.com"},
	}
}
