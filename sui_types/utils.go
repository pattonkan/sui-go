package sui_types

import (
	"strconv"

	"github.com/shopspring/decimal"
)

type SUI float64

func (s SUI) Int64() int64 {
	return int64(s * 1e9)
}
func (s SUI) Uint64() uint64 {
	return uint64(s * 1e9)
}
func (s SUI) Decimal() decimal.Decimal {
	return decimal.NewFromInt(s.Int64())
}
func (s SUI) String() string {
	return strconv.FormatInt(s.Int64(), 10)
}
