package service

import (
	"errors"
	"log"
	"time"
	"watt-flow/config"
	"watt-flow/model"
	"watt-flow/util"

	"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
	logger util.Logger
	env    *config.Environment
}

func NewAuthService(logger util.Logger, env *config.Environment) *AuthService {
	return &AuthService{
		logger: logger,
		env:    env,
	}
}

func (s *AuthService) Authorize(tokenString string) (bool, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.env.JWTSecret), nil
	})
	if token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		return true, claims, nil
	}
	var ve *jwt.ValidationError
	if errors.As(err, &ve) {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			log.Println("malformed token")
			return false, nil, errors.New("malformed token")
		}
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return false, nil, errors.New("token expired or not active")
		}
	}
	return false, nil, errors.New("invalid token")
}

func (s *AuthService) CreateToken(user *model.User) string {
	expirationTime := time.Now().Add(24 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role.RoleToString(),
		"exp":      expirationTime,
	})
	tokenString, err := token.SignedString([]byte(s.env.JWTSecret))
	if err != nil {
		s.logger.Error(err)
	}
	return tokenString
}

func (s *AuthService) CreateActivationToken(user *model.User) string {
	expirationTime := time.Now().Add(time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    user.Email,
		"username": user.Username,
		"exp":      expirationTime,
	})
	tokenString, err := token.SignedString([]byte(s.env.JWTSecret))
	if err != nil {
		s.logger.Error(err)
	}
	return tokenString
}

// jwt.MapClaims, error
