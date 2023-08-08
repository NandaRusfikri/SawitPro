// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"
)

type RepositoryInterface interface {
	GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error)
	UserRegister(input UserRegistrationInput) (output UserRegistrationOutput, err error)
	UserLogin(input UserLoginInput) (output UserLoginOutput, err error)
	FindUser(input FindUserInput) (output FindUserOutput, err error)
	UpdateProfile(input UpdateProfileInput) (FindUserOutput, SchemaError)
	InsertHistoryLogin(user_id int) (output int, err error)
}
