package handler

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"net/http"
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

func (h *BillHandler) SearchBills(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	sortBy := c.DefaultQuery("sortBy", "billing_date")
	sortOrder := c.DefaultQuery("sortOrder", "asc")
	search := c.DefaultQuery("search", "{}")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		h.logger.Error("Invalid page parameter", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		h.logger.Error("Invalid pageSize parameter", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageSize parameter"})
		return
	}

	var searchParams dto.BillSearchParams
	if search != "" && search != "{}" {
		if err := json.Unmarshal([]byte(search), &searchParams); err != nil {
			h.logger.Error("Invalid search parameter", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid search parameter format"})
			return
		}
	}

	claims, exists := c.Get("claims")
	if !exists {
		h.logger.Error("Claims not found in context, middleware might be missing")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User claims not found"})
		c.Abort()
		return
	}

	claimsMap, ok := claims.(jwt.MapClaims)
	if !ok {
		h.logger.Error("Invalid claims format", zap.Any("claims", claims))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user claims format"})
		c.Abort()
		return
	}

	loggedInUserID, ok := claimsMap["id"]
	if !ok {
		h.logger.Error("ID not found in claims")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in claims"})
		c.Abort()
		return
	}

	userIDFloat, ok := loggedInUserID.(float64)
	if !ok {
		h.logger.Error("Failed to cast user ID to float64")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
		c.Abort()
		return
	}

	searchParams.UserID = uint64(userIDFloat)

	params := dto.BillQueryParams{
		Page:      pageInt,
		PageSize:  pageSizeInt,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Search:    searchParams,
	}

	h.logger.Info("Search bills params", params)

	bills, total, err := h.service.SearchBills(&params)
	if err != nil {
		h.logger.Error("Failed to query bills", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bills"})
		return
	}

	var billsDto []dto.BillResponseDto
	for _, bill := range bills {
		billsDto = append(billsDto, dto.BillResponseDto{
			ID:          bill.ID,
			IssueDate:   bill.IssueDate,
			BillingDate: bill.BillingDate,
			SpentPower:  bill.SpentPower,
			Price:       bill.Price,
			Status:      bill.Status,
			Pricelist: dto.PricelistDto{
				ID:           bill.Pricelist.ID,
				ValidFrom:    bill.Pricelist.ValidFrom,
				BlueZone:     bill.Pricelist.BlueZone,
				RedZone:      bill.Pricelist.RedZone,
				GreenZone:    bill.Pricelist.GreenZone,
				BillingPower: bill.Pricelist.BillingPower,
				Tax:          bill.Pricelist.Tax,
			},
			Owner: dto.OwnerDto{
				ID:       bill.OwnerID,
				Username: bill.Owner.Username,
				Email:    bill.Owner.Email,
			},
			Household: dto.HouseholdResultDto{
				Id:              bill.HouseholdID,
				CadastralNumber: bill.Household.CadastralNumber,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"bills": billsDto,
		"total": total,
	})
}

func (h *BillHandler) GetBill(c *gin.Context) {

	id := c.Query("id")

	if id == "" {
		c.JSON(400, gin.H{"error": "ID missing in request"})
		return
	}

	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid id"})
		return
	}
	h.logger.Info("Get bill", idInt)
	bill, err := h.service.FindById(idInt)
	if err != nil {
		h.logger.Error(err)
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	billDto := dto.BillResponseDto{
		ID:          bill.ID,
		IssueDate:   bill.IssueDate,
		BillingDate: bill.BillingDate,
		SpentPower:  bill.SpentPower,
		Price:       bill.Price,
		Status:      bill.Status,
		Pricelist: dto.PricelistDto{
			ID:           bill.Pricelist.ID,
			ValidFrom:    bill.Pricelist.ValidFrom,
			BlueZone:     bill.Pricelist.BlueZone,
			RedZone:      bill.Pricelist.RedZone,
			GreenZone:    bill.Pricelist.GreenZone,
			BillingPower: bill.Pricelist.BillingPower,
			Tax:          bill.Pricelist.Tax,
		},
		Owner: dto.OwnerDto{
			ID:       bill.OwnerID,
			Username: bill.Owner.Username,
			Email:    bill.Owner.Email,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"bill": billDto,
	})
}

func (h *BillHandler) PayBill(c *gin.Context) {
	billIdStr := c.Param("id")

	billId, err := strconv.ParseUint(billIdStr, 10, 64)
	if err != nil {
		h.logger.Error("Invalid bill ID format", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bill ID format"})
		return
	}

	claims, exists := c.Get("claims")
	if !exists {
		h.logger.Error("Claims not found in context, middleware might be missing")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User claims not found"})
		c.Abort()
		return
	}

	claimsMap, ok := claims.(jwt.MapClaims)
	if !ok {
		h.logger.Error("Invalid claims format", zap.Any("claims", claims))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user claims format"})
		c.Abort()
		return
	}

	loggedInUserID, ok := claimsMap["id"]
	if !ok {
		h.logger.Error("ID not found in claims")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in claims"})
		c.Abort()
		return
	}

	userIDFloat, ok := loggedInUserID.(float64)
	if !ok {
		h.logger.Error("Failed to cast user ID to float64")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
		c.Abort()
		return
	}

	err = h.service.PayBill(billId, uint64(userIDFloat))
	if err != nil {
		h.logger.Error("Failed to process payment", err)
		if err.Error() == "bill not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "forbidden: you are not authorized to pay this bill" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "this bill has already been paid" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while processing the payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment successful."})
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

	data, err := h.service.GenerateMonthlyBill(year, month)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to send bill for specified month!"})
		return
	}
	c.JSON(201, gin.H{"data": data})
}

func (h *BillHandler) InitiateBilling(c *gin.Context) {
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

	data, err := h.service.InitiateBillingOffload(year, month)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to send bills for specified month!"})
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
