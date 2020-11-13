package pkg

import (
	"log"
	"time"
)

// TimedActivityPredicate represents a time predicate
type TimedActivityPredicate struct {
	Starts time.Time
	Ends   time.Time
}

// newTimedActivityPredicate creates a new instance for query
func newTimedActivityPredicate(starts, ends time.Time) *TimedActivityPredicate {
	return &TimedActivityPredicate{Starts: starts, Ends: ends}
}

// Evaluate number of activities in the timeframe
func (p *TimedActivityPredicate) Evaluate(activity *Activity) int {
	count := 0
	for date, mealCount := range *activity.Meals {
		// Assumption 5. lookup date range is inclusive i.e. [start, end], not (start, end)
		if date.Before(p.Starts) || date.After(p.Ends) {
			log.Println("Rejecting", date)
			continue
		}

		log.Println("Accepting", date)
		count += mealCount
	}

	return count
}
