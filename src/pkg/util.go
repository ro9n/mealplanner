package pkg

import (
	"time"
)

// DEFAULT layout for date
// Assumption 3: All dates are of layout yyyy-mm-dd
const DEFAULT = "2006-01-02"

// ParseTime literal to time
func ParseTime(literal string) (time.Time, error) {
	return time.Parse(DEFAULT, literal)
}
