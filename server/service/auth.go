package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"watt-flow/config"
	"watt-flow/model"
	"watt-flow/util"
)

type AuthService struct {
	logger util.Logger
}

func NewAuthService(logger util.Logger) *AuthService {
	return &AuthService{
		logger: logger,
	}
}

func (s *AuthService) Authorize(tokenString string) (bool, jwt.MapClaims, error) {
	env := config.Init()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(env.JWTSecret), nil
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
	env := config.Init()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role.String(),
	})

	tokenString, err := token.SignedString([]byte(env.JWTSecret))
	if err != nil {
		s.logger.Error(err)
	}
	return tokenString
}

//jwt.MapClaims, error
