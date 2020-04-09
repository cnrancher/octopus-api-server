package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

const (
	issuer     = "Edge-API"
	ttlSeconds = 604800
)

type CustomClaimsExample struct {
	*jwt.StandardClaims
	TokenType string
	CustomerInfo
}

// Define some custom types were going to use within our tokens
type CustomerInfo struct {
	Name string
	Kind string
}

// NewLoginHandler creates a new AuthLoginHandler
func NewLoginHandler(host string) *LoginHandler {
	return &LoginHandler{
		Host: host,
	}
}

// AuthLoginHandler lists available machine types in GKE
type LoginHandler struct {
	Host string
}

func (l *LoginHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tokenAuthValue := GetTokenAuthFromRequest(req)
	if tokenAuthValue == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte("must authenticate"))
		return
	}

	tokenName, tokenKey := SplitTokenParts(tokenAuthValue)
	if tokenName != "admin" || tokenKey == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte("must authenticate"))
		return
	}

	_, err := NormalizeAndValidateTokenForUser(l.Host, tokenKey, tokenName)
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(err.Error()))
		return
	}

	// create a token
	tokenString, err := createToken(tokenKey)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Sorry, error while Signing Token!"))
		logrus.Printf("Token Signing error: %v\n", err)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(tokenString))
}

func createToken(key string) (string, error) {
	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + ttlSeconds,
		Issuer:    issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// key need to be `your-256-bit-secret`
	return token.SignedString([]byte(key))
}
