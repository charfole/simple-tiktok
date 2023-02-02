package service

import (
	"github.com/charfole/simple-tiktok/common"
	"github.com/charfole/simple-tiktok/dao/mysql"
	"github.com/charfole/simple-tiktok/middleware"
	"github.com/charfole/simple-tiktok/model"
	"golang.org/x/crypto/bcrypt"
)

// Service of UserLogin
func UserLoginService(username string, passWord string) (TokenResponse, error) {
	// 1. prepare the response
	var userResponse = TokenResponse{}

	// 2. check the user is valid or not
	err := IsUserValid(username, passWord)
	if err != nil {
		return userResponse, err
	}

	// 3. if user not found, login fail
	var login model.User
	err = mysql.GetAUser(username, &login)
	if err != nil {
		return userResponse, err
	}

	// 4. if user password not matched, login fail
	if !CheckPassword(login.Password, passWord) {
		return userResponse, common.ErrorPasswordFalse
	}

	// 5. creates a new token for the valid user
	token, err := middleware.CreateToken(login.Model.ID, login.Name)
	if err != nil {
		return userResponse, err
	}

	// 6. return the UserID and token to controller layer
	userResponse = TokenResponse{
		UserID: login.Model.ID,
		Token:  token,
	}
	return userResponse, nil
}

// func IsUserExist(username string, password string, login *model.User) error {
// 	// call the mysql to query a user
// 	err := mysql.GetAUser(username, login)
// 	// user not found or other unpredicted errors
// 	if err != nil {
// 		return err
// 	}

// 	// user found but passwrod not matched
// 	if !CheckPassword(login.Password, password) {
// 		return common.ErrorPasswordFalse
// 	}

// 	// success
// 	return nil
// }

// check the password
func CheckPassword(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPwd)
	// compare the hashed password and original password
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	return err == nil
}
