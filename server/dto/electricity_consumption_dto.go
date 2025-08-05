package dto

type ElectricityConsumptionDto struct {
	HouseholdId string `json:"household_id"`
	Year        int    `json:"year"`
	Month       int    `json:"month"`
}

type MonthlyConsumptionResult struct {
	Year        int     `json:"year"`
	Month       int     `json:"month"`
	MonthName   string  `json:"month_name"`
	Consumption float64 `json:"consumption"`
}

type ElectricityConsumptionResponse struct {
	Data []MonthlyConsumptionResult `json:"data"`
}

type Get12MonthsConsumptionDto struct {
	HouseholdId string `json:"household_id"`
	EndYear     int    `json:"end_year"`
	EndMonth    int    `json:"end_month"`
}

type DailyConsumptionData struct {
	Year        int     `json:"year"`
	Month       int     `json:"month"`
	Day         int     `json:"day"`
	Consumption float64 `json:"consumption"`
}

type DailyConsumptionResponse struct {
	Data []DailyConsumptionData `json:"data"`
}
