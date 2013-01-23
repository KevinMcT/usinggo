package custom

func Mystery() bool {
	someMap := make(map[int]string)
	someMap[0] = "String"
	_, ok := someMap[0]
	return ok
}
