// This file contains types that are used in the repository layer.
package repository

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type UserRegistrationInput struct {
	PhoneNumber string
	FullName    string
	Password    string
}
type UserRegistrationOutput struct {
	ID int
}

type UserLoginInput struct {
	PhoneNumber string
	Password    string
}
type UserLoginOutput struct {
	ID          int
	Message     string
	AccessToken string
	ExpiredAt   int64
}

type FindUserInput struct {
	PhoneNumber string
	ID          int
}
type FindUserOutput struct {
	ID          int
	PhoneNumber string
	FullName    string
}

type UpdateProfileInput struct {
	ID          int
	PhoneNumber string
	FullName    string
}

type SchemaError struct {
	Code    int
	Message string
}
