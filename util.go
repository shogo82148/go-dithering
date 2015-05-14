package dithering

// clamp clamps i to the interval [0, 0xffff].
func clamp(i int32) int32 {
	if i < 0 {
		return 0
	}
	if i > 0xffff {
		return 0xffff
	}
	return i
}
