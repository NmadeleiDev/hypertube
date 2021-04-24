package torrentfile

func UnfoldArray(src [][]string) []string {
	res := make([]string, 0, len(src))

	for _, item := range src {
		res = append(res, item...)
	}

	return res
}

