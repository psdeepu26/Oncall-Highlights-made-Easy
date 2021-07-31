package controllers

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func getEpochTime(t time.Time) int64 {
	return t.Unix()
}

func GetDiffEpochTime(diff_days int) int64 {
	t := time.Now()
	t2 := t.AddDate(0, 0, -diff_days)
	return getEpochTime(t2)
}

func GetParsedTime(timeChange string) time.Time {
	layout := "2006-01-02T15:04:05Z07:00"
	t, _ := time.Parse(layout, timeChange)

	return t
}

func GetDiffTime(t time.Time, timeToDiff time.Time) time.Duration {
	return t.Sub(timeToDiff).Round(time.Minute)
}

func ShortDur(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}

func GetCurrentDateTime() string {
	dt := time.Now()

	return dt.Format("01-02-2006 15:04:05 Monday")
}

func GetCurrentDateTimeInPST() string {
	pst, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Error(err)
	}

	dt := time.Now().UTC().In(pst)

	return dt.Format("01-02-2006 15:04:05 Monday")
}
