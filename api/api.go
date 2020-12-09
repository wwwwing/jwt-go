package api

import (
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"jwtdemo/middleware/jwt"
	"jwtdemo/model"
	"log"
	"net/http"
	"time"
)

type RegisterInfo struct {
	Phone string `json:"phone"`
	Password string `json:"password"`
}

func RegisterUser(c *gin.Context)  {
	var registerInfo RegisterInfo

	if c.BindJSON(&registerInfo) == nil {
		//model
		err := model.Register(registerInfo.Phone, registerInfo.Password)
		if err == nil {
			c.JSON(http.StatusOK,gin.H{
				"status":0,
				"msg":"注册成功",
			})
		}else {
			c.JSON(http.StatusOK,gin.H{
				"status":0,
				"msg":"注册失败" + err.Error(),
			})
		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":-1,
			"msg":"解析数据失败！",
		})
	}
	
}

//登录结果结构
type LoginResult struct {
	Token string `json:"token"`
	model.User
}

func Login(c *gin.Context)  {
	var loginReq model.LoginReq
	if c.ShouldBind(&loginReq) == nil{
		check, user, err := model.LoginCheck(loginReq)
		if check {
			generateToken(c,user)
		}else {
			c.JSON(http.StatusOK,gin.H{
				"status":-1,
				"msg":"登录失败" + err.Error(),
			})
		}
	}else {
		c.JSON(http.StatusOK,gin.H{
			"status":-1,
			"msg":"参数错误",
		})
	}

}

func generateToken(c *gin.Context, user model.User) {
	j := &jwt.JWT{
		[]byte("newtrekWang"),
	}
	claims := jwt.CustomClaims{
		user.Id,
		user.Name,
		user.Phone,
		jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),//签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600),//签名过期时间
			Issuer:    "newtrekWang",                   //签名的发行者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"status":-1,
			"msg":err.Error(),
		})
		return
	}
	log.Println("create token is ",token)

	data := LoginResult{
		User: user,
		Token: token,
	}
	c.JSON(http.StatusOK,gin.H{
		"status":0,
		"msg":"登录成功",
		"data":data,
	})
	return
}

func GetDetail(c *gin.Context) {
	claims := c.MustGet("claims").(*jwt.CustomClaims)
	if claims != nil {
		c.JSON(http.StatusOK,gin.H{
			"status":0,
			"msg":"ok",
			"data":claims,
		})
	}
	return
}





















