package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func VerifyEmailFormat(email string) bool {
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// copy a by b  b->a
func CopyStructFields(a interface{}, b interface{}, fields ...string) (err error) {
	return copier.Copy(a, b)
}

// 有效的地址url
func IsVaildUrl(urlstr string) bool {
	url := urlstr
	if !strings.HasPrefix(urlstr, "https://") && strings.HasPrefix(urlstr, "http://") {
		url = "https://" + urlstr
	}
	pattern := `^(http|https)://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}(/[a-zA-Z0-9_/\-.]*)?$`
	match, _ := regexp.MatchString(pattern, url)
	if match {
		return true
	} else {
		return false
	}
}

// 有效的邮箱地址
func IsVaildEmail(emailAddress string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, emailAddress)
	if match {
		return true
	} else {
		return false
	}
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, "==> "+printCallerNameAndLine()+message)
}

func WithMessage(err error, message string) error {
	return errors.WithMessage(err, "==> "+printCallerNameAndLine()+message)
}

func printCallerNameAndLine() string {
	pc, _, line, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name() + "()@" + strconv.Itoa(line) + ": "
}

func GetSelfFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return cleanUpFuncName(runtime.FuncForPC(pc).Name())
}
func cleanUpFuncName(funcName string) string {
	end := strings.LastIndex(funcName, ".")
	if end == -1 {
		return ""
	}
	return funcName[end+1:]
}
func GetRequestName(urlHost string) string {
	ip := net.ParseIP(urlHost)
	if ip != nil {
		return urlHost
	} else {
		if index := strings.LastIndex(urlHost, ":"); index != -1 {
			return urlHost[:index]
		}
		return urlHost
	}
}

// Get the intersection of two slices
func Intersect(slice1, slice2 []uint32) []uint32 {
	m := make(map[uint32]bool)
	n := make([]uint32, 0)
	for _, v := range slice1 {
		m[v] = true
	}
	for _, v := range slice2 {
		flag, _ := m[v]
		if flag {
			n = append(n, v)
		}
	}
	return n
}

// Get the diff of two slices
func Difference(slice1, slice2 []uint32) []uint32 {
	m := make(map[uint32]bool)
	n := make([]uint32, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v] = true
	}
	for _, v := range slice1 {
		if !m[v] {
			n = append(n, v)
		}
	}

	for _, v := range slice2 {
		if !m[v] {
			n = append(n, v)
		}
	}
	return n
}

// Get the intersection of two slices
func IntersectString(slice1, slice2 []string) []string {
	m := make(map[string]bool)
	n := make([]string, 0)
	for _, v := range slice1 {
		m[v] = true
	}
	for _, v := range slice2 {
		flag, _ := m[v]
		if flag {
			n = append(n, v)
		}
	}
	return n
}

// Get the diff of two slices
func DifferenceString(slice1, slice2 []string) []string {
	m := make(map[string]bool)
	n := make([]string, 0)
	inter := IntersectString(slice1, slice2)
	for _, v := range inter {
		m[v] = true
	}
	for _, v := range slice1 {
		if !m[v] {
			n = append(n, v)
		}
	}

	for _, v := range slice2 {
		if !m[v] {
			n = append(n, v)
		}
	}
	return n
}
func OperationIDGenerator() string {
	return strconv.FormatInt(time.Now().UnixNano()+int64(rand.Uint32()), 10)
}

func RemoveRepeatedStringInList(slc []string) []string {
	var result []string
	tempMap := map[string]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

func Pb2String(pb proto.Message) (string, error) {
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
	}
	return marshaler.MarshalToString(pb)
}

func String2Pb(s string, pb proto.Message) error {
	return proto.Unmarshal([]byte(s), pb)
}

func Map2Pb(m map[string]string) (pb proto.Message, err error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(b, pb)
	if err != nil {
		return nil, err
	}
	return pb, nil
}
func Pb2Map(pb proto.Message) (map[string]interface{}, error) {
	_buffer := bytes.Buffer{}
	jsonbMarshaller := &jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  true,
		EmitDefaults: false,
	}
	_ = jsonbMarshaller.Marshal(&_buffer, pb)
	jsonCnt := _buffer.Bytes()
	var out map[string]interface{}
	err := json.Unmarshal(jsonCnt, &out)
	return out, err
}
func IndexContain(target string, List []string) int {

	for key := 0; key < len(List); key++ {

		if target == List[key] {
			return key
		}
	}
	return -1
}
func FloatCompare(f1, f2 interface{}) (n int, err error) {
	var f1Dec, f2Dec decimal.Decimal
	switch f1.(type) {
	case float64:
		f1Dec = decimal.NewFromFloat(f1.(float64))
		switch f2.(type) {
		case float64:
			f2Dec = decimal.NewFromFloat(f2.(float64))
		case string:
			f2Dec, err = decimal.NewFromString(f2.(string))
			if err != nil {
				return 2, err
			}
		default:
			return 2, errors.New("FloatCompare() expecting to receive float64 or string")
		}
	case string:
		f1Dec, err = decimal.NewFromString(f1.(string))
		if err != nil {
			return 2, err
		}
		switch f2.(type) {
		case float64:
			f2Dec = decimal.NewFromFloat(f2.(float64))
		case string:
			f2Dec, err = decimal.NewFromString(f2.(string))
			if err != nil {
				return 2, err
			}
		default:
			return 2, errors.New("FloatCompare() expecting to receive float64 or string")
		}
	default:
		return 2, errors.New("FloatCompare() expecting to receive float64 or string")
	}
	return f1Dec.Cmp(f2Dec), nil
}
func GetHostnameFromUrl(urlString string) string {
	parsedUrl, err := url.Parse(strings.ToLower(urlString))
	if err != nil {
		return ""
	}
	fmt.Println("request host name:", parsedUrl.Host)
	return parsedUrl.Host
}
