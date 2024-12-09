package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"watt-flow/dto"
	"watt-flow/service"
	"watt-flow/util"

	"github.com/gin-gonic/gin"
)

type OwnershipHandler struct {
	service service.IOwnershipService
	logger  util.Logger
}

func (h OwnershipHandler) CreateOwnershipRequest(c *gin.Context) {
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

func (h OwnershipHandler) GetOwnershipRequestsForUser(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	sortBy := c.DefaultQuery("sortBy", "city")
	sortOrder := c.DefaultQuery("sortOrder", "asc")
	search := c.DefaultQuery("search", "")

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

	var searchParams dto.OwnershipSearchParams
	if search != "" {
		if err := json.Unmarshal([]byte(search), &searchParams); err != nil {
			c.JSON(400, gin.H{"error": "Invalid search parameter"})
			return
		}
	}
	params := dto.OwnershipQueryParams{
		Page:      pageInt,
		PageSize:  pageSizeInt,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Search:    searchParams,
	}
	h.logger.Info(params)
	id := c.Param("id")
	requestId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid household ID"})
		return
	}
	requests, total, err := h.service.GetOwnersRequests(requestId, &params)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to load owners requests"})
		return
	}
	c.JSON(200, gin.H{"requests": requests, "total": total})
}

func (h OwnershipHandler) GetPendingRequests(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	sortBy := c.DefaultQuery("sortBy", "city")
	sortOrder := c.DefaultQuery("sortOrder", "asc")
	search := c.DefaultQuery("search", "")

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

	var searchParams dto.OwnershipSearchParams
	if search != "" {
		if err := json.Unmarshal([]byte(search), &searchParams); err != nil {
			c.JSON(400, gin.H{"error": "Invalid search parameter"})
			return
		}
	}
	params := dto.OwnershipQueryParams{
		Page:      pageInt,
		PageSize:  pageSizeInt,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Search:    searchParams,
	}
	h.logger.Info(params)
	requests, total, err := h.service.GetPendingRequests(&params)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to load owners requests"})
		return
	}
	c.JSON(200, gin.H{"requests": requests, "total": total})
}

func (h OwnershipHandler) AcceptOwnershipRequest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid property ID"})
		return
	}
	err = h.service.AcceptOwnershipRequest(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to accept ownership request"})
		return
	}
	c.JSON(200, gin.H{"message": "Owner ship request accepted successfully"})
}

func NewOwnershipHandler(ownershipService service.IOwnershipService, logger util.Logger) *OwnershipHandler {
	return &OwnershipHandler{
		service: ownershipService,
		logger:  logger,
	}
}
