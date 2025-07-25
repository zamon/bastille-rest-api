package helpers

func RoundToTwoDecimals(val float64) float64 {
	return float64(int(val*100)) / 100
}