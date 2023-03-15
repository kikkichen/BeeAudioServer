package utils

import (
	"time"
)

const (
	formatTime = "Mon Jan 02 15:04:05 +0800 2006"
	StartTime  = "Tue Dec 01 12:30:00 +0800 2022" // 最早时间
)

/**
*	整理时间格式
*
*	@params	timeStr	:	时间字符串
*
 */
func StrToDate(
	timeStr string,
) time.Time {
	parse, _ := time.ParseInLocation(formatTime, timeStr, time.Local)
	return parse
}

// 秒级时间戳转time
func UnixSecondToTime(second int64) time.Time {
	return time.Unix(second, 0)
}

// 毫秒级时间戳转time
func UnixMilliToTime(milli int64) time.Time {
	return time.Unix(milli/1000, (milli%1000)*(1000*1000))
}

// 纳秒级时间戳转time
func UnixNanoToTime(nano int64) time.Time {
	return time.Unix(nano/(1000*1000*1000), nano%(1000*1000*1000))
}
