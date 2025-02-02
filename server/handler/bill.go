package handler

import (
	"fmt"
	"strconv"
	"watt-flow/dto"
	"watt-flow/service"
	"watt-flow/util"

	"github.com/gin-gonic/gin"
)

type BillHandler struct {
	service service.IBillService
	logger  util.Logger
}

func NewBillHandler(service service.IBillService, logger util.Logger) *BillHandler {
	return &BillHandler{
		service: service,
		logger:  logger,
	}
}

func (h BillHandler) Query(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	sortBy := c.DefaultQuery("sortBy", "date")
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

	params := dto.MonthlyBillQueryParams{
		Page:      pageInt,
		PageSize:  pageSizeInt,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
	bills, total, err := h.service.QueryMonthly(&params)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"bills": bills, "total": total})
}

func (h *BillHandler) GetUnsentMonthlyBills(c *gin.Context) {
	cities, err := h.service.GetUnsentMonthlyBills()
	if err != nil {
		h.logger.Error("Error fetching unsent bills", err)
		c.JSON(500, gin.H{"error": "Failed to retrieve bills"})
		return
	}
	c.JSON(200, gin.H{"data": cities})
}

func (h *BillHandler) SendBill(c *gin.Context) {
	var bill dto.SendBillDto
	if err := c.BindJSON(&bill); err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid bill month"})
		return
	}
	var year, month int
	_, err := fmt.Sscanf(bill.Month, "%d-%02d", &year, &month)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to send bill for specified month!"})
		return
	}

	data, err := h.service.SendBill(year, month)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to send bill for specified month!"})
		return
	}
	c.JSON(201, gin.H{"data": data})
}

// func (h BillHandler) CreateBill(c *gin.Context) {
// 	var pricelist dto.NewBill
// 	if err := c.BindJSON(&pricelist); err != nil {
// 		h.logger.Error(err)
// 		c.JSON(400, gin.H{"error": "Invalid pricelist data"})
// 		return
// 	}
//
// 	data, err := h.service.CreateBill(&pricelist)
// 	if err != nil {
// 		h.logger.Error(err)
// 		c.JSON(500, gin.H{"error": "Failed to create pricelist!"})
// 		return
// 	}
// 	c.JSON(201, gin.H{"data": data})
// }
//
// func (h BillHandler) Delete(c *gin.Context) {
// 	idParam := c.Param("id")
// 	id, err := strconv.ParseUint(idParam, 10, 64)
// 	if err != nil {
// 		h.logger.Error("Invalid pricelist ID:", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pricelist ID"})
// 		return
// 	}
//
// 	err = h.service.Delete(id)
// 	if err != nil {
// 		h.logger.Error("Bill not found:", err)
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{"message": "Bill deleted successfully"})
// }
