package handlers

import "time"

type CurrencyRequest struct {
	Date time.Time `json:"date,omitempty"`
}
