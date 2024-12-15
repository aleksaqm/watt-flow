package handler

import (
	"net/http"
	"strings"
	"time"
	"watt-flow/dto"
	"watt-flow/service"
	"watt-flow/util"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type MeetingHandler struct {
	service service.IMeetingService
	logger  util.Logger
}

func NewMeetingHandler(service service.IMeetingService, logger util.Logger) *MeetingHandler {
	return &MeetingHandler{
		service: service,
		logger:  logger,
	}
}

func (h MeetingHandler) GetSlotById(c *gin.Context) {
	id := strings.TrimSpace(c.Query("date"))
	h.logger.Info("Parsing time: ", id)

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not valid id parameter"})
	}

	timeslotDate, err := time.Parse("2006-01-02", id)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid date"})
		return
	}

	data, err := h.service.FindByDate(datatypes.Date(timeslotDate))
	if err != nil {
		h.logger.Error(err)
		c.JSON(404, gin.H{"error": " timeslot not found"})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

func (h MeetingHandler) CreateSlot(c *gin.Context) {
	var timeslotDto dto.TimeSlotDto
	if err := c.BindJSON(&timeslotDto); err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid timeslot data"})
		return
	}

	data, err := h.service.Create(&timeslotDto)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to create timeslot"})
		return
	}
	c.JSON(201, gin.H{"data": data})
}
func (h MeetingHandler) CreateOrUpdateTimeSlot(c *gin.Context) {
	var timeslotDto dto.UpdateTimeSlotDto
	if err := c.BindJSON(&timeslotDto); err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid timeslot data"})
		return
	}

	data, err := h.service.CreateOrUpdate(&timeslotDto)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to update timeslot"})
		return
	}
	c.JSON(201, gin.H{"data": data})
}
