package handler

import (
	"strconv"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/service"
	"watt-flow/util"

	"github.com/gin-gonic/gin"
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

func (h HouseholdHandler) Query(c *gin.Context) {
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

	var searchParams dto.HouseholdSearchParams
	if err := c.BindJSON(&searchParams); err != nil {
		c.JSON(400, gin.H{"error": "Invalid search parameter"})
		return
	}

	params := dto.HouseholdQueryParams{
		Page:      pageInt,
		PageSize:  pageSizeInt,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Search:    searchParams,
	}
	h.logger.Info(params)
	households, total, err := h.service.Query(&params)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"households": households, "total": total})
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

func (h HouseholdHandler) CreateOwnershipRequest(c *gin.Context) {
	//validation to check does household have owner
	var ownershipRequest dto.OwnershipRequestDto
	if err := c.BindJSON(&ownershipRequest); err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid ownership request"})
		return
	}
	request, err := h.service.CreateOwnershipRequest(ownershipRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}
	c.JSON(201, gin.H{"data": request})
}

func (h HouseholdHandler) GetOwnershipRequestsForUser(c *gin.Context) {
	id := c.Param("id")
	requestId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid household ID"})
		return
	}
	requests, err := h.service.GetOwnersRequests(requestId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to load owners requests"})
		return
	}
	c.JSON(200, gin.H{"data": requests})
}

func NewHouseholdHandler(householdService service.IHouseholdService, logger util.Logger) *HouseholdHandler {
	return &HouseholdHandler{
		service: householdService,
		logger:  logger,
	}
}
