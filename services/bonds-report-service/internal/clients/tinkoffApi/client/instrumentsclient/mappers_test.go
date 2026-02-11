//go:build unit

package instrumentsclient

import (
	dtoTinkoff "bonds-report-service/internal/models/dto/tinkoffApi"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapSliceInstrumentShortToDomain(t *testing.T) {
	in := []dtoTinkoff.InstrumentShort{
		{Uid: "1"},
		{Uid: "2"},
	}

	out := MapSliceInstrumentShortToDomain(in)

	require.Len(t, out, 2)
	require.Equal(t, "1", out[0].Uid)
	require.Equal(t, "2", out[1].Uid)
}
