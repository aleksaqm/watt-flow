package service

import (
	"encoding/json"
	"watt-flow/dto"
	"watt-flow/model"
	"watt-flow/repository"

	"gorm.io/datatypes"
)

type IMeetingService interface {
	FindByDate(id datatypes.Date) (*dto.TimeSlotDto, error)
	FindByDateAndClerkId(date datatypes.Date, clerkId uint64) (*dto.TimeSlotDto, error)
	CreateTimeSlot(timeslot *dto.TimeSlotDto) (*dto.TimeSlotDto, error)
	CreateOrUpdate(timeslot *dto.UpdateTimeSlotDto) (*dto.TimeSlotDto, error)
	CreateMeeting(meetingDto *dto.MeetingDTO) (*dto.MeetingDTO, error)
	FindMeetingById(id uint64) (*dto.MeetingDTO, error)
	FindMeetingBySlotId(slotId uint64) (*dto.MeetingDTO, error)
	GetUsersMeetings(userID uint64, params *dto.MeetingQueryParams) ([]dto.UsersMeetingDTO, int64, error)
}

type MeetingService struct {
	slotRepository    repository.TimeSlotRepository
	meetingRepository repository.MeetingRepository
}

func NewMeetingService(timeslotRepository repository.TimeSlotRepository, meetingRepository repository.MeetingRepository) *MeetingService {
	return &MeetingService{
		slotRepository:    timeslotRepository,
		meetingRepository: meetingRepository,
	}
}

func (t *MeetingService) FindByDate(date datatypes.Date) (*dto.TimeSlotDto, error) {
	timeslot, err := t.slotRepository.FindByDate(date)
	if err != nil {
		return nil, err
	}
	var slots [15]uint64
	err = json.Unmarshal(timeslot.Slots, &slots)
	if err != nil {
		return nil, err
	}
	return &dto.TimeSlotDto{
		Date:    timeslot.Date,
		Slots:   slots,
		ClerkId: timeslot.Clerk.Id,
		Id:      timeslot.Id,
	}, nil
}

func (t *MeetingService) FindByDateAndClerkId(date datatypes.Date, clerkId uint64) (*dto.TimeSlotDto, error) {
	timeslot, err := t.slotRepository.FindByDateAndClerkID(date, clerkId)
	if err != nil {
		return nil, err
	}
	var slots [15]uint64
	err = json.Unmarshal(timeslot.Slots, &slots)
	if err != nil {
		return nil, err
	}
	return &dto.TimeSlotDto{
		Date:    timeslot.Date,
		Slots:   slots,
		ClerkId: timeslot.Clerk.Id,
		Id:      timeslot.Id,
	}, nil
}

func (t *MeetingService) FindMeetingById(id uint64) (*dto.MeetingDTO, error) {
	meeting, err := t.meetingRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	return &dto.MeetingDTO{
		ID:        meeting.ID,
		StartTime: meeting.StartTime,
		Duration:  meeting.Duration,
		ClerkID:   meeting.ClerkID,
		UserID:    meeting.UserID,
		Username:  meeting.User.Username,
	}, nil
}

func (t *MeetingService) FindMeetingBySlotId(slotId uint64) (*dto.MeetingDTO, error) {
	meeting, err := t.meetingRepository.FindBySlotId(slotId)
	if err != nil {
		return nil, err
	}
	return &dto.MeetingDTO{
		ID:        meeting.ID,
		StartTime: meeting.StartTime,
		Duration:  meeting.Duration,
		ClerkID:   meeting.ClerkID,
		UserID:    meeting.UserID,
		Username:  meeting.User.Username,
	}, nil
}

func (t *MeetingService) CreateTimeSlot(timeslot *dto.TimeSlotDto) (*dto.TimeSlotDto, error) {
	slotsJson, err := json.Marshal(timeslot.Slots)
	if err != nil {
		return nil, err
	}
	timeSlot := model.TimeSlot{
		ClerkID: timeslot.ClerkId,
		Slots:   slotsJson,
		Date:    timeslot.Date,
	}
	_, err = t.slotRepository.Create(&timeSlot)
	if err != nil {
		return nil, err
	}
	return timeslot, nil
}

func (t *MeetingService) CreateMeeting(meetingDto *dto.MeetingDTO) (*dto.MeetingDTO, error) {
	meeting := model.Meeting{
		StartTime: meetingDto.StartTime,
		Duration:  meetingDto.Duration,
		ClerkID:   meetingDto.ClerkID,
		UserID:    meetingDto.UserID,
	}
	_, err := t.meetingRepository.Create(&meeting)
	if err != nil {
		return nil, err
	}
	return &dto.MeetingDTO{
		ID:        meeting.ID,
		StartTime: meeting.StartTime,
		Duration:  meeting.Duration,
		ClerkID:   meeting.ClerkID,
		UserID:    meeting.UserID,
		Username:  meeting.User.Username,
	}, nil
}

func (t *MeetingService) CreateOrUpdate(update *dto.UpdateTimeSlotDto) (*dto.TimeSlotDto, error) {
	timeslot, err := t.slotRepository.FindByDateAndClerkId(update.Date, update.ClerkId) // changed
	if timeslot != nil {                                                                // update
		var slots [15]uint64
		err = json.Unmarshal(timeslot.Slots, &slots)
		if err != nil {
			return nil, err
		}
		for _, slot := range update.Occupied {
			slots[slot] = update.MeetingId
		}
		slotsJson, err := json.Marshal(slots)
		if err != nil {
			return nil, err
		}
		timeslot.Slots = slotsJson
		newSlot, err := t.slotRepository.Update(timeslot)
		if err != nil {
			return nil, err
		}
		return &dto.TimeSlotDto{
			Date:    newSlot.Date,
			Slots:   slots,
			ClerkId: newSlot.Clerk.Id,
			Id:      newSlot.Id,
		}, nil

	} else { // create new
		emptySlot := make([]uint64, 15)

		for i := range emptySlot {
			emptySlot[i] = 0
		}
		for _, element := range update.Occupied {
			emptySlot[element] = update.MeetingId
		}
		slotsJson, err := json.Marshal(emptySlot)
		if err != nil {
			return nil, err
		}
		newTimeSlot := model.TimeSlot{
			ClerkID: update.ClerkId,
			Slots:   slotsJson,
			Date:    update.Date,
		}
		newSlot, err := t.slotRepository.Create(&newTimeSlot)
		if err != nil {
			return nil, err
		}
		return &dto.TimeSlotDto{
			Date:    newSlot.Date,
			Slots:   [15]uint64(emptySlot),
			ClerkId: newSlot.Clerk.Id,
			Id:      newSlot.Id,
		}, nil

	}
}

func (t *MeetingService) GetUsersMeetings(userID uint64, params *dto.MeetingQueryParams) ([]dto.UsersMeetingDTO, int64, error) {
	meetings, total, err := t.meetingRepository.FindForUser(userID, params)
	if err != nil {
		return nil, 0, err
	}
	results := make([]dto.UsersMeetingDTO, 0)
	for _, result := range meetings {
		mappedRequest, _ := t.MapToUsersMeetingDto(&result)
		results = append(results, mappedRequest)
	}
	return results, total, nil
}

func (t *MeetingService) MapToUsersMeetingDto(meeting *model.Meeting) (dto.UsersMeetingDTO, error) {
	response := dto.UsersMeetingDTO{
		ID:        meeting.ID,
		StartTime: meeting.StartTime,
		Duration:  meeting.Duration,
		Clerk:     meeting.Clerk.Username,
	}
	return response, nil
}
