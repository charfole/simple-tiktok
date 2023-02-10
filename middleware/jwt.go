package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/charfole/simple-tiktok/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// initialize a key to create token
var SecretKey = []byte("charfole simple-tiktok secret key")

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return SecretKey, nil
}

type MyClaims struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

// create a token
func CreateToken(userID uint, userName string) (string, error) {
	expireTime := time.Now().Add(24 * time.Hour) //过期时间
	nowTime := time.Now()                        //当前时间
	claims := MyClaims{
		UserID:   userID,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间戳
			IssuedAt:  nowTime.Unix(),    //当前时间戳
			Issuer:    "charfole",        //颁发者签名
			Subject:   "tiktokToken",     //签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// check a token
func CheckToken(tokenStr string) (*MyClaims, bool) {
	tokenObj, _ := jwt.ParseWithClaims(tokenStr, &MyClaims{}, keyFunc)
	if token, _ := tokenObj.Claims.(*MyClaims); tokenObj.Valid {
		return token, true
	} else {
		return nil, false
	}
}

// encapsulated into a jwt middleware
func JWTMiddleware() gin.HandlerFunc {
	// return a function
	return func(c *gin.Context) {
		fmt.Println("visit JWT middleware!")
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}

		// user not found
		if tokenStr == "" {
			c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "用户不存在"})
			c.Abort() // abort the upcoming call
			return
		}

		// check the token
		tokenClaims, ok := CheckToken(tokenStr)
		if !ok {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  common.ErrorTokenFalse.Error(),
			})
			c.Abort() // abort the upcoming call
			return
		}
		// token expires
		if time.Now().Unix() > tokenClaims.ExpiresAt {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  common.ErrorExpired.Error(),
			})
			c.Abort() // abort the upcoming call
			return
		}

		c.Set("username", tokenClaims.UserName)
		c.Set("user_id", tokenClaims.UserID)

		c.Next()
	}
}
