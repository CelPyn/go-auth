package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"pynenborg.com/go-auth/pkg/domain"
	"pynenborg.com/go-auth/pkg/util/log"
)

type DefaultUserService struct {
	usersPath string
	users map[string]domain.User
}

type users struct {
	Users []domain.User `json:"users"`
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

func (us DefaultUserService) Login(name *string, password*string) error {
	return nil
}

func (us DefaultUserService) Get(name *string) domain.User {
	return us.users[*name]
}
