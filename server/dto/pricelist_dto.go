package dto

type NewPricelist struct {
	Year         int     `json:"year"`
	Month        int     `json:"month"`
	BlueZone     float64 `json:"blue"`
	RedZone      float64 `json:"red"`
	GreenZone    float64 `json:"green"`
	BillingPower float64 `json:"bill_power"`
	Tax          float64 `json:"tax"`
}
