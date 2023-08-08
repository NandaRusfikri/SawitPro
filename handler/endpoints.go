package handler

import (
	"SawitProRecruitment/pkg"
	"SawitProRecruitment/repository"
	"SawitProRecruitment/util"
	"fmt"
	"net/http"

	"SawitProRecruitment/generated"
	"github.com/labstack/echo/v4"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)

func (s *Server) Default(ctx echo.Context) error {

	return ctx.JSON(http.StatusOK, generated.DefaultResponse{
		Version: "1.0.0",
		Author:  "NandaRusfikri",
	})
}

func (s *Server) Hello(ctx echo.Context, params generated.HelloParams) error {

	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) Login(ctx echo.Context) error {

	var u generated.LoginRequest
	if err := ctx.Bind(&u); err != nil {
		return err
	}

	output, err := s.Repository.UserLogin(repository.UserLoginInput{
		PhoneNumber: u.PhoneNumber,
		Password:    u.Password,
	})
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}
	go s.Repository.InsertHistoryLogin(output.ID)

	resp := generated.LoginResponse{
		Id:          output.ID,
		Message:     "Login Success",
		AccessToken: output.AccessToken,
		ExpiresIn:   output.ExpiredAt,
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) UserRegister(ctx echo.Context) error {

	var u generated.UserRegisterRequest
	if err := ctx.Bind(&u); err != nil {
		return err
	}
	if !util.ValidateFullName(u.FullName) {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Full name is not valid",
		})
	}
	if !util.ValidatePhoneNumber(u.PhoneNumber) {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Phone number is not valid",
		})
	}
	if !util.ValidatePassword(u.Password) {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Password must be at least 6 characters, contain at least 1 uppercase letter, 1 lowercase letter, and 1 number",
		})
	}

	output, err := s.Repository.UserRegister(repository.UserRegistrationInput{
		PhoneNumber: u.PhoneNumber,
		Password:    pkg.HashPassword(u.Password),
		FullName:    u.FullName,
	})
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}
	resp := generated.UserRegisterResponse{
		Id:      output.ID,
		Message: "Success",
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) GetMyProfile(ctx echo.Context) error {

	Claims, err := pkg.Auth(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Unauthorized",
		})
	}

	output, err := s.Repository.FindUser(repository.FindUserInput{
		ID: Claims.Id,
	})
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, generated.ProfileResponse{
		Id:          output.ID,
		PhoneNumber: output.PhoneNumber,
		FullName:    output.FullName,
	})
}

func (s *Server) UpdateProfile(ctx echo.Context) error {

	Claims, err := pkg.Auth(ctx)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, map[string]interface{}{
			"message": "Unauthorized",
		})
	}
	var u generated.UpdateProfileRequest
	if err := ctx.Bind(&u); err != nil {
		return err
	}

	output, errDB := s.Repository.UpdateProfile(repository.UpdateProfileInput{
		ID:          Claims.Id,
		PhoneNumber: u.PhoneNumber,
		FullName:    u.FullName,
	})
	if errDB != (repository.SchemaError{}) {
		return ctx.JSON(errDB.Code, map[string]interface{}{
			"message": errDB.Message,
		})
	}

	return ctx.JSON(http.StatusOK, generated.ProfileResponse{
		Id:          output.ID,
		PhoneNumber: output.PhoneNumber,
		FullName:    output.FullName,
	})
}
