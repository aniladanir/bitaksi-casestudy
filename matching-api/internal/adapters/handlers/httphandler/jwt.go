package httphandler

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

type JwtClaims struct {
	jwt.StandardClaims
	Authenticated bool `json:"authenticated"`
}

func (c JwtClaims) Valid() error {
	if !c.Authenticated {
		return errors.New("token is not authenticated")
	}
	return c.StandardClaims.Valid()
}
