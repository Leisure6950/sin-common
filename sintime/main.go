package sintime

import "time"

func Now() sinTime {
	return sinTime{
		Time: time.Now(),
	}
}
func Time(t time.Time) sinTime {
	return sinTime{Time: t}
}
func ParseDef(t string) (time.Time, error) {
	return sinTime{}.ParseDef(t)
}
func ParseDefDouble(t1, t2 string) (time.Time, time.Time, error) {
	return sinTime{}.ParseDefDouble(t1, t2)
}
func FormatDef(t time.Time) string {
	return sinTime{Time: t}.FormatDef()
}
