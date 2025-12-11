package bondreportservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"main.go/internal/models"
	"main.go/lib/valuefromcontext"
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

func (c *Client) GetAccountsList(ctx context.Context) (AccountListResponce, error) {
	const op = "bondreportservice.GetAccountsList"
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return AccountListResponce{}, fmt.Errorf("%s: %w", op, err)
	}
	pth := path.Join("bondReportService", "accounts")
	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   pth,
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return AccountListResponce{}, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set(models.HeaderChatID, chatID)

	resp, err := c.client.Do(req)
	if err != nil {
		return AccountListResponce{}, fmt.Errorf("%s: %w", op, err)
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return AccountListResponce{}, fmt.Errorf("%s: %w", op, err)
	}
	if resp.StatusCode != http.StatusOK {
		var statusErr map[string]string
		err = json.Unmarshal(body, &statusErr)
		if err != nil {
			return AccountListResponce{}, fmt.Errorf("%s: %w", op, err)
		}
		return AccountListResponce{}, fmt.Errorf("%s:"+statusErr["error"], op)
	}

	var accountResponce AccountListResponce
	err = json.Unmarshal(body, &accountResponce)
	if err != nil {
		return AccountListResponce{}, fmt.Errorf("%s: %w", op, err)
	}
	return accountResponce, nil

}

