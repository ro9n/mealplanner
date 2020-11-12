package pkg

import (
	"time"
)

// Activity represents an activity by user
type Activity struct {
	UserID string
	Meals  *map[time.Time]int
}

// RawActivity represents an unprocessed JSON
type RawActivity struct {
	Calendar struct {
		DateDayMap map[string]int `json:"dateToDayId"`
		MealDayMap map[int]int    `json:"mealIdToDayId"`
	} `json:"calendar"`
}
