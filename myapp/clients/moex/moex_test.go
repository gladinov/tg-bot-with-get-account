package moex

import (
	"testing"
	"time"
)

func (c *Client) TestGetSpecifications(t *testing.T) {
	// Arrange
	ticker := "RU000A104Y15"
	date := time.Date(2023, time.October, 17, 0, 0, 0, 0, time.Local)
	expected := 13.63

	// Act
	result, _ := c.GetSpecifications(ticker, date)

	// Assert
	if *result.History.Data[0].YieldToMaturity != expected {
		t.Errorf("incorect result")
	}
}
