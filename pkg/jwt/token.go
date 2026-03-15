package jwt

import (
	"InnerG/config"
	"InnerG/pkg/constants"
	"InnerG/pkg/errno"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	UserId string `json:"user_id"`
	Type   int    `json:"type"`
	jwt.RegisteredClaims
}

func CreateAllToken(userid string) (string, string, error) {
	accessToken, err := CreateToken(constants.TypeAccessToken, userid)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := CreateToken(constants.TypeRefreshToken, userid)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func CreateToken(tokenType int, UserId string) (string, error) {
	if config.Service.PrivateKey == "" {
		return "", errno.AuthError.WithMessage("config empty")
	}

	var expireTime time.Time
	nowTime := time.Now()

	switch tokenType {
	case constants.TypeAccessToken:
		expireTime = nowTime.Add(constants.AccessTokenTTL)
	case constants.TypeRefreshToken:
		expireTime = nowTime.Add(constants.RefreshTokenTTL)
	}
	claims := Claims{
		UserId: UserId,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			Issuer:    constants.Issuer,
		},
	}
	// 加密--数字签名部分
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	key, err := jwt.ParseEdPrivateKeyFromPEM([]byte(config.Service.PrivateKey))
	if err != nil {
		return "", errno.AuthError.WithMessage(fmt.Sprintf("parse private key failed, err: %v", err))
	}

	token, err := tokenStruct.SignedString(key)
	if err != nil {
		return "", errno.AuthError.WithMessage(fmt.Sprintf("sign token failed, err: %v", err))
	}
	return token, nil
}

func CheckToken(token string) (tokeType int, user_id string, err error) {
	if token == "" {
		return -1, "", errno.AuthMissing
	}
	tokenStruct, _, err := new(jwt.Parser).ParseUnverified(token, &Claims{})
	if err != nil {
		return -1, "", errno.AuthInvalid
	}

	unverifiedClaims, ok := tokenStruct.Claims.(*Claims)
	if !ok {
		return -1, "", errno.AuthError.WithMessage("cannot handle claims")
	}

	secret, err := jwt.ParseEdPublicKeyFromPEM([]byte(constants.PublicKey))
	if err != nil {
		return -1, "", errno.AuthError.WithMessage(fmt.Sprintf("parse public key failed, err: %v", err))
	}

	// 使用正确的密钥再次解析 token
	response, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, errno.AuthError.WithMessage(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}
		return secret, nil
	})
	// 验证 token 是否有效
	if err != nil {
		return unverifiedClaims.Type, "", checkError(err, unverifiedClaims.Type)
	}

	if _, ok := response.Claims.(*Claims); ok && response.Valid {
		return unverifiedClaims.Type, unverifiedClaims.UserId, nil
	}

	return -1, "", errno.AuthInvalid

}

func checkError(err error, tokenType int) error {
	var ve *jwt.ValidationError
	if errors.As(err, &ve) {
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			if tokenType == constants.TypeAccessToken {
				return errno.AuthAccessExpired
			}
			return errno.AuthRefreshExpired
		}
	}
	return errno.AuthError.WithMessage(err.Error())
}
