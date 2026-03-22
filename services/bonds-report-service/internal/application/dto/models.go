package dto

type BondReportsResponce struct {
	Media [][]*MediaGroup
}

type AccountListResponce struct {
	Accounts string
}

type MediaGroup struct {
	Reports []*ImageData
}

func NewMediaGroup() *MediaGroup {
	return &MediaGroup{
		Reports: make([]*ImageData, 0),
	}
}

type ImageData struct {
	Name    string
	Data    []byte
	Caption string
}

func NewImageData() *ImageData {
	return &ImageData{}
}

type UnionPortfolioStructureWithSberResponce struct {
	Report string
}
