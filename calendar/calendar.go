package calendar

import (
	"time"
)

type Calendar interface {
	// 遅延して実行させたい時とかどうしようか
	PostSchedule(schedule *Schedule) error
	GetSchedule(at time.Time) (*Schedule, error)
	UpdateSchedule() error
}

type Schedule struct {
	at          time.Time
	title       string
	description string
	remainders  []*Remainder
}

type Remainder struct {
	at time.Time
}
