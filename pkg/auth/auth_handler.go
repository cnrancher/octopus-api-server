package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/cnrancher/octopus-api-server/pkg/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	issuer           = "octopus-api-server"
	ttlSeconds       = 604800 //7 days
	actionQuery      = "action"
	loginActionName  = "login"
	logoutActionName = "logout"
)

// NewHandler creates a new AuthHandler
func NewAuthHandler(ctx context.Context, host string, client *kubernetes.Clientset) *Handler {
	return &Handler{
		context:   ctx,
		Host:      host,
		clientset: client,
	}
}

type Handler struct {
	context   context.Context
	Host      string
	clientset *kubernetes.Clientset
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	action := strings.ToLower(req.URL.Query().Get(actionQuery))

	tokenAuthValue := GetTokenAuthFromRequest(req)
	if tokenAuthValue == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte("must authenticate"))
		return
	}

	// manage logout action
	if action == logoutActionName {
		err := h.removeTokenSecret(tokenAuthValue)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(err.Error()))
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("success logout"))
		return
	}

	tokenName, tokenKey := SplitTokenParts(tokenAuthValue)
	if tokenName != "admin" || tokenKey == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte("must authenticate"))
		return
	}

	_, err := NormalizeAndValidateTokenForUser(h.Host, tokenKey, tokenName)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(err.Error()))
		return
	}

	if action == loginActionName {
		// create a token
		tokenString, err := h.createToken(tokenName)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("Sorry, error while Signing Token!"))
			logrus.Printf("Token Signing error: %v\n", err)
			return
		}

		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(tokenString))
		return
	}

	writer.WriteHeader(http.StatusBadRequest)
	writer.Write([]byte("action unknown or not provided"))
	return
}

func (h *Handler) removeTokenSecret(token string) error {
	name, err := GetJWTSecretTokenName(token)
	if err != nil {
		return err
	}
	return h.clientset.CoreV1().Secrets(TokenNamespace).Delete(h.context, name, metav1.DeleteOptions{})
}

func (h *Handler) createToken(name string) (string, error) {
	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + ttlSeconds,
		Issuer:    issuer,
		Subject:   name,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// key need to be `your-256-bit-secret`
	key, err := util.GenerateRandomKey()
	if err != nil {
		return "", err
	}

	strToken, err := token.SignedString([]byte(key))
	if err != nil {
		logrus.Errorf("failed to sign a JWT token, error: %s", err.Error())
		return "", err
	}

	secretToken := TokenSecretData{
		Issuer:    issuer,
		ExpiresAt: time.Now().Add(ttlSeconds * time.Second).Format(time.RFC3339),
		IssuedAt:  time.Now().Format(time.RFC3339),
		Subject:   name,
		Key:       key,
	}
	secret, err := createTokenSecret(strToken, secretToken)
	if err != nil {
		logrus.Errorf("failed to generate a JWT secret token, error: %s", err.Error())
		return "", err
	}
	_, err = h.clientset.CoreV1().Secrets(TokenNamespace).Create(h.context, &secret, metav1.CreateOptions{})
	if err != nil {
		logrus.Errorf("failed to create token secret, error: %s", err.Error())
		return "", err
	}

	return strToken, nil
}
