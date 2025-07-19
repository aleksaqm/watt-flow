package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"watt-flow/dto"
	"watt-flow/service"
	"watt-flow/util"
)

type HouseholdAccessHandler struct {
	service service.IHouseholdAccessService
	logger  util.Logger
}

func NewHouseholdAccessHandler(service service.IHouseholdAccessService, logger util.Logger) *HouseholdAccessHandler {
	return &HouseholdAccessHandler{
		service: service,
		logger:  logger,
	}
}

func (h *HouseholdAccessHandler) GrantAccess(c *gin.Context) {
	householdIdStr := c.Param("id")
	householdId, err := strconv.ParseUint(householdIdStr, 10, 64)
	if err != nil {
		h.logger.Error("Invalid household ID format", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid household ID format"})
		return
	}

	var req dto.GrantAccessRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body for granting access", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: 'userId' is required"})
		return
	}

	currentUserID := uint64(14)

	//currentUserID, exists := c.Get("userId")
	//if !exists {
	//	h.logger.Error("User ID not found in context")
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	//	return
	//}

	err = h.service.GrantAccess(householdId, req.UserID, currentUserID)
	if err != nil {
		h.logger.Error("Failed to grant access", err)
		if err.Error() == "forbidden: only the owner can grant access" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "household not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Access granted successfully"})
}

func (h *HouseholdAccessHandler) RevokeAccess(c *gin.Context) {
	householdIdStr := c.Param("householdId")
	householdId, err := strconv.ParseUint(householdIdStr, 10, 64)
	if err != nil {
		h.logger.Error("Invalid household ID format", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid household ID format"})
		return
	}

	userIdToRevokeStr := c.Param("userId")
	userIdToRevoke, err := strconv.ParseUint(userIdToRevokeStr, 10, 64)
	if err != nil {
		h.logger.Error("Invalid user ID format for revoking", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	currentUserID := uint64(14)

	//currentUserID, exists := c.Get("userId")
	//if !exists {
	//	h.logger.Error("User ID not found in context")
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	//	return
	//}

	err = h.service.RevokeAccess(householdId, userIdToRevoke, currentUserID)
	if err != nil {
		h.logger.Error("Failed to revoke access", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Access for this user does not exist"})
			return
		}
		if err.Error() == "forbidden: only the owner can revoke access" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke access"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *HouseholdAccessHandler) ListAccess(c *gin.Context) {

	householdIdStr := c.Param("householdId")
	householdId, err := strconv.ParseUint(householdIdStr, 10, 64)

	if err != nil {
		h.logger.Error("Invalid household ID format", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid household ID format"})
		return
	}

	currentUserID := uint64(14)
	accesses, err := h.service.GetUsersWithAccess(householdId, currentUserID)
	h.logger.Info("Accesses: ", accesses)
	if err != nil {
		h.logger.Error("Failed to get access list", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "household not found"})
			return
		}
		if err.Error() == "forbidden: only the owner can view access list" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access list"})
		return
	}

	c.JSON(200, gin.H{"data": accesses})
}
