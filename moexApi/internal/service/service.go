package service

import (
	"encoding/json"
	"io"
	"main/lib/e"
	"net/http"
	"net/url"
	"path"
)

const (
	layout = "2006-01-02"
)

type Service interface {
	GetSpecifications(req SpecificationsRequest) (values Values, err error)
}

type SpecificationService struct {
	host   string
	client http.Client
}

func NewSpecificationService(host string) Service {
	return &SpecificationService{
		host:   host,
		client: http.Client{},
	}
}

func (s *SpecificationService) GetSpecifications(req SpecificationsRequest) (values Values, err error) {
	defer func() { err = e.WrapIfErr("getspecification error", err) }()
	ticker := req.Ticker
	date := req.Date
	var data SpecificationsResponce
	daysMax := 14
	for dayToSubstract := 1; dayToSubstract <= daysMax; dayToSubstract++ {

		// Проверяем условия выхода из цикла
		formatDate := date.Format(layout)
		// uri := fmt.Sprintf("https://iss.moex.com/iss/history/engines/stock/markets/bonds/sessions/3/securities/%s.json", ticker)
		Path := path.Join("iss", "history", "engines", "stock", "markets", "bonds", "sessions", "3", "securities", ticker+".json")
		params := url.Values{}
		params.Add("limit", "1")
		params.Add("iss.meta", "off")
		params.Add("history.columns", "TRADEDATE,MATDATE,OFFERDATE,BUYBACKDATE,YIELDCLOSE,YIELDTOOFFER,FACEVALUE,FACEUNIT,DURATION, SHORTNAME")
		params.Add("limit", "1")
		params.Add("from", formatDate)
		params.Add("to", formatDate)

		body, err := s.doRequest(Path, params)
		if err != nil {
			return Values{}, err
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
			return Values{}, err
		}
		if data.History != nil {

			if len(data.History.Data) != 0 {
				break
			}
		}
		date = date.AddDate(0, 0, -1)
	}
	resp := data.History.Data[0]
	return resp, nil
}

func (s *SpecificationService) doRequest(Path string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can`t do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   s.host,
		Path:   Path,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := s.client.Do(req)
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
