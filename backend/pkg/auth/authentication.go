package auth

import (
	"fmt"
	"net/http"

	"github.com/rancher/steve/pkg/auth"
	"k8s.io/apiserver/pkg/authentication/user"
)

const group = "system:masters"

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
	server string
}

func (a *K3sAuthenticator) Authenticate(req *http.Request) (bool, string, error) {
	tokenAuthValue := GetTokenAuthFromRequest(req)
	if tokenAuthValue == "" {
		return false, "", fmt.Errorf("must authenticate")
	}
	tokenName, tokenKey := SplitTokenParts(tokenAuthValue)
	if tokenName != "admin" || tokenKey == "" {
		return false, "", fmt.Errorf("must authenticate")
	}
	_, err := NormalizeAndValidateTokenForUser(a.server, tokenKey, tokenName)
	if err != nil {
		return false, "", fmt.Errorf("must authenticate")
	}
	return true, tokenName, nil
}

func NewK3sAuthenticator(endpoint string) Authenticator {
	return &K3sAuthenticator{
		server: endpoint,
	}
}
