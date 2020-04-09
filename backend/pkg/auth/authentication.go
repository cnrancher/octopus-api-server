package auth

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/rancher/steve/pkg/auth"
	"k8s.io/apiserver/pkg/authentication/user"
)

const group = "system:masters"
const name = "admin"

type Authenticator interface {
	Authenticate(req *http.Request) (authed bool, user string, err error)
}

func ToAuthMiddleware(a Authenticator) auth.Middleware {
	f := func(req *http.Request) (user.Info, bool, error) {
		authed, u, err := a.Authenticate(req)
		return &user.DefaultInfo{
			Name:   u,
			UID:    u,
			Groups: []string{group},
		}, authed, err
	}
	return auth.ToMiddleware(auth.AuthenticatorFunc(f))
}

type K3sAuthenticator struct {
	server   string
	password string
}

func (a *K3sAuthenticator) Authenticate(req *http.Request) (bool, string, error) {
	tokenAuthValue := GetTokenAuthFromRequest(req)
	if tokenAuthValue == "" {
		return false, "", fmt.Errorf("must authenticate")
	}
	err := validatToken(tokenAuthValue, a.password)
	if err != nil {
		return false, "", err
	}
	return true, name, nil
}

func validatToken(tokenString string, pwd string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(pwd), nil
	})

	if token.Valid {
		fmt.Println("You look nice today")
		return nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return fmt.Errorf("that's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return fmt.Errorf("timing is everything")
		} else {
			return fmt.Errorf("couldn't handle this token: %s", err)
		}
	} else {
		return fmt.Errorf("couldn't handle this token: %s", err)
	}
	return nil
}

func NewK3sAuthenticator(endpoint string, password string) Authenticator {
	return &K3sAuthenticator{
		server:   endpoint,
		password: password,
	}
}
