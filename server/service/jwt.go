package service

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// secret key

var secretKey = []byte("abcd1234!@#$")

// 自定义声明
type ServiceCustiomClaims struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`

	jwt.StandardClaims
}

// jwtKeyFunc 返回密钥
func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return secretKey, nil
}

// Sign 生产token
func Sign(name, uid string) (string, error) {
	// 演示， 设置两分钟过期
	expAt := time.Now().Add(time.Duration(2) * time.Minute).Unix()

	// 创建声明
	claims := ServiceCustiomClaims{
		UserId: uid,
		Name:   name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expAt,
			Issuer: "system",
		},
	}
	// 创建token， 指定加密算法
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// 生成token

	return token.SignedString(secretKey)
}
