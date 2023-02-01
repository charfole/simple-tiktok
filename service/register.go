package service

import (
	"github.com/charfole/simple-tiktok/common"
	"github.com/charfole/simple-tiktok/dao/mysql"
	"github.com/charfole/simple-tiktok/middleware"
	"github.com/charfole/simple-tiktok/model"
	"golang.org/x/crypto/bcrypt"
)

const (
	MaxUsernameLength = 32 // the max length of username
	MaxPasswordLength = 32 // the max length of passwrod
	MinPasswordLength = 6  // the minimun length of password
)

type TokenResponse struct {
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}

func UserRegisterService(userName string, passWord string) (TokenResponse, error) {

	// 1. create the tokenResponse for return
	var tokenResponse = TokenResponse{}

	// 2. check the username and password is valid or not
	err := IsUserValid(userName, passWord)
	// if invalid return error
	if err != nil {
		return tokenResponse, err
	}

	// 3. register for a new user
	newUser, err := CreateRegisterUser(userName, passWord)
	if err != nil {
		return tokenResponse, err
	}

	// 4. create a token for this new user
	token, err := middleware.CreateToken(newUser.ID, newUser.Name)
	if err != nil {
		return tokenResponse, err
	}

	// 5. pack the data and return
	tokenResponse = TokenResponse{
		UserID: newUser.ID,
		Token:  token,
	}
	return tokenResponse, nil
}

// check the username and password is valid or not
func IsUserValid(username string, password string) error {
	// 1. check the length of username is in [1, 32] or not
	if len(username) == 0 {
		return common.ErrorUserNameEmpty
	}
	if len(username) > MaxUsernameLength {
		return common.ErrorUserNameInvalid
	}

	// 2. check the length of password is in [6,32] or not
	if len(password) == 0 {
		return common.ErrorPasswordEmpty
	}
	if len(password) > MaxPasswordLength || len(password) < MinPasswordLength {
		return common.ErrorPasswordInvalid
	}
	return nil
}

// register for a new user
func CreateRegisterUser(username string, password string) (model.User, error) {
	// 1. hash the original password and create user model
	hashPassword, _ := HashAndSalt(password)
	newUser := model.User{
		Name:     username,
		Password: hashPassword,
	}

	// 2. migrate the user model to MySQL "users" table
	err := mysql.DB.AutoMigrate(&model.User{})
	if err != nil {
		return newUser, common.ErrorDBMigrateFail
	}

	// 3.check the user if exists or not
	isExist, err := mysql.IsUserExist(username)

	if err != nil {
		return model.User{}, err
	}

	if isExist {
		// if exists
		return newUser, err
	} else {
		// if not, create a new user to MySQL "users" table
		if err := mysql.CreateAUser(&newUser); err != nil {
			return newUser, err
		}
	}
	return newUser, nil
}

// encrypt the code by bcrypt module
func HashAndSalt(password string) (passwordHash string, err error) {
	// 1. convert the type from string to []byte
	pwd := []byte(password)
	// 2. hash the original password
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return
	}
	// 3. convert back and return the hashed password
	passwordHash = string(hash)
	return
}
