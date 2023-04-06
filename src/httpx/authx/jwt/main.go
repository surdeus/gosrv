package jwt

import (
	"github.com/golang-jwt/jwt"
)

type Claims = jwt.MapClaims
type Token = jwt.Token
type SigningMethod = jwt.SigningMethod
type KeyFunc = jwt.KeyFunc

var (
	SigningMethodEDDSA = jwt.SigningMethodEDDSA
)

func New(method jwt.SigningMethod) (Claims, *Token) {
	token := jwt.New(method)
	claims := token.Claims.(jwt.MapClaims)

	return claims, token
}

func SignedString(token *Token, secretKey string) (string, error) {
	signedString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func Parse(token *Token, tokenString string, keyFunc KeyFunc) {
	token.Parse()
}
