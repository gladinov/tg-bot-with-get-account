package cbr

import (
	"cbr/internal/utils/logging"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gladinov/e"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=HTTPTransport
type HTTPTransport interface {
	DoRequest(ctx context.Context, Path string, query url.Values) ([]byte, error)
}

type Transport struct {
	logger *slog.Logger
	host   string
	client *http.Client
}

func NewTransport(logger *slog.Logger, host string) *Transport {
	return &Transport{
		logger: logger,
		host:   host,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (t *Transport) DoRequest(ctx context.Context, Path string, query url.Values) (_ []byte, err error) {
	const op = "transport.doRequest"
	logg := t.logger.With()
	defer logging.LogOperation_Debug(ctx, logg, op, &err)()

	u := url.URL{
		Scheme: "https",
		Host:   t.host,
		Path:   Path,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		errMsg := "could not create http.NewRequest"
		logging.LoggHTTPError(ctx, logg, req, errMsg, op, err)
		return nil, e.WrapIfErr(errMsg, err)
	}

	req.URL.RawQuery = query.Encode()
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MyApp/1.0)")

	resp, err := t.client.Do(req)
	if err != nil {
		errMsg := "could not do request"
		logging.LoggHTTPError(ctx, logg, req, errMsg, op, err)

		return nil, e.WrapIfErr(errMsg, err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errMsg := "could not read body"
		logging.LoggHTTPError(ctx, logg, req, errMsg, op, err)
		return nil, e.WrapIfErr(errMsg, err)
	}

	return body, nil
}
