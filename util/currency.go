package util

const (
	USD = "USD"
	THB = "THB"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, THB:
		return true
	}
	return false
}
