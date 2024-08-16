package tools

import (
	"log"
	"time"
)

func ISO(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func ToTimeString(t int64) string {
	milli := time.UnixMilli(t)
	return milli.Format("2006-01-02 15:04:05")
}

func ToMill(t string) int64 {
	//2000-01-01 01:00:00
	parse, err := time.Parse("2006-01-02 15:04:05", t)
	if err != nil {
		log.Println(err)
		return 0
	}
	return parse.UnixMilli()
}

func ZeroTime() int64 {
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return date.UnixMilli()
}