func (c *Client) GetUsd(ctx context.Context) (UsdResponce, error) {
	const op = "bondreportservice.GetUsd"
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return UsdResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	pth := path.Join("bondReportService", "getUSD")

	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   pth,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return UsdResponce{}, fmt.Errorf("%s:%w", op, err)
	}

	req.Header.Set(models.HeaderChatID, chatID)

	resp, err := c.client.Do(req)
	if err != nil {
		return UsdResponce{}, fmt.Errorf("%s:%w", op, err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UsdResponce{}, fmt.Errorf("%s:%w", op, err)
	}

	if resp.StatusCode != http.StatusOK {
		var statusErr map[string]string
		err := json.Unmarshal(body, &statusErr)
		if err != nil {
			return UsdResponce{}, fmt.Errorf("%s:%w", op, err)
		}
		return UsdResponce{}, fmt.Errorf("%s:"+statusErr["error"], op)
	}

	var usdResponce UsdResponce
	err = json.Unmarshal(body, &usdResponce)
	if err != nil {
		return UsdResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	return usdResponce, nil
}

func (c *Client) GetBondReportsByFifo(ctx context.Context) error {
	const op = "bondreportservice.GetBondReportsByFifo"
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	pth := path.Join("bondReportService", "getBondReportsByFifo")
	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   pth,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	req.Header.Set(models.HeaderChatID, chatID)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	if resp.StatusCode != http.StatusNoContent {
		var statusErr map[string]string
		err := json.Unmarshal(body, &statusErr)
		if err != nil {
			return fmt.Errorf("%s:%w", op, err)
		}
		return fmt.Errorf("%s:"+statusErr["error"], op)
	}
	return nil

}

func (c *Client) GetBondReports(ctx context.Context) (BondReportsResponce, error) {
	const op = "bondreportservice.GetBondReports"
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return BondReportsResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	pth := path.Join("bondReportService", "getBondReports")
	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   pth,
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return BondReportsResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	req.Header.Set(models.HeaderChatID, chatID)

	resp, err := c.client.Do(req)
	if err != nil {
		return BondReportsResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return BondReportsResponce{}, fmt.Errorf("%s:%w", op, err)
	}

	if resp.StatusCode != http.StatusOK {
		var statusErr map[string]string
		err := json.Unmarshal(body, &statusErr)
		if err != nil {
			return BondReportsResponce{}, fmt.Errorf("%s:%w", op, err)
		}
		return BondReportsResponce{}, fmt.Errorf("%s:"+statusErr["error"], op)
	}
	var bondReportResponce BondReportsResponce
	err = json.Unmarshal(body, &bondReportResponce)
	if err != nil {
		return BondReportsResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	return bondReportResponce, nil
}

func (c *Client) GetPortfolioStructure(ctx context.Context) (PortfolioStructureForEachAccountResponce, error) {
	const op = "bondreportservice.GetPortfolioStructure"
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return PortfolioStructureForEachAccountResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	pth := path.Join("bondReportService", "getPortfolioStructure")
	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   pth,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return PortfolioStructureForEachAccountResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	req.Header.Set(models.HeaderChatID, chatID)

	resp, err := c.client.Do(req)
	if err != nil {
		return PortfolioStructureForEachAccountResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PortfolioStructureForEachAccountResponce{}, fmt.Errorf("%s:%w", op, err)
	}

	if resp.StatusCode != http.StatusOK {
		var statusErr map[string]string
		err := json.Unmarshal(body, &statusErr)
		if err != nil {
			return PortfolioStructureForEachAccountResponce{}, fmt.Errorf("%s:%w", op, err)
		}
		return PortfolioStructureForEachAccountResponce{}, fmt.Errorf("%s:"+statusErr["error"], op)
	}
	var bondReportResponce PortfolioStructureForEachAccountResponce
	err = json.Unmarshal(body, &bondReportResponce)
	if err != nil {
		return PortfolioStructureForEachAccountResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	return bondReportResponce, nil

}

func (c *Client) GetUnionPortfolioStructure(ctx context.Context) (UnionPortfolioStructureResponce, error) {
	const op = "bondreportservice.GetUnionPortfolioStructure"
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return UnionPortfolioStructureResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	pth := path.Join("bondReportService", "getUnionPortfolioStructure")
	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   pth,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return UnionPortfolioStructureResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	req.Header.Set(models.HeaderChatID, chatID)

	resp, err := c.client.Do(req)
	if err != nil {
		return UnionPortfolioStructureResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UnionPortfolioStructureResponce{}, fmt.Errorf("%s:%w", op, err)
	}

	if resp.StatusCode != http.StatusOK {
		var statusErr map[string]string
		err := json.Unmarshal(body, &statusErr)
		if err != nil {
			return UnionPortfolioStructureResponce{}, fmt.Errorf("%s:%w", op, err)
		}
		return UnionPortfolioStructureResponce{}, fmt.Errorf("%s:"+statusErr["error"], op)
	}
	var bondReportResponce UnionPortfolioStructureResponce
	err = json.Unmarshal(body, &bondReportResponce)
	if err != nil {
		return UnionPortfolioStructureResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	return bondReportResponce, nil

}

func (c *Client) GetUnionPortfolioStructureWithSber(ctx context.Context) (UnionPortfolioStructureWithSberResponce, error) {
	const op = "bondreportservice.GetUnionPortfolioStructureWithSber"
	chatID, err := valuefromcontext.GetChatIDFromCtxStr(ctx)
	if err != nil {
		return UnionPortfolioStructureWithSberResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	pth := path.Join("bondReportService", "getUnionPortfolioStructureWithSber")
	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   pth,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return UnionPortfolioStructureWithSberResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	req.Header.Set(models.HeaderChatID, chatID)

	resp, err := c.client.Do(req)
	if err != nil {
		return UnionPortfolioStructureWithSberResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UnionPortfolioStructureWithSberResponce{}, fmt.Errorf("%s:%w", op, err)
	}

	if resp.StatusCode != http.StatusOK {
		var statusErr map[string]string
		err := json.Unmarshal(body, &statusErr)
		if err != nil {
			return UnionPortfolioStructureWithSberResponce{}, fmt.Errorf("%s:%w", op, err)
		}
		return UnionPortfolioStructureWithSberResponce{}, fmt.Errorf("%s:"+statusErr["error"], op)
	}
	var bondReportResponce UnionPortfolioStructureWithSberResponce
	err = json.Unmarshal(body, &bondReportResponce)
	if err != nil {
		return UnionPortfolioStructureWithSberResponce{}, fmt.Errorf("%s:%w", op, err)
	}
	return bondReportResponce, nil

}
