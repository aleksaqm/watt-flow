package handler

import (
	"strconv"
	"time"
	"watt-flow/service"
	"watt-flow/util"

	"github.com/gin-gonic/gin"
)

type ElectricityConsumptionHandler struct {
	service service.IElectricityConsumptionService
	logger  util.Logger
}

func NewElectricityConsumptionHandler(service service.IElectricityConsumptionService, logger util.Logger) *ElectricityConsumptionHandler {
	return &ElectricityConsumptionHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ElectricityConsumptionHandler) GetMonthlyConsumption(c *gin.Context) {
	householdId := c.Param("id")

	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
	monthStr := c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month())))

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid year parameter"})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid month parameter"})
		return
	}

	data, err := h.service.GetMonthlyConsumption(householdId, year, month)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to get consumption data"})
		return
	}

	c.JSON(200, gin.H{"data": data})
}

func (h *ElectricityConsumptionHandler) Get12MonthsConsumption(c *gin.Context) {
	householdId := c.Param("id")

	endYearStr := c.DefaultQuery("endYear", strconv.Itoa(time.Now().Year()))
	endMonthStr := c.DefaultQuery("endMonth", strconv.Itoa(int(time.Now().Month())))

	endYear, err := strconv.Atoi(endYearStr)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid endYear parameter"})
		return
	}

	endMonth, err := strconv.Atoi(endMonthStr)
	if err != nil || endMonth < 1 || endMonth > 12 {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid endMonth parameter"})
		return
	}

	data, err := h.service.Get12MonthsConsumption(householdId, endYear, endMonth)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to get 12 months consumption data"})
		return
	}

	c.JSON(200, gin.H{"data": data.Data})
}

func (h *ElectricityConsumptionHandler) GetDailyConsumption(c *gin.Context) {
	householdId := c.Param("id")

	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if yearStr == "" || monthStr == "" {
		c.JSON(400, gin.H{"error": "Year and month parameters are required"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid year parameter"})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid month parameter"})
		return
	}

	if month < 1 || month > 12 {
		c.JSON(400, gin.H{"error": "Month must be between 1 and 12"})
		return
	}

	data, err := h.service.GetDailyConsumption(householdId, year, month)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to get daily consumption data"})
		return
	}

	c.JSON(200, gin.H{"data": data.Data})
}
