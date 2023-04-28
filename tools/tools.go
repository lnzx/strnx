package tools

import (
	"math/rand"
	"strings"
	"time"
)

// RandStringBytesMaskImprSrc https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func RandStringBytesMaskImprSrc(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func IfThen(trueVal, defaultVal string) string {
	if trueVal == "" {
		return defaultVal
	}
	return trueVal
}

var (
	BackOffset = -1 * time.Millisecond
)

// GetBeforeDay 1945-10-10 12:12:12 --> 1945-10-9 12:12:12
func GetBeforeDay(t time.Time) time.Time {
	before := t.AddDate(0, 0, -1)
	return before
}

// GetDayStart 1945-10-10 12:12:12 --> 1945-10-10 00:00:00
func GetDayStart(t time.Time) time.Time {
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return start
}

// GetMonthRange 1945-10-10 12:12:12 --> [1945-10-01 00:00:00, 1945-10-30 23:59:59]
func GetMonthRange(t time.Time) (time.Time, time.Time) {
	today := GetDayStart(t)
	dayOffset := t.Day()*-1 + 1
	start := today.AddDate(0, 0, dayOffset)
	end := start.AddDate(0, 1, 0).Add(BackOffset)
	return start, end
}

func GetDate(timestamp string) string {
	return strings.Split(timestamp, "T")[0]
}

// GetMonthRangeDate 1945-10-10 12:12:12 --> [1945-10-01, 1945-10-31]
func GetMonthRangeDate(t time.Time) (string, string) {
	start, end := GetMonthRange(t)
	return DateFormat(start), DateFormat(end)
}

func DateFormat(t time.Time) string {
	return t.Format("2006-01-02")
}
