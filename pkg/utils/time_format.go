/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/2/22 11:52).
 */
package utils

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	TimeOffset = 8 * 3600  //8 hour offset
	HalfOffset = 12 * 3600 //Half-day hourly offset
)

// Get the current timestamp by Second
func GetCurrentTimestampBySecond() int64 {
	return time.Now().Unix()
}

// Convert timestamp to time.Time type
func UnixSecondToTime(second int64) time.Time {
	return time.Unix(second, 0)
}

// Convert nano timestamp to time.Time type
func UnixNanoSecondToTime(nanoSecond int64) time.Time {
	return time.Unix(0, nanoSecond)
}
func UnixMillSecondToTime(millSecond int64) time.Time {
	return time.Unix(0, millSecond*1e6)
}

// Get the current timestamp by Nano
func GetCurrentTimestampByNano() int64 {
	return time.Now().UnixNano()
}

// Get the current timestamp by Mill
func GetCurrentTimestampByMill() int64 {
	return time.Now().UnixNano() / 1e6
}

// Get the timestamp at 0 o'clock of the day
func GetCurDayZeroTimestamp() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	return t.Unix() - TimeOffset
}

// Get the timestamp at 12 o'clock on the day
func GetCurDayHalfTimestamp() int64 {
	return GetCurDayZeroTimestamp() + HalfOffset

}

// Get the formatted time at 0 o'clock of the day, the format is "2006-01-02_00-00-00"
func GetCurDayZeroTimeFormat() string {
	return time.Unix(GetCurDayZeroTimestamp(), 0).Format("2006-01-02_15-04-05")
}

// Get the formatted time at 12 o'clock of the day, the format is "2006-01-02_12-00-00"
func GetCurDayHalfTimeFormat() string {
	return time.Unix(GetCurDayZeroTimestamp()+HalfOffset, 0).Format("2006-01-02_15-04-05")
}
func GetTimeStampByFormat(datetime string) string {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	timestamp := tmp.Unix()
	return strconv.FormatInt(timestamp, 10)
}

func TimeStringFormatTimeUnix(timeFormat string, timeSrc string) int64 {
	tm, _ := time.Parse(timeFormat, timeSrc)
	return tm.Unix()
}

func TimeStringToTime(timeString string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", timeString)
	return t, err
}

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02")
}

func GetDateTimeBeginTimeAndEndTime() (time.Time, time.Time) {
	//1.获取当前时区
	loc, _ := time.LoadLocation("Local")
	loc = loc

	//2.今日日期字符串
	date := time.Now().Format("2006-01-02")

	//3.拼接成当天0点时间字符串
	startDate := date + " 00:00:00.000"
	//得到0点日期 2021-04-24 00:00:00 +0800 CST
	startTime, _ := time.Parse("2006-01-02 15:04:05.000", startDate)
	fmt.Println(startTime)
	//4.拼接成当天23点时间字符串
	endDate := date + " 23:59:59.999"
	//得到23点日期 2021-04-24 23:59:59 +0800 CST
	endTime, _ := time.Parse("2006-01-02 15:04:05.000", endDate)
	fmt.Println(endTime)

	return startTime, endTime
}

func GetDateTimeBeginTimeAndEndTimeByInputTime(time2 time.Time) (time.Time, time.Time) {

	//2.今日日期字符串
	date := time2.Format("2006-01-02")

	//3.拼接成当天0点时间字符串
	startDate := date + " 00:00:00.000"
	//得到0点日期 2021-04-24 00:00:00 +0800 CST
	startTime, _ := time.Parse("2006-01-02 15:04:05.000", startDate)
	fmt.Println(startTime)
	//4.拼接成当天23点时间字符串
	endDate := date + " 23:59:59.999"
	//得到23点日期 2021-04-24 23:59:59 +0800 CST
	endTime, _ := time.Parse("2006-01-02 15:04:05.000", endDate)
	fmt.Println(endTime)

	return startTime, endTime
}

// 取当前时间到其他时间的差值
func SubDemo(ts string) (error, time.Duration) {
	now := time.Now()
	_, err := time.Parse("2006-01-02 15:04:05", ts)
	if err != nil {
		fmt.Printf("parse string err:%v\n", err)
		return err, 0
	}
	// 按照东八区的时区格式解析一个字符串
	tlocal, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Printf("Parse a string according to the time zone format of Dongba district err:%v\n", err)
		return err, 0
	}
	// 按照指定的时区解析时间
	t, err := time.ParseInLocation("2006-01-02 15:04:05", ts, tlocal)
	if err != nil {
		fmt.Printf("Resolve the time according to the specified time zone:%v\n", err)
		return err, 0
	}
	// 计算时间的差值
	reverseTime := now.Sub(t)
	return nil, reverseTime
}

const (
	timePattern = `(\d{4}[-/\.]\d{1,2}[-/\.]\d{1,2})[:\sT-]*(\d{0,2}:{0,1}\d{0,2}:{0,1}\d{0,2}){0,1}\.{0,1}(\d{0,9})([\sZ]{0,1})([\+-]{0,1})([:\d]*)`
)

var (
	TimeTpl      = "2006-01-02 15:04:05"
	formatKeyTpl = map[byte]string{
		'd': "02",
		'D': "Mon",
		'w': "Monday",
		'N': "Monday",
		'S': "02",
		'l': "Monday",
		'F': "January",
		'm': "01",
		'M': "Jan",
		'n': "1",
		'Y': "2006",
		'y': "06",
		'a': "pm",
		'A': "PM",
		'g': "3",
		'h': "03",
		'H': "15",
		'i': "04",
		's': "05",
		'O': "-0700",
		'P': "-07:00",
		'T': "MST",
		'u': "000000",
		'c': "2006-01-02T15:04:05-07:00",
		'r': "Mon, 02 Jan 06 15:04 MST",
	}
	GetLocationName = func(zone int) string {
		switch zone {
		case 8:
			return "Asia/Shanghai"
		}
		return "UTC"
	}
)

