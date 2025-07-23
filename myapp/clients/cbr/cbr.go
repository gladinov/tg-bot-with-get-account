package cbr

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"golang.org/x/text/encoding/charmap"
	"main.go/lib/e"
)

const (
	layout = "02/01/2006"
)

type Valute struct {
	NumCode   string `xml:"NumCode"`
	CharCode  string `xml:"CharCode"`
	Nominal   string `xml:"Nominal"`
	Name      string `xml:"Name"`
	Value     string `xml:"Value"`
	VunitRate string `xml:"VunitRate"`
}

type ValCurs struct {
	Date   string   `xml:"Date,attr"`
	Valute []Valute `xml:"Valute"`
}

type Client struct {
	host   string
	client http.Client
}

func New(host string) *Client {
	return &Client{
		host:   host,
		client: http.Client{},
	}
}

func (c *Client) GetAllCurrencies(date time.Time) (curr ValCurs, err error) {
	defer func() { err = e.WrapIfErr("getAllCurrencies error", err) }()
	formatDate := date.Format(layout)
	Path := path.Join("scripts", "XML_daily.asp")
	params := url.Values{}
	params.Add("date_req", formatDate)
	body, err := c.doRequest(Path, params)
	if err != nil {
		return curr, err
	}

	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.CharsetReader = func(label string, input io.Reader) (io.Reader, error) {
		if label == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return input, nil
	}

	err = decoder.Decode(&curr)
	if err != nil {
		return curr, err
	}

	return curr, nil
}

func (c *Client) doRequest(Path string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can`t do request", err) }()
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   Path,
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MyApp/1.0)")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
