package handler

import (
	"net/http"
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

func (h PricelistHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.logger.Error("Invalid pricelist ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pricelist ID"})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		h.logger.Error("Pricelist not found:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Pricelist not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pricelist deleted successfully"})
}
