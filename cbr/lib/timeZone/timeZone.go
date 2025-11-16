package timezone

import (
	"fmt"
	"time"
)

const (
	moscowLoc = "Europe/Moscow"
)

func ToMoscowTime(t time.Time) (time.Time, error) {
	const op = "lib.timeZone.ToMoscowTime"
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.Time{}, fmt.Errorf("op:%s failed to load Moscow location", op)

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
		return &time.Location{}, fmt.Errorf("op: %s, error: failed to load Moscow location", op)
	}
	return location, nil
}
