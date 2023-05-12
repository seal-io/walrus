package timex

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-module/carbon"
)

const (
	Year    = "year"
	Quarter = "quarter"
	Month   = "month"
	Week    = "week"
	Day     = "day"
)

// TimezoneInPosix is in posix timezone string format.
// Time zone Asia/Shanghai in posix is UTC-8.
func TimezoneInPosix(offset int) string {
	timeZone := "UTC"
	if offset != 0 {
		utcOffSig := "-"
		utcOffHrs := offset / 60 / 60

		if utcOffHrs < 0 {
			utcOffSig = "+"
			utcOffHrs = 0 - utcOffHrs
		}

		timeZone = fmt.Sprintf("UTC%s%d", utcOffSig, utcOffHrs)
	}
	return timeZone
}

// StartTimeOfHour returns the start time of the hour.
func StartTimeOfHour(t time.Time, loc *time.Location) time.Time {
	return carbon.FromStdTime(t).SetLocation(loc).StartOfHour().ToStdTime()
}

// StartTimeOfNextHour returns the start time of the next hour.
func StartTimeOfNextHour(t time.Time, loc *time.Location) time.Time {
	return carbon.FromStdTime(t).AddHours(1).SetLocation(loc).StartOfHour().ToStdTime()
}

func StartOfDay(t time.Time, loc *time.Location) time.Time {
	return carbon.FromStdTime(t).SetLocation(loc).StartOfDay().ToStdTime()
}

func StartOfWeek(t time.Time, loc *time.Location) time.Time {
	return carbon.FromStdTime(t).SetLocation(loc).StartOfWeek().ToStdTime()
}

func StartOfMonth(t time.Time, loc *time.Location) time.Time {
	return carbon.FromStdTime(t).SetLocation(loc).StartOfMonth().ToStdTime()
}

func StartOfQuarter(t time.Time, loc *time.Location) time.Time {
	return carbon.FromStdTime(t).SetLocation(loc).StartOfQuarter().ToStdTime()
}

func StartOfYear(t time.Time, loc *time.Location) time.Time {
	return carbon.FromStdTime(t).SetLocation(loc).StartOfYear().ToStdTime()
}

// GetTimeSeries group by step (day, week, month, quarter, year) with time range,
// return the time series.
func GetTimeSeries(start, end time.Time, step string, loc *time.Location) ([]time.Time, error) {
	var timeSeries []time.Time
	switch step {
	case Day:
		start = StartOfDay(start, loc)
		startTimeOfEnd := StartOfDay(end, loc)
		if !startTimeOfEnd.Equal(end) {
			end = startTimeOfEnd
		} else {
			end = startTimeOfEnd.AddDate(0, 0, -1)
		}
		for !start.After(end) {
			timeSeries = append(timeSeries, StartOfDay(start, loc))
			start = start.AddDate(0, 0, 1)
		}
	case Week:
		start = StartOfWeek(start, loc)
		startTimeOfEnd := StartOfWeek(end, loc)
		if !startTimeOfEnd.Equal(end) {
			end = startTimeOfEnd
		} else {
			end = startTimeOfEnd.AddDate(0, 0, -7)
		}
		for !start.After(end) {
			timeSeries = append(timeSeries, StartOfWeek(start, loc))
			start = start.AddDate(0, 0, 7)
		}
	case Month:
		start = StartOfMonth(start, loc)
		startTimeOfEnd := StartOfMonth(end, loc)
		if !startTimeOfEnd.Equal(end) {
			end = startTimeOfEnd
		} else {
			end = startTimeOfEnd.AddDate(0, -1, 0)
		}
		for !start.After(end) {
			timeSeries = append(timeSeries, StartOfMonth(start, loc))
			start = start.AddDate(0, 1, 0)
		}
	case Quarter:
		start = StartOfQuarter(start, loc)
		startTimeOfEnd := StartOfQuarter(end, loc)
		if !startTimeOfEnd.Equal(end) {
			end = startTimeOfEnd
		} else {
			end = startTimeOfEnd.AddDate(0, -3, 0)
		}
		for !start.After(end) {
			timeSeries = append(timeSeries, StartOfQuarter(start, loc))
			start = start.AddDate(0, 3, 0)
		}
	case Year:
		start = StartOfYear(start, loc)
		startTimeOfEnd := StartOfYear(end, loc)
		if !startTimeOfEnd.Equal(end) {
			end = startTimeOfEnd
		} else {
			end = startTimeOfEnd.AddDate(-1, 0, 0)
		}
		for !start.After(end) {
			timeSeries = append(timeSeries, StartOfYear(start, loc))
			start = start.AddDate(1, 0, 0)
		}
	default:
		return nil, errors.New("invalid time unit")
	}

	return timeSeries, nil
}
