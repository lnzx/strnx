package tools

func IfThen(trueVal, defaultVal string) string {
	if trueVal == "" {
		return defaultVal
	}
	return trueVal
}
