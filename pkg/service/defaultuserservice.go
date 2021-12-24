package service

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"os"
	"pynenborg.com/go-auth/pkg/domain"
	"pynenborg.com/go-auth/pkg/util/log"
	"time"
)

const (
	PRIVATE_SECRET = "super-secret-private-key"
)

type DefaultUserService struct {
	usersPath string
	users map[string]domain.User
}

type users struct {
	Users []domain.User `json:"users"`
}

type customClaims struct {
	Username string `json:"username"`
	Grants []string `json:"grants"`
	jwt.StandardClaims
}

func NewUserService(usersPath string) DefaultUserService {
	users := readUserFile(usersPath)
	log.Info("Read all users: %s", users)
	return DefaultUserService{usersPath: usersPath, users: users}
}

func readUserFile(usersPath string) map[string]domain.User {
	userFile, err := os.Open(usersPath)
	if err != nil {
		fmt.Printf("Error reading usersFile at [%s]: %s\n", usersPath, err)
		// create an empty map and return it
		return make(map[string]domain.User)
	}

	defer func(userFile *os.File) {
		err := userFile.Close()
		if err != nil {
			fmt.Printf("Error closing usersFile at [%s]: %s\n", usersPath, err)
		}
	}(userFile)

	return unmarshal(userFile)
}

func unmarshal(userFile *os.File) map[string]domain.User {
	content, _ := ioutil.ReadAll(userFile)

	var userSlice users
	err := json.Unmarshal(content, &userSlice)
	if err != nil {
		fmt.Printf("Error closing usersFile at [%s]: %s\n", userFile.Name(), err)
		return map[string]domain.User{}
	}

	users := make(map[string]domain.User)
	for _, user := range userSlice.Users {
		users[user.Name] = user
	}

	return users
}

func (us DefaultUserService) Login(name string, password string) (string, *domain.HttpError) {
	user, err := us.Get(name)

	if err != nil {
		return "", err
	}

	passwordMatches := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil

	if !passwordMatches {
		return "", &domain.HttpError{
			Message: fmt.Sprintf("Login attempt for user %s failed: password did not match", name),
			Status: http.StatusUnauthorized,
		}
	}

	return generateJwt(user)
}

func generateJwt(user domain.User) (string, *domain.HttpError) {
	claims := createClaims(user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, signError := token.SignedString([]byte(PRIVATE_SECRET))

	if signError != nil {
		return "", &domain.HttpError{
			Message: fmt.Sprintf("Something went wrong while generating JWT for user [%s]", user.Name),
			Status:  http.StatusInternalServerError,
		}
	}

	return signedString, nil
}

func createClaims(user domain.User) customClaims {
	expiresAt := time.Now().Unix() + 3600

	return customClaims{
		Username: user.Name,
		Grants: user.Grants,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expiresAt, Issuer: "go-auth.pynenborg.com"},
	}
}

func (us DefaultUserService) Get(name string) (domain.User, *domain.HttpError) {
	user := us.users[name]

	if user.Name == "" {
		return user, &domain.HttpError{
			Message: fmt.Sprintf("User [%s] could not be found", name),
			Status:  http.StatusNotFound,
		}
	}

	return user, nil
}
