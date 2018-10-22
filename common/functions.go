package common

func Map(vs []PlayerType, f func(PlayerType) bool) []bool {
	vsm := make([]bool, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
