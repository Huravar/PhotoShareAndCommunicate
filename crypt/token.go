package crypt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
	"photo_service/model"
	"strconv"
	"time"
)

type UserTokenInfo struct {
	jwt.RegisteredClaims
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Identity uint   `json:"identity"`
}
type UserTokenBasicInfo struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Identity uint   `json:"identity"`
}

func GenerateToken(UserbasicInfo model.BasicUserInformation) (string, string, error) {
	UserClaim := UserTokenInfo{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.Itoa(int(UserbasicInfo.ID)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-5 * time.Second)),
		},
		UserId:   strconv.Itoa(int(UserbasicInfo.ID)),
		UserName: UserbasicInfo.UserName,
		Identity: UserbasicInfo.Identity,
	}
	TokenKey, _ := GenerateSecureSecret(32)
	UserToken := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	SignedToken, err := UserToken.SignedString([]byte(TokenKey))
	if err != nil {
		return "", "", fmt.Errorf("token 生成失败：%v", err.Error())
	}
	return SignedToken, TokenKey, nil
}

func GenerateSecureSecret(length int) (string, error) {
	secret := make([]byte, length)
	if _, err := rand.Read(secret); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(secret), nil
}

func ParasedAndVerify(PToken string, PTokenKey string) (UserTokenBasicInfo, error) {
	ItemParsedToken, err := jwt.ParseWithClaims(
		PToken,           // 待解析的 Token 字符串
		&UserTokenInfo{}, // 自定义声明结构体的指针（用于存储解析后的数据）
		func(t *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("无效签名方法: %v", t.Header["alg"])
			}
			return []byte(PTokenKey), nil // 返回用于验证签名的密钥
		},
	)
	if ItemParasClaim, ok := ItemParsedToken.Claims.(*UserTokenInfo); ok && ItemParsedToken.Valid {
		return UserTokenBasicInfo{
			UserId:   ItemParasClaim.UserId,
			UserName: ItemParasClaim.UserName,
			Identity: ItemParasClaim.Identity,
		}, nil
	} else {
		return UserTokenBasicInfo{}, err
	}

}
