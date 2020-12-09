package model

import "fmt"

const (
	dbName = "jwtDemo"
	userBucket = "user"
)


// User 用户类
type User struct {
	Id         string `json:"userId"`
	Name       string `json:"userName"`
	Gender     string `json:"gender"`
	Phone      string `json:"userMobile"`
	Pwd        string `json:"pwd"`
	Permission string `json:"permission"`
}

// LoginReq 登录请求参数类
type LoginReq struct {
	Phone string `json:"mobile"`
	Pwd   string `json:"pwd"`
}
// 序列化

func dumpUser(user User) []byte  {
	return nil
}

// 反序列化
func loadUser(jsonByte []byte) User {
	return User{}
}

//用户注册
func Register(phone string,password string) error {
	if CheckUser(phone) {
		return fmt.Errorf("用户已存在")
	}
	return nil
}

//检查用户是否存在
func CheckUser(phone string) bool  {
	if phone != "" {
		return true
	}
	return false
}

//登录限制
func LoginCheck(loinreq LoginReq) (bool, User, error) {

	resultUser := User{
		Id: "1",
		Name: "哈哈哈",
		Gender: "男",
		Phone: "18010000000",
		Pwd: "123456",
		Permission: "1",
	}

	return true,resultUser,nil
}