package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"watt-flow/model"
	"watt-flow/service"
	"watt-flow/util"
)

type AddressHandler struct {
	service service.IAddressService
	logger  util.Logger
}

func NewAddressHandler(addressService service.IAddressService, logger util.Logger) *AddressHandler {
	return &AddressHandler{
		service: addressService,
		logger:  logger,
	}
}

func (h *AddressHandler) Create(c *gin.Context) {
	var address model.Address
	if err := c.ShouldBindJSON(&address); err != nil {
		h.logger.Error("Invalid address data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address data"})
		return
	}

	createdAddress, err := h.service.Create(&address)
	if err != nil {
		h.logger.Error("Failed to create address:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": createdAddress})
}

func (h *AddressHandler) FindById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.logger.Error("Invalid address ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	address, err := h.service.FindById(id)
	if err != nil {
		h.logger.Error("Address not found:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": address})
}

func (h *AddressHandler) FindAll(c *gin.Context) {
	addresses, err := h.service.FindAll()
	if err != nil {
		h.logger.Error("Failed to retrieve addresses:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve addresses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": addresses})
}

func (h *AddressHandler) Update(c *gin.Context) {
	var address model.Address
	if err := c.ShouldBindJSON(&address); err != nil {
		h.logger.Error("Invalid address data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address data"})
		return
	}

	updatedAddress, err := h.service.Update(&address)
	if err != nil {
		h.logger.Error("Failed to update address:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedAddress})
}

func (h *AddressHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.logger.Error("Invalid address ID:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		h.logger.Error("Address not found:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
}
