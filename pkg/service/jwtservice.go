package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"pynenborg.com/go-auth/pkg/domain"
	"regexp"
)

const (
	privateSecret = "super-secret-private-key"
)

func GenerateJwt(user domain.User) (string, *domain.HttpError) {
	claims := createClaims(user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, signError := token.SignedString([]byte(privateSecret))

	if signError != nil {
		return "", &domain.HttpError{
			Message: fmt.Sprintf("Something went wrong while generating JWT for user [%s]", user.Name),
			Status:  http.StatusInternalServerError,
		}
	}

	return signedString, nil
}

func ValidateBearer(token string) (*domain.HttpError) {
	bearerRegex, _ := regexp.Compile("Bearer (?P<jwt>.*)")
	matches := bearerRegex.MatchString(token)

	if !matches {
		return &domain.HttpError{
			Message: "Invalid Bearer Token: invalid format",
			Status:  http.StatusUnauthorized,
		}
	}

	match := bearerRegex.FindStringSubmatch(token)
	result := make(map[string]string)
	for i, name := range bearerRegex.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}


	jsonToken := result["jwt"]
	fmt.Println(jsonToken)

	jwToken, err := jwt.Parse(jsonToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(privateSecret), nil
	})

	if err != nil {
		return &domain.HttpError{
			Message: "Invalid or expired Bearer Token",
			Status:  http.StatusUnauthorized,
		}
	}

	if claims, ok := jwToken.Claims.(jwt.MapClaims); ok && jwToken.Valid {
		fmt.Println(claims["exp"])
	} else {
		return &domain.HttpError{
			Message: "Invalid or expired Bearer Token",
			Status:  http.StatusUnauthorized,
		}
	}

	return nil
}