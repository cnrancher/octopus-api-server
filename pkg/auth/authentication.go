package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/rancher/steve/pkg/auth"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/client-go/kubernetes"
)

const (
	group       = "system:masters"
	subjectName = "Subject"
)

type Authenticator interface {
	Authenticate(req *http.Request) (authed bool, user string, err error)
}

func NewK3sAuthenticator(endpoint string, client *kubernetes.Clientset, ctx context.Context) Authenticator {
	return &K3sAuthenticator{
		server:    endpoint,
		clientset: client,
		context:   ctx,
	}
}

type K3sAuthenticator struct {
	server    string
	clientset *kubernetes.Clientset
	context   context.Context
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

func (a *K3sAuthenticator) Authenticate(req *http.Request) (bool, string, error) {
	tokenAuthValue := GetTokenAuthFromRequest(req)
	if tokenAuthValue == "" {
		return false, "", fmt.Errorf("must authenticate")
	}

	secret, err := a.getJWTTokenSecret(tokenAuthValue)
	if err != nil {
		return false, "", fmt.Errorf("invalid token, error: %s", err.Error())
	}

	err = a.validJWTToken(tokenAuthValue, secret)
	if err != nil {
		return false, "", err
	}
	return true, string(secret.Data[subjectName]), nil
}

func (a *K3sAuthenticator) getJWTTokenSecret(token string) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	name, err := GetJWTSecretTokenName(token)
	if err != nil {
		return secret, err
	}
	secret, err = a.clientset.CoreV1().Secrets(tokenNamespace).Get(a.context, name, metav1.GetOptions{})
	if err != nil {
		return secret, err
	}
	return secret, nil
}

func (a *K3sAuthenticator) validJWTToken(tokenString string, secret *corev1.Secret) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret.Data[tokenSecretKeyName], nil
	})

	if token.Valid {
		return nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		// delete token secret if is not valid
		err = a.clientset.CoreV1().Secrets(tokenNamespace).Delete(a.context, secret.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
		// Token is either expired or not active yet
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return fmt.Errorf("that's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return fmt.Errorf("token timing is not valid")
		} else {
			return fmt.Errorf("couldn't handle this token: %s", err)
		}
	} else {
		err = a.clientset.CoreV1().Secrets(tokenNamespace).Delete(a.context, secret.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
		return fmt.Errorf("couldn't handle this token: %s", err)
	}
	return nil
}
