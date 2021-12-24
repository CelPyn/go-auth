package domain

type UserService interface {

	Login(name string, password string) (string, *HttpError)

	Get(name string) (User, *HttpError)

}
