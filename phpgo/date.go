package phpgo

import (
	"strings"
	"time"
)

// DateStrToTime strtotime
func DateStrToTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	layouts := []string{
		"2006-01-02 15:04:05 -0700 MST",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05 -0700 MST",
		"2006/01/02 15:04:05 -0700",
		"2006/01/02 15:04:05",
		"2006-01-02 -0700 MST",
		"2006-01-02 -0700",
		"2006-01-02",
		"2006/01/02 -0700 MST",
		"2006/01/02 -0700",
		"2006/01/02",
		"2006-01-02 15:04:05 -0700 -0700",
		"2006/01/02 15:04:05 -0700 -0700",
		"2006-01-02 -0700 -0700",
		"2006/01/02 -0700 -0700",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}

	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, value)
		if err == nil {
			return t
		}
	}
	panic(err)
}

// DateFormat 跟 PHP 中 date 类似的使用方式，如果 ts 没传递，则使用当前时间
func DateFormat(format string, ts ...time.Time) string {
	patterns := []string{
		// 年
		"Y", "2006", // 4 位数字完整表示的年份
		"y", "06", // 2 位数字表示的年份

		// 月
		"m", "01", // 数字表示的月份，有前导零
		"n", "1", // 数字表示的月份，没有前导零
		"M", "Jan", // 三个字母缩写表示的月份
		"F", "January", // 月份，完整的文本格式，例如 January 或者 March

		// 日
		"d", "02", // 月份中的第几天，有前导零的 2 位数字
		"j", "2", // 月份中的第几天，没有前导零

		"D", "Mon", // 星期几，文本表示，3 个字母
		"l", "Monday", // 星期几，完整的文本格式;L的小写字母

		// 时间
		"g", "3", // 小时，12 小时格式，没有前导零
		"G", "15", // 小时，24 小时格式，没有前导零
		"h", "03", // 小时，12 小时格式，有前导零
		"H", "15", // 小时，24 小时格式，有前导零

		"a", "pm", // 小写的上午和下午值
		"A", "PM", // 小写的上午和下午值

		"i", "04", // 有前导零的分钟数
		"s", "05", // 秒数，有前导零
	}
	replacer := strings.NewReplacer(patterns...)
	format = replacer.Replace(format)

	t := time.Now()
	if len(ts) > 0 {
		t = ts[0]
	}
	return t.Format(format)
}

/*
DefaultFormat       = "2006-01-02 15:04:05"
DateFormat          = "2006-01-02"
FormattedDateFormat = "Jan 2, 2006"
TimeFormat          = "15:04:05"
HourMinuteFormat    = "15:04"
HourFormat          = "15"
DayDateTimeFormat   = "Mon, Aug 2, 2006 3:04 PM"
CookieFormat        = "Monday, 02-Jan-2006 15:04:05 MST"
RFC822Format        = "Mon, 02 Jan 06 15:04:05 -0700"
RFC1036Format       = "Mon, 02 Jan 06 15:04:05 -0700"
RFC2822Format       = "Mon, 02 Jan 2006 15:04:05 -0700"
RFC3339Format       = "2006-01-02T15:04:05-07:00"
RSSFormat           = "Mon, 02 Jan 2006 15:04:05 -0700"
*/

