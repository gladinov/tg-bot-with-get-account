package moex

import (
	"testing"
	"time"
)

const (
	moexHost = "iss.moex.com"
)

func TestGetSpecifications(t *testing.T) {
	client := New(moexHost)
	// Arrange
	ticker := "RU000A1053P7"
	date := time.Now()
	expected := 195.0
	// expected := ""

	// Act
	result, _ := client.GetSpecifications(ticker, date)

	// Assert
	var get float64
	if result.History.Data[0].YieldToOffer != nil {
		get = *result.History.Data[0].YieldToMaturity
		if get != expected {
			t.Errorf("incorect result: expected: %v , get %v", expected, get)
		}
	} else {
		t.Errorf("get nil")
	}

}
