package dto

import "time"

type FluxQueryConsumptionDto struct {
	TimePeriod  string
	GroupPeriod string
	City        string
	Precision   string
	StartDate   time.Time
	EndDate     time.Time
	Realtime    bool
}
