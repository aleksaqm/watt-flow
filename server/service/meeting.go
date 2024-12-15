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
	Create(timeslot *dto.TimeSlotDto) (*dto.TimeSlotDto, error)
	CreateOrUpdate(timeslot *dto.UpdateTimeSlotDto) (*dto.TimeSlotDto, error)
}

type MeetingService struct {
	slotRepository *repository.TimeSlotRepository
}

func NewMeetingService(timeslotRepository *repository.TimeSlotRepository) *MeetingService {
	return &MeetingService{
		slotRepository: timeslotRepository,
	}
}

func (t *MeetingService) FindByDate(date datatypes.Date) (*dto.TimeSlotDto, error) {
	timeslot, err := t.slotRepository.FindByDate(date)
	if err != nil {
		return nil, err
	}
	var slots [15]bool
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

func (t *MeetingService) Create(timeslot *dto.TimeSlotDto) (*dto.TimeSlotDto, error) {
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
func (t *MeetingService) CreateOrUpdate(update *dto.UpdateTimeSlotDto) (*dto.TimeSlotDto, error) {
	timeslot, err := t.slotRepository.FindByDate(update.Date)
	if timeslot != nil { //update
		var slots [15]bool
		err = json.Unmarshal(timeslot.Slots, &slots)
		if err != nil {
			return nil, err
		}
		for _, slot := range update.Occupied {
			slots[slot] = false
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

	} else { //create new
		boolArray := make([]bool, 15)

		for i := range boolArray {
			boolArray[i] = true
		}
		for i := range update.Occupied {
			boolArray[i] = false
		}
		slotsJson, err := json.Marshal(boolArray)
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
			Slots:   [15]bool(boolArray),
			ClerkId: newSlot.Clerk.Id,
			Id:      newSlot.Id,
		}, nil

	}
}
