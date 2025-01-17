package handler

import (
	"watt-flow/model"
	"watt-flow/service"
	"watt-flow/util"

	"github.com/gin-gonic/gin"
)

type PricelistHandler struct {
	service service.IPricelistService
	logger  util.Logger
}

func NewPricelistHandler(service service.IPricelistService, logger util.Logger) *PricelistHandler {
	return &PricelistHandler{
		service: service,
		logger:  logger,
	}
}

func (h PricelistHandler) CreatePricelist(c *gin.Context) {
	var pricelist model.Pricelist
	if err := c.BindJSON(&pricelist); err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid pricelist data"})
		return
	}

	data, err := h.service.CreatePricelist(&pricelist)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to create pricelist"})
		return
	}
	c.JSON(201, gin.H{"data": data})
}
