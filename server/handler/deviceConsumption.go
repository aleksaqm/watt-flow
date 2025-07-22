package handler

import (
	"watt-flow/dto"
	"watt-flow/service"
	"watt-flow/util"

	"github.com/gin-gonic/gin"
)

type DeviceConsumptionHandler struct {
	service service.IDeviceConsumptionService
	logger  util.Logger
}

func (h DeviceConsumptionHandler) QueryConsumption(c *gin.Context) {
	var queryParams dto.FluxQueryConsumptionDto
	if err := c.BindJSON(&queryParams); err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid device consumption query params"})
		return
	}
	data, err := h.service.QueryConsumption(queryParams)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to query device consumption"})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

func NewDeviceConsumptionHandler(deviceConsumptionService service.IDeviceConsumptionService, logger util.Logger) *DeviceConsumptionHandler {
	return &DeviceConsumptionHandler{
		service: deviceConsumptionService,
		logger:  logger,
	}
}
