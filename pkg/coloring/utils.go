package coloring

func float64toUnit8(x float64) uint8 {
	if x > 255 {
		x = 255
	} else if x < 0 {
		x = 0
	}
	return uint8(x + 0.5)
}
