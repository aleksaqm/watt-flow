package main

import (
	"math"
	"math/rand"
	"time"
)

type Household struct {
	rng            *rand.Rand
	address        *Location
	baseLoad       float64
	peakMultiplier float64
	seasonImpact   float64
}

type Location struct {
	City   string
	Street string
	Number string
}

func newLocation(city string, street string, number string) *Location {
	return &Location{
		City:   city,
		Street: street,
		Number: number,
	}
}

func newHousehold(seed int64, location *Location) *Household {
	rng := rand.New(rand.NewSource(seed))
	return &Household{
		rng:            rng,
		baseLoad:       0.7 + (rng.Float64() * 0.6),
		peakMultiplier: 2.5 + rng.Float64(),
		seasonImpact:   0.3 + (rng.Float64() * 0.6),
		address:        location,
	}
}

func (hs *Household) calculateDailyPattern(hour float64) float64 {
	// Morning peak (7-9 AM)
	if hour >= 7 && hour <= 9 {
		return hs.peakMultiplier * 0.8
	}
	// Evening peak (18-22)
	if hour >= 18 && hour <= 22 {
		return hs.peakMultiplier
	}
	// Nighttime (23-5)
	if hour >= 23 || hour <= 5 {
		return 0.8
	}
	// Mid-day (9-18)
	if hour > 9 && hour < 18 {
		return 1.5
	}
	// Early morning (5-7)
	return 1.2
}

func (hs *Household) calculateSeasonalFactor(month time.Month) float64 {
	monthAngle := float64(month-1) * (2 * math.Pi / 12)
	// Create a sinusoidal pattern with peak in winter (December/January)
	// and trough in summer (June/July)
	seasonalVariation := math.Cos(monthAngle)*hs.seasonImpact + 1.0
	return seasonalVariation
}

func (hs *Household) SimulateConsumption(timestamp time.Time) float64 {
	// Time-of-day factor
	hour := float64(timestamp.Hour())
	dailyPattern := hs.calculateDailyPattern(hour)

	// Seasonal factor
	month := timestamp.Month()
	seasonalFactor := hs.calculateSeasonalFactor(month)

	// Random variation (±15% to simulate appliance usage)
	randomFactor := 0.85 + (hs.rng.Float64() * 0.3)

	// Calculate total consumption
	consumption := hs.baseLoad * dailyPattern * seasonalFactor * randomFactor

	return consumption
}
