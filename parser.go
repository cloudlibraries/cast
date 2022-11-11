package cast

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func parseInt(s string) (int64, error) {
	var foundZero bool
loop:
	for i := len(s); i > 0; i-- {
		switch s[i-1] {
		case '.':
			if foundZero {
				s = s[:i-1]
				break loop
			}
		case '0':
			foundZero = true
		default:
			break loop
		}
	}
	return strconv.ParseInt(s, 0, 0)
}

type timeFormatType int

const (
	timeFormatNoTimezone timeFormatType = iota
	timeFormatNamedTimezone
	timeFormatNumericTimezone
	timeFormatNumericAndNamedTimezone
	timeFormatTimeOnly
)

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[timeFormatNoTimezone-0]
	_ = x[timeFormatNamedTimezone-1]
	_ = x[timeFormatNumericTimezone-2]
	_ = x[timeFormatNumericAndNamedTimezone-3]
	_ = x[timeFormatTimeOnly-4]
}

const _timeFormatType_name = "timeFormatNoTimezonetimeFormatNamedTimezonetimeFormatNumericTimezonetimeFormatNumericAndNamedTimezonetimeFormatTimeOnly"

var _timeFormatType_index = [...]uint8{0, 20, 43, 68, 101, 119}

func (i timeFormatType) String() string {
	if i < 0 || i >= timeFormatType(len(_timeFormatType_index)-1) {
		return "timeFormatType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _timeFormatType_name[_timeFormatType_index[i]:_timeFormatType_index[i+1]]
}

type timeFormat struct {
	format string
	typ    timeFormatType
}

func (f timeFormat) hasTimezone() bool {
	// We don't include the formats with only named timezones, see
	// https://github.com/golang/go/issues/19694#issuecomment-289103522
	return f.typ >= timeFormatNumericTimezone && f.typ <= timeFormatNumericAndNamedTimezone
}

var (
	timeFormats = []timeFormat{
		{time.RFC3339, timeFormatNumericTimezone},
		{"2006-01-02T15:04:05", timeFormatNoTimezone}, // iso8601 without timezone
		{time.RFC1123Z, timeFormatNumericTimezone},
		{time.RFC1123, timeFormatNamedTimezone},
		{time.RFC822Z, timeFormatNumericTimezone},
		{time.RFC822, timeFormatNamedTimezone},
		{time.RFC850, timeFormatNamedTimezone},
		{"2006-01-02 15:04:05.999999999 -0700 MST", timeFormatNumericAndNamedTimezone}, // Time.String()
		{"2006-01-02T15:04:05-0700", timeFormatNumericTimezone},                        // RFC3339 without timezone hh:mm colon
		{"2006-01-02 15:04:05Z0700", timeFormatNumericTimezone},                        // RFC3339 without T or timezone hh:mm colon
		{"2006-01-02 15:04:05", timeFormatNoTimezone},
		{time.ANSIC, timeFormatNoTimezone},
		{time.UnixDate, timeFormatNamedTimezone},
		{time.RubyDate, timeFormatNumericTimezone},
		{"2006-01-02 15:04:05Z07:00", timeFormatNumericTimezone},
		{"2006-01-02", timeFormatNoTimezone},
		{"02 Jan 2006", timeFormatNoTimezone},
		{"2006-01-02 15:04:05 -07:00", timeFormatNumericTimezone},
		{"2006-01-02 15:04:05 -0700", timeFormatNumericTimezone},
		{time.Kitchen, timeFormatTimeOnly},
		{time.Stamp, timeFormatTimeOnly},
		{time.StampMilli, timeFormatTimeOnly},
		{time.StampMicro, timeFormatTimeOnly},
		{time.StampNano, timeFormatTimeOnly},
	}

	location = time.Local
)

func parseTime(s string) (time.Time, error) {
	return parseTimeWith(s, location, timeFormats)
}

func parseTimeWith(s string, location *time.Location, formats []timeFormat) (d time.Time, e error) {
	for _, format := range formats {
		if d, e = time.Parse(format.format, s); e == nil {

			// Some time formats have a zone name, but no offset, so it gets
			// put in that zone name (not the default one passed in to us), but
			// without that zone's offset. So set the location manually.
			if format.typ <= timeFormatNamedTimezone {
				if location == nil {
					location = time.Local
				}
				year, month, day := d.Date()
				hour, min, sec := d.Clock()
				d = time.Date(year, month, day, hour, min, sec, d.Nanosecond(), location)
			}

			return
		}
	}
	return d, fmt.Errorf("unable to parse date: %s", s)
}

var (
	durationRegExp       *regexp.Regexp
	durationRegExpGroups = []string{
		`<years>[\+|\-]?\d+y`,
		`<months>[\+|\-]?\d+M`,
		`<days>[\+|\-]?\d+d`,
		`<hours>[\+|\-]?\d+h`,
		`<minutes>[\+|\-]?\d+m`,
		`<seconds>[\+|\-]?\d+s`,
		`<milliseconds>[\+|\-]?\d+ms`,
		`<microseconds>[\+|\-]?\d+us`,
		`<nanoseconds>[\+|\-]?\d+ns`,
	}
)

func init() {
	var buf = new(bytes.Buffer)
	for _, group := range durationRegExpGroups {
		buf.WriteString(`(?P`)
		buf.WriteString(group)
		buf.WriteString(`)?`)
	}
	durationRegExp = regexp.MustCompile(buf.String())
}

func parseDuration(s string) (time.Duration, error) {
	return parseDurationWith(s, time.Now())
}

func parseDurationWith(s string, tm time.Time) (time.Duration, error) {
	matches := durationRegExp.FindStringSubmatch(s)
	if len(matches) == 0 {
		return 0, nil
	}

	nums := []int{}
	for index := 1; index < len(matches); index++ {
		s := matches[index]
		if len(s) == 0 {
			nums = append(nums, 0)
			continue
		}
		for s[len(s)-1] < '0' || s[len(s)-1] > '9' {
			s = s[:len(s)-1]
		}
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
		nums = append(nums, int(n))
	}

	if len(nums) == 0 {
		return 0, fmt.Errorf("parse duration `%s` failed", s)
	}

	duration := tm.AddDate(nums[0], nums[1], nums[2]).Add(
		time.Duration(nums[3]) * time.Hour,
	).Add(
		time.Duration(nums[4]) * time.Minute,
	).Add(
		time.Duration(nums[5]) * time.Second,
	).Add(
		time.Duration(nums[6]) * time.Millisecond,
	).Add(
		time.Duration(nums[7]) * time.Microsecond,
	).Add(
		time.Duration(nums[8]) * time.Nanosecond,
	).Sub(tm)

	return duration, nil
}
