package handler

import (
	"strconv"
	"watt-flow/model"
	"watt-flow/service"
	"watt-flow/util"

	"github.com/gin-gonic/gin"
)

type DeviceStatusHandler struct {
	service service.IDeviceStatusService
	logger  util.Logger
}

func (h DeviceStatusHandler) GetByAddress(c *gin.Context) {
	address := c.Param("address")
	data, err := h.service.FindByAddress(address)
	if err != nil {
		h.logger.Error(err)
		c.JSON(404, gin.H{"error": "Device status not found"})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

func (h DeviceStatusHandler) Query(c *gin.Context) {
	err := h.service.Query()
	if err != nil {
		h.logger.Error(err)
		c.JSON(404, gin.H{"error": "Device status not found"})
		return
	}
	c.JSON(200, gin.H{})
}

func (h DeviceStatusHandler) GetByHouseholdID(c *gin.Context) {
	id := c.Param("household_id")
	householdID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid household ID"})
		return
	}
	data, err := h.service.FindByHouseholdID(householdID)
	if err != nil {
		h.logger.Error(err)
		c.JSON(404, gin.H{"error": "Device status not found for this household"})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

func (h DeviceStatusHandler) Create(c *gin.Context) {
	var deviceStatus model.DeviceStatus
	if err := c.BindJSON(&deviceStatus); err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid device status data"})
		return
	}
	data, err := h.service.Create(&deviceStatus)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to create device status"})
		return
	}
	c.JSON(201, gin.H{"data": data})
}

func (h DeviceStatusHandler) Update(c *gin.Context) {
	var deviceStatus model.DeviceStatus
	if err := c.BindJSON(&deviceStatus); err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid device status data"})
		return
	}
	data, err := h.service.Update(&deviceStatus)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to update device status"})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

func (h DeviceStatusHandler) Delete(c *gin.Context) {
	address := c.Param("address")
	err := h.service.Delete(address)
	if err != nil {
		h.logger.Error(err)
		c.JSON(404, gin.H{"error": "Device status not found"})
		return
	}
	c.JSON(200, gin.H{"message": "Device status deleted"})
}

func NewDeviceStatusHandler(deviceStatusService service.IDeviceStatusService, logger util.Logger) *DeviceStatusHandler {
	return &DeviceStatusHandler{
		service: deviceStatusService,
		logger:  logger,
	}
}
