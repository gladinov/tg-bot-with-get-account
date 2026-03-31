package cbr

import (
	"bytes"
	"cbr/internal/models"
	"cbr/internal/utils/logging"
	"context"
	"encoding/xml"
	"io"
	"log/slog"

	"github.com/gladinov/e"
	"golang.org/x/text/encoding/charmap"
)

func parseCurrencies(ctx context.Context, logger *slog.Logger, data []byte) (_ models.CurrenciesResponce, err error) {
	const op = "cbr.parseCurrencies"
	logg := logger.With()
	defer logging.LogOperation_Debug(ctx, logg, op, &err)()

	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = func(label string, input io.Reader) (io.Reader, error) {
		if label == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return input, nil
	}
	var curr models.CurrenciesResponce
	err = decoder.Decode(&curr)
	if err != nil {
		return models.CurrenciesResponce{}, e.WrapIfErr("could not decode Xml file", err)
	}

	return curr, nil
}
