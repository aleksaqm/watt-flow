package handler

import (
	"watt-flow/service"
	"watt-flow/util"

	"github.com/gin-gonic/gin"
)

type CityHandler struct {
	service service.ICityService
	logger  util.Logger
}

func NewCityHandler(service service.ICityService, logger util.Logger) *CityHandler {
	return &CityHandler{
		service: service,
		logger:  logger,
	}
}

func (h *CityHandler) GetAllCities(c *gin.Context) {
	cities, err := h.service.GetAllCities()
	if err != nil {
		h.logger.Error("Error fetching city list", err)
		c.JSON(500, gin.H{"error": "Failed to retrieve cities"})
		return
	}
	c.JSON(200, gin.H{"data": cities})
}
