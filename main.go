package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jwtdemo/api"
	"jwtdemo/middleware/jwt"
)

func main() {
	r := gin.Default()
	//r.POST("login")
	fmt.Println("this is a homepage")
	r.POST("/login",api.Login)
	r.POST("/register",api.RegisterUser)

	group := r.Group("/data")
	group.Use(jwt.JwtAuth())
	{
		group.GET("/get",api.GetDetail)
	}
	r.Run(":8888")
}
