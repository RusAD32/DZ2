package middleware

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
)

var (
	SecretKey      []byte      = []byte("OY8X4BtTYuBc73aStnGJtdcXKljWJjcOSbarTOUXb5lniOUTKaXb1axwnrMjc2P2jIzhhb")
	emptyValidFunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	}
)

var JwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: emptyValidFunc,
	SigningMethod:       jwt.SigningMethodHS256,
})
