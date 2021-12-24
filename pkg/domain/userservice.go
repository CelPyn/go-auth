package domain

type UserService interface {

	Login(name *string, password *string) error

	Get(name *string) User

}
