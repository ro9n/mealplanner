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

type predicate interface {
	Evaluate(activity *Activity) int
}

// Query represents query execution Query
type Query interface {
	Apply(directory string) []string
}

type activeQuery struct {
	period predicate
}

type superActiveQuery struct {
	period predicate
}

type boredQuery struct {
	period    predicate
	preceding predicate
}

// NewQuery creates query
func NewQuery(strategy, starts, ends string) (Query, error) {
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
		return newActiveQuery(s, e), nil
	case SUPERACTIVE:
		return newSuperActiveQuery(s, e), nil
	case BORED:
		return newBoredQuery(s, e), nil
	default:
		return nil, nil
	}
}

func newActiveQuery(starts, ends time.Time) *activeQuery {
	return &activeQuery{newTimedActivityPredicate(starts, ends)}
}

func newSuperActiveQuery(starts, ends time.Time) *superActiveQuery {
	return &superActiveQuery{newTimedActivityPredicate(starts, ends)}
}

func newBoredQuery(starts, ends time.Time) *boredQuery {
	// Assumption 4. preceding time refers to same size window ending before start date
	// e.g. if start=16/1 end=30/1, preceding start=1/1 end = 15/1
	d := ends.Sub(starts) + 24*time.Hour
	return &boredQuery{newTimedActivityPredicate(starts, ends), newTimedActivityPredicate(starts.Add(-d), ends.Add(-d))}
}

func (c activeQuery) Apply(directory string) []string {
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

func (c superActiveQuery) Apply(directory string) []string {
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

func (c boredQuery) Apply(directory string) []string {
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
