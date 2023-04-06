package authx

import (
	"math/rand"
	"encoding/base64"
)

type Token string
type Tokens[S any] map[Token] S

type Sessions[S any] struct {
	Tokens Tokens[S]
	TokenSize int
}

func New[S any]() Sessions[S] {
	ret := Sessions[S]{}
	ret.Tokens = make(Tokens[S], 50)
	ret.TokenSize = 16

	return ret
}

func (auths Sessions[S])generateToken() Token {
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

func (auths Sessions[S])New(v S) Token {
	token := auths.generateToken()
	auths.Tokens[token] = v

	return token
}

func (auths Sessions[S])Get(token Token) (any, bool) {
	v, ok := auths.Tokens[token]
	return v, ok
}

func (auths Sessions[S])EncodeForClient(token Token) string {
	return base64.StdEncoding.EncodeToString([]byte(token))
}

func (auths Sessions[S])DecodeForServer(token string) (Token, error) {
	bts, err := base64.StdEncoding.DecodeString(token)
	return Token(bts), err
}

