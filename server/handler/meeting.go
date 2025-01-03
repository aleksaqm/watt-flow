package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func (h MeetingHandler) GetSlotByDateAndClerkId(c *gin.Context) {
	id := c.DefaultQuery("clerk_id", "0")
	clerkId, err := strconv.ParseUint(id, 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clerk id"})
		return
	}
	date := strings.TrimSpace(c.Query("date"))
	h.logger.Info("Parsing time: ", date)

	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not valid id parameter"})
	}

	timeslotDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid date"})
		return
	}

	data, err := h.service.FindByDateAndClerkId(datatypes.Date(timeslotDate), clerkId)
	if err != nil {
		h.logger.Error(err)
		c.JSON(404, gin.H{"error": " timeslot not found"})
		return
	}
	c.JSON(200, gin.H{"data": data})

}

func (h MeetingHandler) GetMeetingById(c *gin.Context) {
	id := c.Param("id")
	meetingId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid meeting ID"})
		return
	}
	data, err := h.service.FindMeetingById(meetingId)
	if err != nil {
		h.logger.Error(err)
		c.JSON(404, gin.H{"error": "meeting not found"})
		return
	}
	c.JSON(200, gin.H{"data": data})
}

func (h MeetingHandler) GetUsersMeetings(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}
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
	var searchParams dto.MeetingSearchParams
	if search != "" {
		if err := json.Unmarshal([]byte(search), &searchParams); err != nil {
			c.JSON(400, gin.H{"error": "Invalid search parameter"})
			return
		}
	}
	params := dto.MeetingQueryParams{
		Page:      pageInt,
		PageSize:  pageSizeInt,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Search:    searchParams,
	}
	h.logger.Info(params)

	data, total, err := h.service.GetUsersMeetings(userId, &params)
	if err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Failed to load meetings"})
		return
	}
	c.JSON(200, gin.H{"meetings": data, "total": total})
}

func (h MeetingHandler) CreateSlot(c *gin.Context) {
	var timeslotDto dto.TimeSlotDto
	if err := c.BindJSON(&timeslotDto); err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid timeslot data"})
		return
	}

	data, err := h.service.CreateTimeSlot(&timeslotDto)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to create timeslot"})
		return
	}
	c.JSON(201, gin.H{"data": data})
}

func (h MeetingHandler) CreateNewMeeting(c *gin.Context) {
	var meeting dto.NewMeetingDTO
	if err := c.BindJSON(&meeting); err != nil {
		h.logger.Error(err)
		c.JSON(400, gin.H{"error": "Invalid meeting data"})
		return
	}
	timeslotDto := meeting.TimeSlot
	meetingDto := meeting.Meeting

	newMeeting, err := h.service.CreateMeeting(&meetingDto)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to create meeting"})
		return
	}

	timeslotDto.MeetingId = newMeeting.ID
	_, err = h.service.CreateOrUpdate(&timeslotDto)
	if err != nil {
		h.logger.Error(err)
		c.JSON(500, gin.H{"error": "Failed to create or update timeslot"})
		return
	}

	c.JSON(201, gin.H{"data": newMeeting})
}
