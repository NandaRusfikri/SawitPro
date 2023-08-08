package handler

import (
	"SawitProRecruitment/generated"
	"SawitProRecruitment/repository"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const token = "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTY1MTAzMjEsImlkIjoyMCwiZnVsbF9uYW1lIjoiTmFuZGFSdXNmaWtyaSIsInBob25lX251bWJlciI6Iis2MjgyMzEyNjY0NTY0In0.QKGwpK2_YWQeMkcIjhnP-ltKLOrfAmhzhTu4BGbQb2L7FuEyCIaf7vzd0VNZxca3WxktpXpfssuBZuNNRv8Gi8BHPknD5ReCOH5mswLNaRWl882je5po-Rtwy71VqJb2t3r0kbri9_mQ6E8yds2tGNkvFvVSEwNiMW7v2HFrd3vY7KfgqXkaMja_AAf50DRaqmAwkpIKBCC2afqYWdZUlwQWnGvTjSL9BVCoC_TCbm4SEs_qJq9NPjgGEyBhzwQZHQ3mBJZ8DlbFIfwKtiLp6Kp3AK78ZeWG-4yANp-0FevYh6YKQCSc4kQWZh9uTDDPIzuqd7labpWEpOkRa64p4w"

func TestDefault(t *testing.T) {
	e := echo.New()

	server := &Server{
		Repository: &MockRepository{},
	}

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := server.Default(ctx)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResponse := `{"version":"1.0.0","author":"NandaRusfikri"}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}
}

func TestHello(t *testing.T) {
	e := echo.New()

	server := &Server{
		Repository: &MockRepository{},
	}

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := server.Hello(ctx, generated.HelloParams{Id: 1})
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResponse := `{"message":"Hello User 1"}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}
}

func TestLogin(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new Server instance
	server := &Server{
		// Initialize the repository with a mock repository for testing
		Repository: &MockRepository{},
	}

	// Bind the login request JSON to the context
	reqBody := `{"phone_number": "+6282312664564", "password": "Password1!"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Call the Login method
	err := server.Login(ctx)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// Validate the response JSON
		expectedResponse := `{"id":1,"message":"Login Success","access_token":"testaccesstoken","expires_in":3600}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}
}

func TestUserRegister(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new Server instance
	server := &Server{
		// Initialize the repository with a mock repository for testing
		Repository: &MockRepository{},
	}

	// Bind the user register request JSON to the context
	reqBody := `{"full_name": "NandaRusfikri", "phone_number": "+6282312664564", "password": "Password1!"}`
	req := httptest.NewRequest(http.MethodPost, "/user/register", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Call the UserRegister method
	err := server.UserRegister(ctx)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// Validate the response JSON
		expectedResponse := `{"id":1,"message":"Success"}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}
}

func TestGetMyProfile(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new Server instance
	server := &Server{
		// Initialize the repository with a mock repository for testing
		Repository: &MockRepository{},
	}

	// Create a new request with the authorization header containing a JWT token
	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	req.Header.Set(echo.HeaderAuthorization, token)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Call the GetMyProfile method
	err := server.GetMyProfile(ctx)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// Validate the response JSON
		expectedResponse := `{"id":1,"phone_number":"1234567890","full_name":"John Doe"}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}
}

func TestUpdateProfile(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new Server instance
	server := &Server{
		// Initialize the repository with a mock repository for testing
		Repository: &MockRepository{},
	}

	// Bind the update profile request JSON to the context
	reqBody := `{"full_name": "John Doe Jr.", "phone_number": "9876543210"}`
	req := httptest.NewRequest(http.MethodPost, "/update-profile", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, token)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Call the UpdateProfile method
	err := server.UpdateProfile(ctx)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// Validate the response JSON
		expectedResponse := `{"id":1,"phone_number":"9876543210","full_name":"John Doe Jr."}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}
}

// MockRepository is a mock implementation of the repository.RepositoryInterface for testing
type MockRepository struct{}

func (m *MockRepository) GetTestById(ctx context.Context, input repository.GetTestByIdInput) (output repository.GetTestByIdOutput, err error) {
	return repository.GetTestByIdOutput{
		Name: "testname",
	}, nil
}

func (m *MockRepository) UserLogin(input repository.UserLoginInput) (repository.UserLoginOutput, error) {
	return repository.UserLoginOutput{
		ID:          1,
		AccessToken: "testaccesstoken",
		ExpiredAt:   3600,
	}, nil

}

func (m *MockRepository) UserRegister(input repository.UserRegistrationInput) (repository.UserRegistrationOutput, error) {
	return repository.UserRegistrationOutput{
		ID: 1,
	}, nil
}

func (m *MockRepository) FindUser(input repository.FindUserInput) (repository.FindUserOutput, error) {
	return repository.FindUserOutput{
		ID:          1,
		PhoneNumber: "1234567890",
		FullName:    "John Doe",
	}, nil
}

func (m *MockRepository) UpdateProfile(input repository.UpdateProfileInput) (repository.FindUserOutput, repository.SchemaError) {
	return repository.FindUserOutput{
		ID:          1,
		PhoneNumber: "9876543210",
		FullName:    "John Doe Jr.",
	}, repository.SchemaError{}
}
func (m *MockRepository) InsertHistoryLogin(user_id int) (output int, err error) {
	return 1, nil
}
