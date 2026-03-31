//go:build unit

package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseNullFloat64(t *testing.T) {
	cases := []struct {
		name        string
		input       any
		expected    NullFloat64
		expectedErr bool
	}{
		{
			name:  "Error string",
			input: "Hello World",
			expected: NullFloat64{
				IsSet: true,
			},
			expectedErr: true,
		},
		{
			name:  "Error structure",
			input: NullFloat64{},
			expected: NullFloat64{
				IsSet: true,
			},
			expectedErr: true,
		},
		{
			name:  "float64",
			input: 27.31,
			expected: NullFloat64{
				Value:  float64(27.31),
				IsSet:  true,
				IsNull: false,
			},
			expectedErr: false,
		},
		{
			name:  "null",
			input: nil,
			expected: NullFloat64{
				IsSet:  true,
				IsNull: true,
			},
			expectedErr: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseNullFloat64(tc.input)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected.Value, got.Value)
			require.Equal(t, tc.expected.IsSet, got.IsSet)
			require.Equal(t, tc.expected.IsNull, got.IsNull)
		})
	}
}

func TestParseNullString(t *testing.T) {
	cases := []struct {
		name        string
		input       any
		expected    NullString
		expectedErr bool
	}{
		{
			name:  "Error float64",
			input: 27.31,
			expected: NullString{
				IsSet: true,
			},
			expectedErr: true,
		},
		{
			name:  "Error structure",
			input: NullFloat64{},
			expected: NullString{
				IsSet: true,
			},
			expectedErr: true,
		},
		{
			name:  "string",
			input: "Hello World",
			expected: NullString{
				Value:  "Hello World",
				IsSet:  true,
				IsNull: false,
			},
			expectedErr: false,
		},
		{
			name:  "null",
			input: nil,
			expected: NullString{
				IsSet:  true,
				IsNull: true,
			},
			expectedErr: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseNullString(tc.input)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected.Value, got.Value)
			require.Equal(t, tc.expected.IsSet, got.IsSet)
			require.Equal(t, tc.expected.IsNull, got.IsNull)
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	cases := []struct {
		name        string
		input       []byte
		expected    Values
		expectedErr bool
	}{
		{
			name: "Correct",
			input: []byte(`[
			"2025-05-21",
			"2027-09-28",
			null,
			"2026-04-06",
			33.25,
			null,
			1000,
			"RUB",
			283,
			"СибСтекП04"
			]`),
			expected: Values{
				ShortName: NullString{
					Value:  "СибСтекП04",
					IsSet:  true,
					IsNull: false,
				},
				TradeDate: NullString{
					Value:  "2025-05-21",
					IsSet:  true,
					IsNull: false,
				},
				MaturityDate: NullString{
					Value:  "2027-09-28",
					IsSet:  true,
					IsNull: false,
				},
				OfferDate: NullString{
					IsSet:  true,
					IsNull: true,
				},
				BuybackDate: NullString{
					Value:  "2026-04-06",
					IsSet:  true,
					IsNull: false,
				},
				YieldToMaturity: NullFloat64{
					Value:  33.25,
					IsSet:  true,
					IsNull: false,
				},
				YieldToOffer: NullFloat64{
					IsSet:  true,
					IsNull: true,
				},
				FaceValue: NullFloat64{
					Value:  1000,
					IsSet:  true,
					IsNull: false,
				},
				FaceUnit: NullString{
					Value:  "RUB",
					IsSet:  true,
					IsNull: false,
				},
				Duration: NullFloat64{
					Value:  283,
					IsSet:  true,
					IsNull: false,
				},
			},
			expectedErr: false,
		},
		{
			name: "Error: less 10 elements",
			input: []byte(`[
			"2025-05-21",
			"2027-09-28",
			null,
			"2026-04-06",
			33.25,
			null,
			1000,
			"RUB",
			283
			]`),
			expected:    Values{},
			expectedErr: true,
		},
		{
			name: "Error: only strings",
			input: []byte(`[
			"2025-05-21",
			"2025-05-21",
			"2025-05-21",
			"2026-04-06",
			"2025-05-21",
			"2025-05-21",
			"2025-05-21",
			"2025-05-21",
			"2025-05-21",
			"2025-05-21"
			]`),
			expected:    Values{},
			expectedErr: true,
		},
		{
			name: "Error: only floats",
			input: []byte(`[
			27,
			27,
			27,
			27,
			27,
			27,
			27,
			27,
			27,
			27"
			]`),
			expected:    Values{},
			expectedErr: true,
		},
		{
			name: "only nulls",
			input: []byte(`[
			null,
			null,
			null,
			null,
			null,
			null,
			null,
			null,
			null,
			null
			]`),
			expected: Values{
				ShortName: NullString{
					IsSet:  true,
					IsNull: true,
				},
				TradeDate: NullString{
					IsSet:  true,
					IsNull: true,
				},
				MaturityDate: NullString{
					IsSet:  true,
					IsNull: true,
				},
				OfferDate: NullString{
					IsSet:  true,
					IsNull: true,
				},
				BuybackDate: NullString{
					IsSet:  true,
					IsNull: true,
				},
				YieldToMaturity: NullFloat64{
					IsSet:  true,
					IsNull: true,
				},
				YieldToOffer: NullFloat64{
					IsSet:  true,
					IsNull: true,
				},
				FaceValue: NullFloat64{
					IsSet:  true,
					IsNull: true,
				},
				FaceUnit: NullString{
					IsSet:  true,
					IsNull: true,
				},
				Duration: NullFloat64{
					IsSet:  true,
					IsNull: true,
				},
			},
			expectedErr: false,
		},
		{
			name:        "crash unmarshall",
			input:       []byte(`[ "2025-05-21", `),
			expected:    Values{},
			expectedErr: true,
		},
		{
			name:        "not an array",
			input:       []byte(`{"key": "value"}`),
			expected:    Values{},
			expectedErr: true,
		},
		{
			name: "error in element 0 (TRADEDATE)",
			input: []byte(`[
                ["invalid"], "2027-09-28", null, "2026-04-06", 33.25, null, 1000, "RUB", 283, "СибСтекП04"
            ]`),
			expected:    Values{},
			expectedErr: true,
		},
		{
			name: "error in element 1 (MATDATE)",
			input: []byte(`[
                "2025-05-21", ["invalid"], null, "2026-04-06", 33.25, null, 1000, "RUB", 283, "СибСтекП04"
            ]`),
			expected:    Values{},
			expectedErr: true,
		},
		{
			name: "error in element 2 (OFFERDATE)",
			input: []byte(`[
                "2025-05-21", "2027-09-28", ["invalid"], "2026-04-06", 33.25, null, 1000, "RUB", 283, "СибСтекП04"
            ]`),
			expected:    Values{},
			expectedErr: true,
		},
		{
			name: "error in element 3 (BUYBACKDATE)",
			input: []byte(`[
                "2025-05-21", "2027-09-28", null, ["invalid"], 33.25, null, 1000, "RUB", 283, "СибСтекП04"
            ]`),
			expected:    Values{},
			expectedErr: true,
		},
		{
			name: "error in element 4 (YIELDCLOSE)",
			input: []byte(`[
                "2025-05-21", "2027-09-28", null, "2026-04-06", "not_a_number", null, 1000, "RUB", 283, "СибСтекП04"
            ]`),
			expected:    Values{},
			expectedErr: true,
		},
		{
			name: "error in element 5 (YIELDTOOFFER)",
			input: []byte(`[
                "2025-05-21", "2027-09-28", null, "2026-04-06", 33.25, "not_a_number", 1000, "RUB", 283, "СибСтекП04"
            ]`),
			expected:    Values{},
			expectedErr: true,
		},
		{
			name: "error in element 6 (FACEVALUE)",
			input: []byte(`[
                "2025-05-21", "2027-09-28", null, "2026-04-06", 33.25, null, "not_a_number", "RUB", 283, "СибСтекП04"
            ]`),
			expected:    Values{},
			expectedErr: true,
		},
		{
			name: "error in element 7 (FACEUNIT)",
			input: []byte(`[
                "2025-05-21", "2027-09-28", null, "2026-04-06", 33.25, null, 1000, ["invalid"], 283, "СибСтекП04"
            ]`),
			expected:    Values{},
			expectedErr: true,
		},
		{
			name: "error in element 8 (DURATION)",
			input: []byte(`[
                "2025-05-21", "2027-09-28", null, "2026-04-06", 33.25, null, 1000, "RUB", "not_a_number", "СибСтекП04"
            ]`),
			expected:    Values{},
			expectedErr: true,
		},
		{
			name: "error in element 9 (SHORTNAME)",
			input: []byte(`[
                "2025-05-21", "2027-09-28", null, "2026-04-06", 33.25, null, 1000, "RUB", 283, ["invalid"]
            ]`),
			expected:    Values{},
			expectedErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var values Values
			err := json.Unmarshal(tc.input, &values)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assertNullStringEqual(t, "ShortName", tc.expected.ShortName, values.ShortName)
			assertNullStringEqual(t, "TradeDate", tc.expected.TradeDate, values.TradeDate)
			assertNullStringEqual(t, "MaturityDate", tc.expected.MaturityDate, values.MaturityDate)
			assertNullStringEqual(t, "OfferDate", tc.expected.OfferDate, values.OfferDate)
			assertNullStringEqual(t, "BuybackDate", tc.expected.BuybackDate, values.BuybackDate)
			assertNullFloat64Equal(t, "YieldToMaturity", tc.expected.YieldToMaturity, values.YieldToMaturity)
			assertNullFloat64Equal(t, "YieldToOffer", tc.expected.YieldToOffer, values.YieldToOffer)
			assertNullFloat64Equal(t, "FaceValue", tc.expected.FaceValue, values.FaceValue)
			assertNullStringEqual(t, "FaceUnit", tc.expected.FaceUnit, values.FaceUnit)
			assertNullFloat64Equal(t, "Duration", tc.expected.Duration, values.Duration)
		})
	}
}

func assertNullStringEqual(t *testing.T, fieldName string, expected, actual NullString) {
	t.Helper()
	require.Equal(t, expected.Value, actual.Value, "%s.Value", fieldName)
	require.Equal(t, expected.IsNull, actual.IsNull, "%s.IsNull", fieldName)
	require.Equal(t, expected.IsSet, actual.IsSet, "%s.IsSet", fieldName)
}

func assertNullFloat64Equal(t *testing.T, fieldName string, expected, actual NullFloat64) {
	t.Helper()
	require.Equal(t, expected.Value, actual.Value, "%s.Value", fieldName)
	require.Equal(t, expected.IsNull, actual.IsNull, "%s.IsNull", fieldName)
	require.Equal(t, expected.IsSet, actual.IsSet, "%s.IsSet", fieldName)
}
