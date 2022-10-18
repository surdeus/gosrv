package jwt

import (
	"github.com/golang-jwt/jwt"
)

type SigningMethod struct {
	jwt.SigningMethod
}

type Claims = jwt.MapClaims
type Token = jwt.Token

var (
	SigningMethodEdDSA = jwt.SignindMethodEdDSA
)

func New(method jwt.SigningMethod) (jwt.MapClaims, jwt.Token) {
	token := jwt.New(method)
	claims := token.Claims.(jwt.MapClaims)

	return Claims{claims}, Token{token}
}


