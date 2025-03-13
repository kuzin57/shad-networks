package convert

func FloatAnyToInt(val any) int {
	return int(val.(float64))
}
