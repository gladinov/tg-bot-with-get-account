package service

import "bonds-report-service/internal/application/visualization"

type BondReportsResponce struct {
	Media [][]*visualization.MediaGroup
}
