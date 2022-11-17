package utils

import "time"

const (
	TimezoneDatePattern = "2006-01-02T15:04:05Z"
)

// 将UTC日期转换为日期字符串
func FormatDate(date time.Time, pattern string) string {
	return date.In(time.FixedZone("UTC", 0)).Format(pattern)
}

// 将UTC日期字符串转换为时间戳毫秒数
func ParseDate(dateStr string, pattern string) (int64, error) {
	date, err := time.Parse(pattern, dateStr)
	if err != nil {
		return -1, err
	}
	return date.UTC().UnixNano() / int64(time.Millisecond), nil
}
