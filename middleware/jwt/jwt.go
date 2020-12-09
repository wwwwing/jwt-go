package jwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func JwtAuth() gin.HandlerFunc  {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == ""{
			 c.JSON(http.StatusOK,gin.H{
			 	"status":-1,
			 	"msg":"token不存在",
			 })
			 c.Abort()
			return
		}
		log.Println("get token is : ",token)

		jwt := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := jwt.ParseToken(token)
		if err != nil {
			if err == TokenInvalid {
				c.JSON(http.StatusOK,gin.H{
					"status":-1,
					"msg":"授权已过期",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusOK,gin.H{
				"status":-1,
				"msg":err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("claims",claims)
	}
}


// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

//解析token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString,&CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey,nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims,nil
	}
	return nil,TokenInvalid
}

//刷新token
func (j *JWT) RefreshToken(tokenString string) (string,error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString,&CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey,nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "",TokenInvalid
}

//生成token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	ID    string `json:"userId"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	jwt.StandardClaims
}

// 新建一个jwt实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

// 获取signKey
func GetSignKey() string {
	return SignKey
}

// 这是SignKey
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}


// 一些常量
var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "newtrekWang"
)