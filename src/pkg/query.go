package pkg

import (
	"log"
	"time"
)

// TimedActivityQuery represents a time query
type TimedActivityQuery struct {
	Starts time.Time
	Ends   time.Time
}

// NewTimedActivityQuery creates a new instance for query
func newTimedActivityQuery(starts, ends time.Time) *TimedActivityQuery {
	return &TimedActivityQuery{Starts: starts, Ends: ends}
}

// Evaluate number of activities in the timeframe
func (query *TimedActivityQuery) Evaluate(activity *Activity) int {
	count := 0
	for date, mealCount := range *activity.Meals {
		// Assumption 5: lookup date range is inclusive i.e. [start, end], not (start, end)
		if date.Before(query.Starts) || date.After(query.Ends) {
			log.Println("Rejecting", date)
			continue
		}

		log.Println("Accepting", date)
		count += mealCount
	}

	return count
}
