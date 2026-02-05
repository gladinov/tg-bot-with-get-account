//go:build unit

package moex

import (
	"testing"

	domainModel "bonds-report-service/internal/models/domain"
	moexDTO "bonds-report-service/internal/models/dto/moex"

	"github.com/stretchr/testify/assert"
)

func TestMapNullStringFromDTOToDomain(t *testing.T) {
	tests := []struct {
		name string
		dto  moexDTO.NullString
		want domainModel.NullString
	}{
		{
			name: "value is set",
			dto: moexDTO.NullString{
				Value:  "ABC",
				IsSet:  true,
				IsNull: false,
			},
			want: domainModel.NewNullString("ABC", true, false),
		},
		{
			name: "explicit null",
			dto: moexDTO.NullString{
				Value:  "",
				IsSet:  true,
				IsNull: true,
			},
			want: domainModel.NewNullString("", true, true),
		},
		{
			name: "not set",
			dto: moexDTO.NullString{
				Value:  "",
				IsSet:  false,
				IsNull: false,
			},
			want: domainModel.NewNullString("", false, false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapNullStringFromDTOToDomain(tt.dto)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapNullFloat64FromDTOToDomain(t *testing.T) {
	tests := []struct {
		name string
		dto  moexDTO.NullFloat64
		want domainModel.NullFloat64
	}{
		{
			name: "value is set",
			dto: moexDTO.NullFloat64{
				Value:  12.34,
				IsSet:  true,
				IsNull: false,
			},
			want: domainModel.NewNullFloat64(12.34, true, false),
		},
		{
			name: "explicit null",
			dto: moexDTO.NullFloat64{
				Value:  0,
				IsSet:  true,
				IsNull: true,
			},
			want: domainModel.NewNullFloat64(0, true, true),
		},
		{
			name: "not set",
			dto: moexDTO.NullFloat64{
				Value:  0,
				IsSet:  false,
				IsNull: false,
			},
			want: domainModel.NewNullFloat64(0, false, false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapNullFloat64FromDTOToDomain(tt.dto)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapValueFromDTOToDomain(t *testing.T) {
	dto := moexDTO.Values{
		ShortName: moexDTO.NullString{
			Value:  "OFZ 26238",
			IsSet:  true,
			IsNull: false,
		},
		TradeDate: moexDTO.NullString{
			Value:  "2024-01-01",
			IsSet:  true,
			IsNull: false,
		},
		MaturityDate: moexDTO.NullString{
			Value:  "2034-01-01",
			IsSet:  true,
			IsNull: false,
		},
		OfferDate: moexDTO.NullString{
			Value:  "",
			IsSet:  true,
			IsNull: true,
		},
		BuybackDate: moexDTO.NullString{
			IsSet:  false,
			IsNull: false,
		},
		YieldToMaturity: moexDTO.NullFloat64{
			Value:  9.87,
			IsSet:  true,
			IsNull: false,
		},
		YieldToOffer: moexDTO.NullFloat64{
			IsSet:  false,
			IsNull: false,
		},
		FaceValue: moexDTO.NullFloat64{
			Value:  1000,
			IsSet:  true,
			IsNull: false,
		},
		FaceUnit: moexDTO.NullString{
			Value:  "RUB",
			IsSet:  true,
			IsNull: false,
		},
		Duration: moexDTO.NullFloat64{
			Value:  4.56,
			IsSet:  true,
			IsNull: false,
		},
	}

	got := MapValueFromDTOToDomain(dto)

	want := domainModel.ValuesMoex{
		ShortName:       domainModel.NewNullString("OFZ 26238", true, false),
		TradeDate:       domainModel.NewNullString("2024-01-01", true, false),
		MaturityDate:    domainModel.NewNullString("2034-01-01", true, false),
		OfferDate:       domainModel.NewNullString("", true, true),
		BuybackDate:     domainModel.NewNullString("", false, false),
		YieldToMaturity: domainModel.NewNullFloat64(9.87, true, false),
		YieldToOffer:    domainModel.NewNullFloat64(0, false, false),
		FaceValue:       domainModel.NewNullFloat64(1000, true, false),
		FaceUnit:        domainModel.NewNullString("RUB", true, false),
		Duration:        domainModel.NewNullFloat64(4.56, true, false),
	}

	assert.Equal(t, want, got)
}
