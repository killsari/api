package utils

import "time"

func ParseInLocation(s string) (t time.Time, err error) {
	return time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
}

func ParseString(t time.Time) (s string) {
	return t.Format("2006-01-02 15:04:05")
}
