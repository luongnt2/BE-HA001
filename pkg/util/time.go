package util

import (
	"fmt"
	"time"
)

const (
	DefaultFormatDate = "02/01/2006"
	DefaultFormatTime = "15:04"
	DefaultFormatYear = "2006"
)

func TimeToUnix(t *time.Time) int64 {
	if t != nil {
		return t.Unix()
	}

	return 0
}

func GetVietnameseDateLabel(t time.Time) string {
	if DateEqual(t, time.Now()) {
		return "hôm nay"
	}

	return fmt.Sprintf("ngày %s", t.Format("02/01/2006"))
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func GetGMT7TimeZone() *time.Location {
	timezone, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		panic(fmt.Sprintf("time zone error: %v", err))
	}

	return timezone
}

func NewTimePointer(t time.Time) *time.Time {
	return &t
}
