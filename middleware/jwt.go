package middleware

import (
	"net/http"
	"time"

	"github.com/charfole/simple-tiktok/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// initialize a key to create token
var Key = []byte("charfole simple-tiktok secret key")

type MyClaims struct {
	UserId   uint   `json:"user_id"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

// create a token
func CreateToken(userId uint, userName string) (string, error) {
	expireTime := time.Now().Add(24 * time.Hour) //过期时间
	nowTime := time.Now()                        //当前时间
	claims := MyClaims{
		UserId:   userId,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间戳
			IssuedAt:  nowTime.Unix(),    //当前时间戳
			Issuer:    "charfole",        //颁发者签名
			Subject:   "tiktokToken",     //签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(Key)
}

// check a token
func CheckToken(token string) (*MyClaims, bool) {
	tokenObj, _ := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return Key, nil
	})
	if key, _ := tokenObj.Claims.(*MyClaims); tokenObj.Valid {
		return key, true
	} else {
		return nil, false
	}
}

// encapsulated into a jwt middleware
func JWTMiddleware() gin.HandlerFunc {
	// return a function
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}

		// user doesn't exist
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
		c.Set("user_id", tokenClaims.UserId)

		c.Next()
	}
}
