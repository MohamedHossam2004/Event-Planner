package rabbit

func In(str string, list ...string) bool {
	for _, v := range list {
		if str == v {
			return true
		}
	}
	return false
}
