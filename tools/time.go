package tools

import (
	"strings"
	"time"
)

var (
	BackOffset = -1 * time.Millisecond
)

// GetBeforeDay 1945-10-10 12:12:12 --> [1945-10-9 12:12:12, 1945-10-10 12:12:12]
func GetBeforeDay(t time.Time) (time.Time, time.Time) {
	before := t.AddDate(0, 0, -1)
	return before, t
}

func GetBeforeDayN(t time.Time, n int) (time.Time, time.Time) {
	before := t.AddDate(0, 0, -n)
	return before, t
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

func GetDateString(timestamp string) string {
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
