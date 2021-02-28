package decimals

import "github.com/shopspring/decimal"

func NewBalance(balance, amount float64) float64 {
	res, _ := decimal.NewFromFloat(balance).Add(decimal.NewFromFloat(amount)).Round(3).Float64()
	return res
}

func NegativeAmount(amount float64) float64 {
	negative := decimal.NewFromInt(-1)
	result, _ := negative.Mul(decimal.NewFromFloat(amount)).Float64()
	return result
}
