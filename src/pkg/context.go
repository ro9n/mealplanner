package pkg

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// Predicates for ingestion
const (
	ACTIVE      = "active"
	SUPERACTIVE = "superactive"
	BORED       = "bored"
)

type query interface {
	Evaluate(activity *Activity) int
}

// Context represents query execution context
type Context interface {
	Apply(directory string) []string
}

type activeContext struct {
	period query
}

type superActiveContext struct {
	period query
}

type boredContext struct {
	period    query
	preceding query
}

// NewContext creates query context
func NewContext(strategy, starts, ends string) (Context, error) {
	if !(strategy == ACTIVE || strategy == SUPERACTIVE || strategy == BORED) {
		return nil, errors.New(fmt.Sprint("Invalid argument", strategy))
	}

	s, err := ParseTime(starts)
	if err != nil {
		return nil, err
	}

	e, err := ParseTime(ends)
	if err != nil {
		return nil, err
	}

	switch strategy {
	case ACTIVE:
		return newActiveContext(s, e), nil
	case SUPERACTIVE:
		return newSuperActiveContext(s, e), nil
	case BORED:
		return newBoredContext(s, e), nil
	default:
		return nil, nil
	}
}

func newActiveContext(starts, ends time.Time) *activeContext {
	return &activeContext{newTimedActivityQuery(starts, ends)}
}

func newSuperActiveContext(starts, ends time.Time) *superActiveContext {
	return &superActiveContext{newTimedActivityQuery(starts, ends)}
}

func newBoredContext(starts, ends time.Time) *boredContext {
	// Assumption 4: preceding time refers to same size window ending before start date
	// e.g. if start=16/1 end=30/1, preceding start=1/1 end = 15/1
	d := ends.Sub(starts) + 24*time.Hour
	return &boredContext{newTimedActivityQuery(starts, ends), newTimedActivityQuery(starts.Add(-d), ends.Add(-d))}
}

func (c activeContext) Apply(directory string) []string {
	var users []string

	w := NewWorker(NewJSONReader())

	for activity := range w.process(directory) {
		count := c.period.Evaluate(&activity)
		log.Println(activity.UserID, count)

		if count >= 5 {
			users = append(users, activity.UserID)
		}
	}

	return users
}

func (c superActiveContext) Apply(directory string) []string {
	var users []string

	w := NewWorker(NewJSONReader())

	for activity := range w.process(directory) {
		count := c.period.Evaluate(&activity)
		log.Println(activity.UserID, count)

		if count > 10 {
			users = append(users, activity.UserID)
		}
	}

	return users
}

func (c boredContext) Apply(directory string) []string {
	var users []string

	w := NewWorker(NewJSONReader())

	for activity := range w.process(directory) {
		count := c.period.Evaluate(&activity)
		log.Println(activity.UserID, count)

		if count < 5 {
			precedingCount := c.preceding.Evaluate(&activity)
			log.Println(activity.UserID, precedingCount)

			if precedingCount >= 5 {
				users = append(users, activity.UserID)
			}
		}
	}

	return users
}
