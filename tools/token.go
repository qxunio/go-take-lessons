package tools

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go-take-lessons/configs"
	"time"
)

const KEY = "688e58432579e36d33c45f796e9b94db"

// 生成token
func GenToken(userId string) (string, error) {
	claims := &jwt.StandardClaims{
		Id:        userId,
		Audience:  "user",
		ExpiresAt: time.Now().Unix() + configs.Conf.App.LoginTime,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "qxun",
		NotBefore: time.Now().Unix(),
		Subject:   "user-profile",
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(KEY))

	if err != nil {
		return "", err
	}
	return token, nil
}

// 校验token, 返回usid
func ValidateToken(token string) (string, error) {
	t, pe := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(KEY), nil
	})
	if pe != nil {
		return "", pe
	}
	if claims, ok := t.Claims.(*jwt.StandardClaims); ok && t.Valid {
		return claims.Id, nil
	}
	return "", errors.New("Token Validate Fail ")
}
