package utils

import (
	"time"

	"github.com/gladinov/e"
)

const (
	moscowLoc = "Europe/Moscow"
	layout    = "02/01/2006"
)

func ToMoscowTime(t time.Time) (time.Time, error) {
	const op = "lib.timeZone.ToMoscowTime"
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.Time{}, e.WrapIfErr("failed to convert time in Moscow location", err)
	}
	return t.In(location), nil
}

func GetStartSingleExchangeRateRubble(location *time.Location) time.Time {
	return time.Date(1992, time.July, 1, 0, 0, 0, 0, location)
}

func GetMoscowLocation() (*time.Location, error) {
	const op = "service.getMoscowLocation"
	location, err := time.LoadLocation(moscowLoc)
	if err != nil {
		return &time.Location{}, e.WrapIfErr("failed to load Moscow location", err)
	}
	return location, nil
}

func MustGetMoscowLocation() *time.Location {
	const op = "service.MustGetMoscowLocation"
	timeLocation, err := GetMoscowLocation()
	if err != nil {
		panic(err)
	}
	return timeLocation
}

func NormalizeDate(date, now, startDate time.Time) string {
	switch {
	case date.After(now):
		return now.Format(layout)
	case date.Before(startDate):
		return startDate.Format(layout)
	default:
		return date.Format(layout)
	}
}
