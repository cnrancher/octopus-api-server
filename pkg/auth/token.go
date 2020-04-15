package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/cnrancher/edge-api-server/pkg/util"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	CookieName      = "R_SESS"
	AuthHeaderName  = "Authorization"
	AuthValuePrefix = "Bearer"
	BasicAuthPrefix = "Basic"

	usernameLabel      = "authn.management.edge.io/token-username"
	edgeApiLabel       = "authn.management.edge.io/edge-api"
	tokenNamespace     = "kube-system"
	nameLength         = 8
	tokenSecretKeyName = "Key"
)

type TokenSecretData struct {
	Issuer    string `json:"issuer,omitempty"`
	ExpiresAt string `json:"expiresAt,omitempty"`
	IssuedAt  string `json:"issuedAt,omitempty"`
	Subject   string `json:"subject,omitempty"`
	Key       string `json:"key,omitempty"`
}

func SplitTokenParts(tokenID string) (string, string) {
	parts := strings.Split(tokenID, ":")
	if len(parts) != 2 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

func GetTokenAuthFromRequest(req *http.Request) string {
	var tokenAuthValue string
	authHeader := req.Header.Get(AuthHeaderName)
	authHeader = strings.TrimSpace(authHeader)

	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if strings.EqualFold(parts[0], AuthValuePrefix) {
			if len(parts) > 1 {
				tokenAuthValue = strings.TrimSpace(parts[1])
			}
		} else if strings.EqualFold(parts[0], BasicAuthPrefix) {
			if len(parts) > 1 {
				base64Value := strings.TrimSpace(parts[1])
				data, err := base64.URLEncoding.DecodeString(base64Value)
				if err != nil {
					logrus.Errorf("Error %v parsing %v header", err, AuthHeaderName)
				} else {
					tokenAuthValue = string(data)
				}
			}
		}
	} else {
		cookie, err := req.Cookie(CookieName)
		if err == nil {
			tokenAuthValue = cookie.Value
		}
	}
	return tokenAuthValue
}

func GetJWTSecretTokenName(token string) (string, error) {
	var name = ""
	parts, err := SplitJWTTokenParts(token)
	if err != nil {
		return name, err
	}
	name = strings.Trim(strings.ToLower(parts[2]), "_")
	name = name[len(name)-nameLength:]
	return fmt.Sprintf("jwt-%s-secret", name), nil
}

func SplitJWTTokenParts(token string) ([]string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return parts, fmt.Errorf("invalid token, expected length 3 but only got %d", len(parts))
	}
	return parts, nil
}

func createTokenSecret(token string, secretToken TokenSecretData) (corev1.Secret, error) {
	name, err := GetJWTSecretTokenName(token)
	if err != nil {
		return corev1.Secret{}, err
	}

	strData := util.StructToStrMap(&secretToken, 4)
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				usernameLabel: secretToken.Subject,
				edgeApiLabel:  "true",
			},
		},
		StringData: strData,
	}, nil
}
