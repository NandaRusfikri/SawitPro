package repository

import (
	"SawitProRecruitment/entities"
	"SawitProRecruitment/pkg"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"time"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}

func (r *Repository) UserRegister(input UserRegistrationInput) (output UserRegistrationOutput, err error) {

	query := "SELECT COUNT(*) FROM users WHERE phone_number = $1"
	var count int
	err = r.Db.QueryRow(query, input.PhoneNumber).Scan(&count)
	if err != nil {
		log.Errorln("Error executing query:", err)
		return output, err
	}
	if count > 0 {
		log.Errorln("phoneNumber already exists in database.")
		return output, errors.New("Phone number already exists in database.")
	}

	insertUserSQL := "INSERT INTO users (phone_number, full_name, password) VALUES ($1, $2, $3) RETURNING id"
	err = r.Db.QueryRow(insertUserSQL, input.PhoneNumber, input.FullName, input.Password).Scan(&output.ID)
	if err != nil {
		log.Errorln("Failed to insert new user:", err)
		return output, err
	}
	return
}

func (r *Repository) FindUser(input FindUserInput) (output FindUserOutput, err error) {

	query := "SELECT id, full_name, phone_number FROM users WHERE id is not null "
	if input.ID > 0 {
		query += fmt.Sprintf("AND id = %d ", input.ID)
	}
	if input.PhoneNumber != "" {
		query += fmt.Sprintf("AND phone_number = '%s' ", input.PhoneNumber)
	}

	err = r.Db.QueryRow(query).Scan(&output.ID, &output.FullName, &output.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return output, errors.New("user not found")
		} else {
			log.Errorln("Failed to scan row:", err)
			return output, errors.New("user not found")
		}
	}

	return
}

func (r *Repository) UserLogin(input UserLoginInput) (output UserLoginOutput, err error) {
	query := "SELECT id, full_name, phone_number, password FROM users WHERE phone_number = $1"
	var EntityUser entities.EntityUser
	err = r.Db.QueryRow(query, input.PhoneNumber).Scan(&EntityUser.ID, &EntityUser.FullName, &EntityUser.PhoneNumber, &EntityUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Errorln("No rows were returned!")
			return output, errors.New("user not found")
		} else {
			log.Errorln("Error executing query:", err)
			return output, errors.New("user not found")
		}
	}

	CheckPassword := pkg.ComparePassword(EntityUser.Password, input.Password)
	if CheckPassword != nil {
		return output, errors.New("user not found")
	}

	expiredAt := time.Now().Add(time.Hour * time.Duration(1440))
	claims := pkg.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt.Unix(),
		},
		Id:          EntityUser.ID,
		FullName:    EntityUser.FullName,
		PhoneNumber: EntityUser.PhoneNumber,
	}

	access_token, ErrorJWT := pkg.Sign(claims)
	if ErrorJWT != nil {
		return output, errors.New("error creating token")
	}

	output.ID = EntityUser.ID
	output.AccessToken = access_token
	output.ExpiredAt = expiredAt.Unix()
	return
}

func (r *Repository) UpdateProfile(input UpdateProfileInput) (FindUserOutput, SchemaError) {

	var User FindUserOutput
	CheckUser := "SELECT id, full_name, phone_number FROM users WHERE id = $1"

	err := r.Db.QueryRow(CheckUser, input.ID).Scan(&User.ID, &User.FullName, &User.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Errorln("No rows were returned!")
		} else {
			log.Errorln("Failed to scan row:", err)
			return FindUserOutput{}, SchemaError{
				Code:    400,
				Message: "User not found",
			}
		}
	}

	if input.FullName != "" {
		User.FullName = input.FullName
	}
	if input.PhoneNumber != "" {
		CheckPhoneSQL := "SELECT count(*) FROM users WHERE phone_number = $1 AND id != $2"
		var count int
		err = r.Db.QueryRow(CheckPhoneSQL, input.PhoneNumber, input.ID).Scan(&count)
		if err != nil {
			log.Errorln("Error executing query:", err)
			return FindUserOutput{}, SchemaError{
				Code:    400,
				Message: "Error executing query",
			}
		}
		if count > 0 {
			return FindUserOutput{}, SchemaError{
				Code:    409,
				Message: "Phone number already exists in database.",
			}
		}
		User.PhoneNumber = input.PhoneNumber
	}

	queryUpdate := "UPDATE users SET full_name = $1, phone_number = $2 WHERE id = $3"

	fmt.Printf("User: %+v\n", User)
	_, err = r.Db.Query(queryUpdate, User.FullName, User.PhoneNumber, User.ID)
	if err != nil {
		log.Errorln("Failed to Update user:", err)
		return FindUserOutput{}, SchemaError{
			Code:    500,
			Message: "Failed to Update user",
		}
	}

	//User.ID = input.ID
	//User.FullName = input.FullName
	//User.PhoneNumber = input.PhoneNumber

	return User, SchemaError{}
}

func (r *Repository) InsertHistoryLogin(user_id int) (output int, err error) {

	insertUserSQL := "INSERT INTO history_login (user_id) VALUES ($1) RETURNING id"
	err = r.Db.QueryRow(insertUserSQL, user_id).Scan(&output)
	if err != nil {
		log.Errorln("Failed to insert history login:", err)
		return output, err
	}
	return
}
