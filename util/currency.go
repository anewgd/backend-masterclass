package util

const (
	USD = "USD"
	EUR = "EUR"
	ETB = "ETB"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, ETB:
		return true
	}

	return false
}
