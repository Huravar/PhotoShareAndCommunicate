package crypt

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	CryPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码加密发生错误：%v", err.Error())
	} else {
		return string(CryPassword), nil
	}

}

func CheckPasswordHash(password string, CryPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(CryPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}
