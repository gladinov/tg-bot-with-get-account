package other_func

import (
	"math"

	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

func CastMoney(v *pb.Quotation) float64 {
	if v != nil {
		r := float64(v.Units) + float64(v.Nano/1e9)
		return r
	}
	return 0
}

func MoneyValue(v *pb.MoneyValue) float64 {
	if v != nil {
		r := float64(v.Units) + float64(v.Nano/1e9)
		return r
	}
	return 0
}

// округление до precision знаков после запятой
func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
