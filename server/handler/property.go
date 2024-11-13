package handler

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"watt-flow/model"
	"watt-flow/service"
	"watt-flow/util"
)

type PropertyHandler struct {
	service service.IPropertyService
	logger  util.Logger
}

func (p PropertyHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	propertyId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		p.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid property ID"})
		return
	}
	data, err := p.service.FindById(propertyId)
	p.logger.Info(data)
	if err != nil {
		p.logger.Error(err)
		c.JSON(404, gin.H{"error": "Property not found"})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

func (p PropertyHandler) Create(c *gin.Context) {
	var property model.Property
	if err := c.BindJSON(&property); err != nil {
		p.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid property data"})
		return
	}
	p.logger.Info(property)
	data, err := p.service.Create(&property)
	if err != nil {
		p.logger.Error(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"data": data})
}

func (p PropertyHandler) Update(c *gin.Context) {
	var property model.Property
	if err := c.BindJSON(&property); err != nil {
		p.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid property data"})
		return
	}
	data, err := p.service.Update(&property)
	if err != nil {
		p.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to update property"})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

func (p PropertyHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	propertyId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		p.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid property ID"})
		return
	}
	err = p.service.Delete(propertyId)
	if err != nil {
		p.logger.Error(err)
		c.JSON(404, gin.H{"error": "Property not found"})
		return
	}
	c.JSON(200, gin.H{"message": "Property deleted"})
}

func (p PropertyHandler) FindByStatus(c *gin.Context) {
	status := c.Query("status")
	parsedStatus, err := model.ParsePropertyStatus(status)
	if err != nil {
		p.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid property status"})
		return
	}
	properties, err := p.service.FindByStatus(parsedStatus)
	if err != nil {
		p.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to retrieve properties by status"})
		return
	}
	c.JSON(200, gin.H{"data": properties})
}

func NewPropertyHandler(propertyService service.IPropertyService, logger util.Logger) *PropertyHandler {
	return &PropertyHandler{
		service: propertyService,
		logger:  logger,
	}
}
