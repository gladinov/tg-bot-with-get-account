package moex

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"main.go/lib/e"
)

const (
	layout = "2006-01-02"
)

type Client struct {
	host   string
	client http.Client
}

func New(host string) *Client {
	return &Client{
		host:   host,
		client: http.Client{}}
}

type Yields struct {
	History *History `json:"history"`
}

type History struct {
	Data []Values `json:"data"`
}

type Values struct {
	TradeDate       *string  `json:"TRADEDATE"`    // Торговая дата(на момент которой рассчитаны остальные данные)
	MaturityDate    *string  `json:"MATDATE"`      // Дата погашения
	OfferDate       *string  `json:"OFFERDATE"`    // Дата Оферты
	BuybackDate     *string  `json:"BUYBACKDATE"`  // дата обратного выкупа
	YieldToMaturity *float64 `json:"YIELDCLOSE"`   // Доходность к погашению при покупке
	YieldToOffer    *float64 `json:"YIELDTOOFFER"` // Доходность к оферте при покупке
	FaceValue       *float64 `json:"FACEVALUE"`
	FaceUnit        *float64 `json:"FACEUNIT"` // номинальная стоимость облигации
	Duration        *float64 `json:"DURATION"` // дюрация (средневзвешенный срок платежей)

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

func (c *Client) GetSpecifications(ticker string, date time.Time) (yields *Yields, err error) {
	// Ищем котировки на протяжении последних 14 дней.(для исключения попадания на выходные и нерабочие для биржи дни)
	defer func() { err = e.WrapIfErr("getspecification error", err) }()
	var data *Yields
	daysMax := 14
	for dayToSubstract := 1; dayToSubstract <= daysMax; dayToSubstract++ {

		// Проверяем условия выхода из цикла
		formatDate := date.Format(layout)
		// uri := fmt.Sprintf("https://iss.moex.com/iss/history/engines/stock/markets/bonds/sessions/3/securities/%s.json", ticker)
		Path := path.Join("iss", "history", "engines", "stock", "markets", "bonds", "sessions", "3", "securities", ticker+".json")
		params := url.Values{}
		params.Add("limit", "1")
		params.Add("iss.meta", "off")
		params.Add("history.columns", "TRADEDATE,MATDATE,OFFERDATE,BUYBACKDATE,YIELDCLOSE,YIELDTOOFFER,FACEVALUE, FACEUNIT,DURATION")
		params.Add("limit", "1")
		params.Add("from", formatDate)
		params.Add("to", formatDate)

		body, err := c.doRequest(Path, params)
		if err != nil {
			return data, err
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}
		if data.History != nil {

			if len(data.History.Data) != 0 {
				break
			}
		}
		date = date.AddDate(0, 0, -1)
	}

	return data, nil

}

func (d *Values) UnmarshalJSON(data []byte) error {
	dataSlice := make([]any, 8)
	err := json.Unmarshal(data, &dataSlice)
	if err != nil {
		return errors.New("CustomFloat64: UnmarshalJSON: " + err.Error())
	}
	d.TradeDate = checkStringNull(dataSlice[0])
	d.MaturityDate = checkStringNull(dataSlice[1])
	d.OfferDate = checkStringNull(dataSlice[2])
	d.BuybackDate = checkStringNull(dataSlice[3])
	d.YieldToMaturity = checkFloa64Null(dataSlice[4])
	d.YieldToOffer = checkFloa64Null(dataSlice[5])
	d.FaceValue = checkFloa64Null(dataSlice[6])
	d.FaceUnit = checkFloa64Null(dataSlice[7])
	d.Duration = checkFloa64Null(dataSlice[8])

	return nil
}

func checkFloa64Null(a any) *float64 {
	if FloatVal, ok := a.(float64); ok {
		return &FloatVal
	} else {
		return nil
	}
}

func checkStringNull(a any) *string {
	if StringVal, ok := a.(string); ok {
		return &StringVal
	} else {
		return nil
	}
}
