package torrentfile

func UnfoldArray(src [][]string) []string {
	res := make([]string, 0, len(src))

	for _, item := range src {
		res = append(res, item...)
	}

	return res
}

func StrArrayIdx(haystack []string, needle string) int {
	for i, item := range haystack {
		if item == needle {
			return i
		}
	}

	return -1
}
