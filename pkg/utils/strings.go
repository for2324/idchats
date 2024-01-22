/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/4/8 15:09).
 */
package utils

import (
	"Open_IM/pkg/common/constant"
	"encoding/json"
	"math"
	"math/big"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}
func Interface2JsonString(dt interface{}) string {
	byteString, _ := json.Marshal(dt)
	return string(byteString)
}
func StringToInt(i string) int {
	j, _ := strconv.Atoi(i)
	return j
}
func StringToInt64(i string) int64 {
	j, _ := strconv.ParseInt(i, 10, 64)
	return j
}
func StringToInt32(i string) int32 {
	j, _ := strconv.ParseInt(i, 10, 64)
	return int32(j)
}
func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func Uint32ToString(i uint32) string {
	return strconv.FormatInt(int64(i), 10)
}

func IsDigit(str string) bool {
	for _, x := range []rune(str) {
		if !unicode.IsDigit(x) {
			return false
		}
	}
	return true
}
func Float64ToString(v float64) string {
	return strconv.FormatFloat(v, 'f', 10, 64)
}

// judge a string whether in the  string list
func IsContain(target string, List []string) bool {
	for _, element := range List {

		if target == element {
			return true
		}
	}
	return false
}
func IsContainEqual(target string, List []string) bool {
	for _, element := range List {

		if strings.EqualFold(element, target) {
			return true
		}
	}
	return false
}
func IsContainInt32(target int32, List []int32) bool {
	for _, element := range List {
		if target == element {
			return true
		}
	}
	return false
}
func IsContainInt(target int, List []int) bool {
	for _, element := range List {
		if target == element {
			return true
		}
	}
	return false
}
func InterfaceArrayToStringArray(data []interface{}) (i []string) {
	for _, param := range data {
		i = append(i, param.(string))
	}
	return i
}
func StructToJsonString(param interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func StructToJsonBytes(param interface{}) []byte {
	dataType, _ := json.Marshal(param)
	return dataType
}

// The incoming parameter must be a pointer
func JsonStringToStruct(s string, args interface{}) error {
	err := json.Unmarshal([]byte(s), args)
	return err
}

func GetMsgID(sendID string) string {
	t := Int64ToString(GetCurrentTimestampByNano())
	return Md5(t + sendID + Int64ToString(rand.Int63n(GetCurrentTimestampByNano())))
}
func GetConversationIDBySessionType(sourceID string, sessionType int) string {
	switch sessionType {
	case constant.SingleChatType:
		return "single_" + sourceID
	case constant.GroupChatType:
		return "group_" + sourceID
	case constant.SuperGroupChatType:
		return "super_group_" + sourceID
	case constant.NotificationChatType:
		return "notification_" + sourceID
	}
	return ""

}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
func UInt64ToString(i uint64) string {
	return strconv.FormatUint(i, 10)
}
func RemoveDuplicateElement(idList []string) []string {
	result := make([]string, 0, len(idList))
	temp := map[string]struct{}{}
	for _, item := range idList {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
func String2bytes(str string) []byte {
	if str == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(str), len(str))
}
func Bytes2string(bs []byte) string {
	if len(bs) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}

// 从 big int 对象创建字节数组
func BigIntToBytes(i *big.Int) []byte {
	return i.Bytes()
}

// 从字节数组创建 big int 对象
func BytesToBigInt(b []byte) *big.Int {
	i := new(big.Int)
	i.SetBytes(b)
	return i
}
func IsFloat64Zero(num float64) bool {
	// 使用一个很小的精度范围来判断
	epsilon := 1e-9
	return math.Abs(num) < epsilon
}
func IsEmailValid(email string) bool {
	// 定义正则表达式的模式
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// 编译正则表达式
	regex := regexp.MustCompile(pattern)

	// 使用正则表达式进行匹配
	return regex.MatchString(email)
}
