package handler

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/service"
	"watt-flow/util"
)

type HouseholdHandler struct {
	service service.IHouseholdService
	logger  util.Logger
}

func (h HouseholdHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	householdId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid household ID"})
		return
	}
	data, err := h.service.FindById(householdId)
	if err != nil {
		h.logger.Error(err)
		c.JSON(404, gin.H{"error": "Household not found"})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

func (h HouseholdHandler) Create(c *gin.Context) {
	var household dto.CreateHouseholdDto
	if err := c.BindJSON(&household); err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid household data"})
		return
	}
	data, err := h.service.Create(&household)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to create household"})
		return
	}
	c.JSON(201, gin.H{"data": data})
}

func (h HouseholdHandler) Update(c *gin.Context) {
	var household model.Household
	if err := c.BindJSON(&household); err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid household data"})
		return
	}
	data, err := h.service.Update(&household)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to update household"})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

func (h HouseholdHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	householdId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid household ID"})
		return
	}
	err = h.service.Delete(householdId)
	if err != nil {
		h.logger.Error(err)
		c.JSON(404, gin.H{"error": "Household not found"})
		return
	}
	c.JSON(200, gin.H{"message": "Household deleted"})
}

func (h HouseholdHandler) FindByStatus(c *gin.Context) {
	status := c.Query("status")
	parsedStatus, err := model.ParseHouseholdStatus(status)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid household status"})
		return
	}
	households, err := h.service.FindByStatus(parsedStatus)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to retrieve households by status"})
		return
	}
	c.JSON(200, gin.H{"data": households})
}

func NewHouseholdHandler(householdService service.IHouseholdService, logger util.Logger) *HouseholdHandler {
	return &HouseholdHandler{
		service: householdService,
		logger:  logger,
	}
}
