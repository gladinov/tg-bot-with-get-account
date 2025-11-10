package service

import "time"

type Service struct {
	Client *Client
}

func NewService(client *Client) *Service {
	return &Service{Client: client}
}

func (s *Service) GetAllCurrencies(date time.Time) (ValCurs, error) {
	return s.Client.GetAllCurrencies(date)
}
