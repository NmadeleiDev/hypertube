package p2p

func IntArrayContain(haystack []int, needle int) bool {
	if haystack == nil {
		return false
	}

	for _, val := range haystack {
		if val == needle {
			return true
		}
	}
	return false
}
