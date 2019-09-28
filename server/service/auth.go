package service

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AuthService providers authentication service
type AuthService interface {
	Auth(string, string) (string, error)
}

func NewAuthService(key []byte, clients map[string]string) AuthService {
	service := authService{key, clients}

	return service
}

type authService struct {
	key    []byte
	clients map[string]string
}

type CustomClaims struct {
	ClientID string `json:"clientId"`
	jwt.StandardClaims
}

func CustomClaimsFactory() jwt.Claims {
	return &CustomClaims{}
}
const expiration  =  14400

func generateToken(signingKey []byte, clientID string) (string, error) {
	claims := CustomClaims{
		clientID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * expiration).Unix(),
			//IssuedAt: jwt.TimeFunc().Unix(),
			Issuer:    "system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 生成Token
	return token.SignedString(signingKey)
}

func (as authService) Auth(clientID string, clientSecret string) (string, error) {
	if as.clients[clientID] == clientSecret {
		signed, err := generateToken(as.key, clientID)
		if err != nil {
			return "", errors.New(err.Error())
		}

		return signed, nil
	}

	return "", ErrAuth
}

// ErrAuth is returned when credentials are incorrect
var ErrAuth = errors.New("Incorrect credentials")