package dto

import "time"

type FluxQueryConsumptionDto struct {
	TimePeriod  string
	GroupPeriod string
	DeviceId    string
	Precision   string
	StartDate   time.Time
	EndDate     time.Time
	Realtime    bool
}

type ConsumptionQueryResult struct {
	Rows []ConsumptionQueryResultRow
}

type ConsumptionQueryResultRow struct {
	TimeField time.Time
	Value     float64
}
