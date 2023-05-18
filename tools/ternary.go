package tools

import (
	"fmt"
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
	size := fields[0]
	if len(size) < 4 {
		size = fmt.Sprintf("%s%s", strings.Repeat(" ", 4-len(size)), size)
	}
	used := fields[1]
	if len(used) < 4 {
		used = fmt.Sprintf("%s%s", strings.Repeat(" ", 4-len(used)), used)
	}
	percent := "(" + fields[2] + ")"
	if len(percent) < 6 {
		percent = fmt.Sprintf("%s%s", strings.Repeat(" ", 6-len(percent)), percent)
	}
	return used + "/" + size + percent
}
