package time

import "time"

func GetCurrentTime() time.Time {
	return time.Now()
}

// 时间减法
func SubTime(t time.Time, d time.Duration) time.Time {
	return t.Add(-d)
}

// 当前时间前n分钟的时间
func GetTimeBeforeMinutes(minutes int) time.Time {
	return SubTime(GetCurrentTime(), time.Duration(minutes)*time.Minute)
}