// DateDate date
func DateDate(format string, timestamp int64) string {
	// 时间戳小于等于0时，默认是当前时间
	if timestamp <= 0 {
		timestamp = time.Now().Unix()
	}
	now := time.Unix(timestamp, 0)

	// 原始的 format
	originFormat := format

	// 日
	// d
	if strings.ContainsRune(originFormat, rune('d')) {
		// -1 表示全部替换
		format = strings.Replace(format, "d", "02", -1)
	}
	// D
	if strings.ContainsRune(originFormat, rune('D')) {
		format = strings.Replace(format, "D", "Mon", -1)
	}
	// j
	if strings.ContainsRune(originFormat, rune('j')) {
		format = strings.Replace(format, "j", "2", -1)
	}
	// todo l
	if strings.ContainsRune(originFormat, rune('l')) {
		format = strings.Replace(format, "l", "Monday", -1)
	}
	// todo N
	if strings.ContainsRune(originFormat, rune('N')) {
		now.Weekday()
	}
	// todo S
	if strings.ContainsRune(originFormat, rune('S')) {

	}
	// todo w
	if strings.ContainsRune(originFormat, rune('w')) {
		now.Weekday()
	}
	// todo z
	if strings.ContainsRune(originFormat, rune('z')) {
		now.YearDay()
	}

	// 星期
	// todo W
	if strings.ContainsRune(originFormat, rune('W')) {

	}

	// 月
	// todo F
	if strings.ContainsRune(originFormat, rune('F')) {
		format = strings.Replace(format, "F", "January", -1)
	}
	// m
	if strings.ContainsRune(originFormat, rune('m')) {
		format = strings.Replace(format, "m", "01", -1)
	}
	// todo M
	if strings.ContainsRune(originFormat, rune('M')) {
		format = strings.Replace(format, "M", "Jan", -1)
	}
	// todo n
	if strings.ContainsRune(originFormat, rune('n')) {
		format = strings.Replace(format, "n", "1", -1)
	}
	// todo t
	if strings.ContainsRune(originFormat, rune('t')) {

	}

	// 年
	// todo L 判断是否为闰年
	if strings.ContainsRune(originFormat, rune('L')) {
	}
	// todo o
	if strings.ContainsRune(originFormat, rune('o')) {
	}
	// Y
	if strings.ContainsRune(originFormat, rune('Y')) {
		format = strings.Replace(format, "Y", "2006", -1)
	}
	if strings.ContainsRune(originFormat, rune('y')) {
		format = strings.Replace(format, "y", "06", -1)
	}

	// 时间
	// a
	if strings.ContainsRune(originFormat, rune('a')) {
		format = strings.Replace(format, "a", "pm", -1)
	}
	// A
	if strings.ContainsRune(originFormat, rune('A')) {
		format = strings.Replace(format, "A", "PM", -1)
	}
	// todo B
	if strings.ContainsRune(originFormat, rune('B')) {
	}
	// g
	if strings.ContainsRune(originFormat, rune('g')) {
		format = strings.Replace(format, "g", "3", -1)
	}
	// todo G
	if strings.ContainsRune(originFormat, rune('G')) {
	}
	// todo h
	if strings.ContainsRune(originFormat, rune('h')) {
	}
	// H
	if strings.ContainsRune(originFormat, rune('H')) {
		format = strings.Replace(format, "H", "15", -1)
	}
	// i
	if strings.ContainsRune(originFormat, rune('i')) {
		format = strings.Replace(format, "i", "04", -1)
	}
	// s
	if strings.ContainsRune(originFormat, rune('s')) {
		format = strings.Replace(format, "s", "05", -1)
	}
	// todo u
	if strings.ContainsRune(originFormat, rune('u')) {
	}

	// 时区
	// e
	if strings.ContainsRune(originFormat, rune('e')) {
		format = strings.Replace(format, "e", "MST", -1)
	}
	// todo I 判断是否为夏令时
	if strings.ContainsRune(originFormat, rune('I')) {
	}
	// O
	if strings.ContainsRune(originFormat, rune('O')) {
		format = strings.Replace(format, "O", "-0700", -1)
	}
	// P
	if strings.ContainsRune(originFormat, rune('P')) {
		format = strings.Replace(format, "P", "-07:00", -1)
	}
	// todo T
	if strings.ContainsRune(originFormat, rune('T')) {
	}
	// todo Z
	if strings.ContainsRune(originFormat, rune('T')) {
	}

	// 完整的时间/日期
	// c
	if strings.ContainsRune(originFormat, rune('c')) {
		format = strings.Replace(format, "c", "2006-01-02T15:04:05-07:00", -1)
	}
	// r
	if strings.ContainsRune(originFormat, rune('r')) {
		format = strings.Replace(format, "r", "Mon, 02 Jan 2006 15:04:05 -0700", -1)
	}
	// todo U
	if strings.ContainsRune(originFormat, rune('U')) {
		// format = strings.Replace(format, "U", fmt.Sprintf("%d", timestamp), -1)
	}

	return now.Format(format)
}

// DatemdHi m-d H:i
func DatemdHi() string {
	return time.Now().Format("01-02 15:04")
}

// DateYmdHi Y-m-d H:i
func DateYmdHi() string {
	return time.Now().Format("2006-01-02 15:04")
}

// DateYmd Y-m-d
func DateYmd() string {
	return time.Now().Format("2006-01-02")
}

// DateHis H:i:s
func DateHis() string {
	return time.Now().Format("15:04:05")
}

// DateYmdHis Y-m-d H:i:s
func DateYmdHis() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
