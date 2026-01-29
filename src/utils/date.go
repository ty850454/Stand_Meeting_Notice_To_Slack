package utils

import (
	"fmt"
	"strconv"
	"time"
)

// GetCurrentIntDateAndYear 获取当前日期和年份
func GetCurrentIntDateAndYear() (int64, int) {
	now := time.Now()
	return TimeToIntDate(now), now.Year()
}

// ParseStrToTime dateStr 解析为date
func ParseStrToTime(dateStr string) (time.Time, error) {
	return time.Parse("20060102", dateStr)
}

// FormatIntDate 将int64格式的日期转换为yyyy-MM-dd格式的字符串
func FormatIntDate(date int64) string {
	dateStr := strconv.FormatInt(date, 10)
	if len(dateStr) != 8 {
		return dateStr
	}

	return fmt.Sprintf("%s-%s-%s", dateStr[:4], dateStr[4:6], dateStr[6:])
}

// TimeToIntDate 将日期解析为int64格式
func TimeToIntDate(date time.Time) int64 {

	return int64(date.Year()*10000 + int(date.Month())*100 + date.Day())
}

// IntDateToTime 将 int date 解析为 Time
func IntDateToTime(date int64) time.Time {

	year := date / 10000
	month := date / 100 % 100
	day := date % 100

	return time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.Local)
}
