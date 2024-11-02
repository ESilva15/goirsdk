package conversions

func MsToKph(v float32) int {
	return int((3600 * v) / 1000)
}
