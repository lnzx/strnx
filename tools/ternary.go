package tools

import (
	"strconv"
)

func IfThen(trueVal, defaultVal string) string {
	if trueVal == "" {
		return defaultVal
	}
	return trueVal
}

func ToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}
