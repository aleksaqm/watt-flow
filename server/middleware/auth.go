package middleware

import (
	"watt-flow/service"
	"watt-flow/util"
)

type AuthMiddleware struct {
	service service.AuthService
	logger  util.Logger
}

func NewAuthMiddleware(service service.AuthService, logger util.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		service: service,
		logger:  logger,
	}
}
