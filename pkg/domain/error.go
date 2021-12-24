package domain

type HttpError struct {
	Cause error
	Message string
	Status int
}
