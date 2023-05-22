package tools

import (
	"strconv"
)

func IfThen(exp bool, trueVal, defaultVal string) string {
	if exp {
		return trueVal
	}
	return defaultVal
}

func ToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}
