package auth

import (
	"math/rand"
	"encoding/base64"
)

type Token string
type Tokens map[Token] any

type Sessions struct {
	Tokens Tokens
	TokenSize int
}

func New() Sessions {
	ret := Sessions{}
	ret.Tokens = make(Tokens)
	ret.TokenSize = 16

	return ret
}

func (auths Sessions)generateToken() Token {
	var (
	    ok bool
	)

	token := make([]byte, auths.TokenSize)
	for {
		rand.Read(token)
		_, ok = auths.Tokens[Token(token)]
		if !ok {
			break
		}
	}

	return Token(token)
}

func (auths Sessions)New(v any) Token {
	token := auths.generateToken()
	auths.Tokens[token] = v

	return token
}

func (auths Sessions)Get(token Token) (any, bool) {
	v, ok := auths.Tokens[token]
	return v, ok
}

func (auths Sessions)EncodeForClient(token Token) string {
	return base64.StdEncoding.EncodeToString([]byte(token))
}

func (auths Sessions)DecodeForServer(token string) (Token, error) {
	bts, err := base64.StdEncoding.DecodeString(token)
	return Token(bts), err
}