// FormatTlp format template
func FormatTlp(format string) string {
	runes := []rune(format)
	buffer := bytes.NewBuffer(nil)
	for i := 0; i < len(runes); i++ {
		switch runes[i] {
		case '\\':
			if i < len(runes)-1 {
				buffer.WriteRune(runes[i+1])
				i += 1
				continue
			} else {
				return buffer.String()
			}
		default:
			if runes[i] > 255 {
				buffer.WriteRune(runes[i])
				break
			}
			if f, ok := formatKeyTpl[byte(runes[i])]; ok {
				buffer.WriteString(f)
			} else {
				buffer.WriteRune(runes[i])
			}
		}
	}
	return buffer.String()
}

type regexMapStruct struct {
	Value *regexp.Regexp
	Time  int64
	sync.RWMutex
}

var (
	l                 sync.RWMutex
	regexCache             = map[string]*regexMapStruct{}
	regexCacheTimeout uint = 1800
)

func getRegexpCompile(pattern string) (r *regexp.Regexp, err error) {
	l.RLock()
	var data *regexMapStruct
	var ok bool
	data, ok = regexCache[pattern]
	l.RUnlock()
	if ok {
		r = data.Value
		return
	}
	r, err = regexp.Compile(pattern)
	if err != nil {
		return
	}
	l.Lock()
	regexCache[pattern] = &regexMapStruct{Value: r, Time: time.Now().Unix()}
	l.Unlock()
	return
}

// RegexExtract extract matching text
func RegexExtract(pattern string, str string) ([]string, error) {
	r, err := getRegexpCompile(pattern)
	if err == nil {
		return r.FindStringSubmatch(str), nil
	}
	return nil, err
}
func init() {
	go func() {
		ticker := time.NewTicker(600 * time.Second)
		for range ticker.C {
			clearRegexpCompile()
		}
	}()
}
func clearRegexpCompile() {
	newRegexCache := map[string]*regexMapStruct{}
	l.Lock()
	defer l.Unlock()
	if len(regexCache) == 0 {
		return
	}
	now := time.Now().Unix()
	for k := range regexCache {
		if uint(now-regexCache[k].Time) <= regexCacheTimeout {
			newRegexCache[k] = &regexMapStruct{Value: regexCache[k].Value, Time: now}
		}
	}
	regexCache = newRegexCache
}

const (
	// PadRight Right padding character
	PadRight PadType = iota
	// PadLeft Left padding character
	PadLeft
	// PadSides Two-sided padding characters,If the two sides are not equal, the right side takes precedence.
	PadSides
	letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type (
	// ru is a pseudorandom number generator
	ru struct {
		x uint32
	}
	PadType uint8
)

// Len string length (utf8)
func Len(str string) int {
	// strings.Count(str,"")-1
	return utf8.RuneCountInString(str)
}

// Pad String padding
func Pad(raw string, length int, padStr string, padType PadType) string {
	l := length - Len(raw)
	if l <= 0 {
		return raw
	}
	if padType == PadRight {
		raw = fmt.Sprintf("%s%s", raw, strings.Repeat(padStr, l))
	} else if padType == PadLeft {
		raw = fmt.Sprintf("%s%s", strings.Repeat(padStr, l), raw)
	} else {
		left := 0
		right := 0
		if l > 1 {
			left = l / 2
			right = (l / 2) + (l % 2)
		}

		raw = fmt.Sprintf("%s%s%s", strings.Repeat(padStr, left), raw, strings.Repeat(padStr, right))
	}
	return raw
}
func ParseTime(str string, format ...string) (t time.Time, err error) {
	if len(format) > 0 {
		return time.ParseInLocation(FormatTlp(format[0]), str, time.Local)
	}

	var year, month, day, hour, min, sec string
	match, err := RegexExtract(timePattern, str)
	if err != nil {
		return
	}
	matchLen := len(match)
	if matchLen == 0 {
		err = errors.New("cannot parse")
		return
	}
	if matchLen > 1 && match[1] != "" {
		for k, v := range match {
			match[k] = strings.TrimSpace(v)
		}
		arr := make([]string, 3)
		for _, v := range []string{"-", "/", "."} {
			arr = strings.Split(match[1], v)
			if len(arr) >= 3 {
				break
			}
		}
		if len(arr) < 3 {
			err = errors.New("cannot parse date")
			return
		}
		year = arr[0]
		month = Pad(arr[1], 2, "0", PadLeft)
		day = Pad(arr[2], 2, "0", PadLeft)
	}

	if len(match[2]) > 0 {
		s := strings.Replace(match[2], ":", "", -1)
		if len(s) < 6 {
			s += strings.Repeat("0", 6-len(s))
		}
		hour = Pad(s[0:2], 2, "0", PadLeft)
		min = Pad(s[2:4], 2, "0", PadLeft)
		sec = Pad(s[4:6], 2, "0", PadLeft)
		return time.ParseInLocation(TimeTpl, year+"-"+month+"-"+day+" "+hour+":"+min+":"+sec, time.Local)
	}
	return time.ParseInLocation("2006-01-02", year+"-"+month+"-"+day, time.Local)
}
