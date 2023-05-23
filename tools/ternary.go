package tools

import (
	"strconv"
	"strings"
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

func ParseTraffic(s string) string {
	if s != "" && strings.Contains(s, "Not enough data") {
		return ""
	}
	return s
}
