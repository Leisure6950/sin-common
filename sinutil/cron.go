package sinutil

// TODO(sg): move file to right place
import (
	"fmt"
	"time"

	"gopkg.in/robfig/cron.v3"
)

var _ cron.Schedule = (*ConstantDelayAlignSchedule)(nil)

type ConstantDelayAlignSchedule struct {
	delay   time.Duration
	aligned bool
}

func EveryAlign(s string) (*ConstantDelayAlignSchedule, error) {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return &ConstantDelayAlignSchedule{}, err
	}

	if duration < time.Second {
		duration = time.Second
	}
	// precision: second
	duration -= time.Duration(duration.Nanoseconds()) % time.Second

	return &ConstantDelayAlignSchedule{
		delay: duration,
	}, nil
}

func EverySecondAlign(secs int) *ConstantDelayAlignSchedule {
	spec := fmt.Sprintf("%ds", secs)
	sched, _ := EveryAlign(spec)
	return sched
}

func (sched *ConstantDelayAlignSchedule) Next(t time.Time) time.Time {
	if !sched.aligned {
		sched.aligned = true
		t = t.Truncate(sched.delay)
	}

	return t.Add(sched.delay - time.Duration(t.Nanosecond())*time.Nanosecond)
}
