package handler

import (
	"strconv"
	"watt-flow/dto"
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

func (h PricelistHandler) Query(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	sortBy := c.DefaultQuery("sortBy", "city")
	sortOrder := c.DefaultQuery("sortOrder", "asc")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid page parameter"})
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid pageSize parameter"})
		return
	}

	params := dto.PricelistQueryParams{
		Page:      pageInt,
		PageSize:  pageSizeInt,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
	pricelists, total, err := h.service.Query(&params)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"pricelists": pricelists, "total": total})
}

func (h PricelistHandler) CreatePricelist(c *gin.Context) {
	var pricelist dto.NewPricelist
	if err := c.BindJSON(&pricelist); err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid pricelist data"})
		return
	}

	data, err := h.service.CreatePricelist(&pricelist)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to create pricelist!"})
		return
	}
	c.JSON(201, gin.H{"data": data})
}
