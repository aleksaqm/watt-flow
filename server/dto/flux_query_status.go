package dto

import "time"

type FluxQueryStatusDto struct {
	TimePeriod  string
	GroupPeriod string
	DeviceId    string
	Precision   string
	StartDate   time.Time
	EndDate     time.Time
	Realtime    bool
}
