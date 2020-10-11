package sintime

import (
	"time"
)

var (
	_timeZone, _timeOffset = time.Now().Zone()
)

type sinTime struct {
	time.Time
	err error
}

func (t sinTime) Error() error {
	return t.err
}

// 本地时区日起始时间
func (t sinTime) DayBegin() sinTime {
	return sinTime{Time: time.Unix(int64((t.Unix()+int64(_timeOffset))/(24*60*60))*24*60*60-int64(_timeOffset), 0)}
}

// 本地时区日结束时间
func (t sinTime) DayEnd() sinTime {
	return sinTime{Time: time.Unix(int64((t.Unix()+int64(_timeOffset))/(24*60*60))*24*60*60-int64(_timeOffset)+60*60*24-1, 0)}
}

// 本地时区日相对时间
func (t sinTime) DayRelativeSecond() int64 {
	return t.Unix() - (int64((t.Unix()+int64(_timeOffset))/(24*60*60))*24*60*60 - int64(_timeOffset))
}

// 以本地时区格式化时间为字符串
func (t sinTime) FormatDef() string {
	return t.Format("2006-01-02 15:04:05")
}

// 以本地时区解析时间字符串
func (t sinTime) ParseDef(v string) (time.Time, error) {
	if t.err != nil {
		return time.Time{}, t.err
	}
	var tmp time.Time
	switch len(v) {
	case 10:
		tmp, t.err = time.ParseInLocation("2006-01-02", v, time.Local)
	case 13:
		tmp, t.err = time.ParseInLocation("2006-01-02 15", v, time.Local)
	case 16:
		tmp, t.err = time.ParseInLocation("2006-01-02 15:04", v, time.Local)
	case 19:
		tmp, t.err = time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
	case 23:
		tmp, t.err = time.ParseInLocation("2006-01-02 15:04:05.000", v, time.Local)
	default:
		tmp, t.err = time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
	}
	return tmp, t.err
}

// 以本地时区解析两个时间字符串
func (t sinTime) ParseDefDouble(t1, t2 string) (time.Time, time.Time, error) {
	nvt1, err := t.ParseDef(t1)
	if err != nil {
		return nvt1, nvt1, err
	}
	nvt2, err := t.ParseDef(t2)
	if err != nil {
		return nvt1, nvt2, err
	}
	return nvt1, nvt2, nil
}
