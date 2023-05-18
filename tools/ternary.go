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

func WrapDisk(output string) string {
	fields := strings.Fields(output)
	return fields[1] + "/" + fields[0] + "(" + fields[2] + ")"
}
