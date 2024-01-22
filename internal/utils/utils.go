package utils

import (
	"encoding/base64"
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"reflect"
	"strconv"
	"time"
)

func JsonDataList(resp interface{}) []map[string]interface{} {
	var list []proto.Message
	if reflect.TypeOf(resp).Kind() == reflect.Slice {
		s := reflect.ValueOf(resp)
		for i := 0; i < s.Len(); i++ {
			ele := s.Index(i)
			list = append(list, ele.Interface().(proto.Message))
		}
	}

	result := make([]map[string]interface{}, 0)
	for _, v := range list {
		m := ProtoToMap(v, false)
		result = append(result, m)
	}
	return result
}

// JsonDataListV2 不忽略零值
func JsonDataListV2(resp interface{}) []map[string]interface{} {
	var list []proto.Message
	if reflect.TypeOf(resp).Kind() == reflect.Slice {
		s := reflect.ValueOf(resp)
		for i := 0; i < s.Len(); i++ {
			ele := s.Index(i)
			list = append(list, ele.Interface().(proto.Message))
		}
	}

	result := make([]map[string]interface{}, 0)
	for _, v := range list {
		m := ProtoToMapV2(v, false)
		result = append(result, m)
	}
	return result
}

func JsonDataOne(pb proto.Message) map[string]interface{} {
	return ProtoToMap(pb, false)
}

func ProtoToMap(pb proto.Message, idFix bool) map[string]interface{} {
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
	}

	s, _ := marshaler.MarshalToString(pb)
	out := make(map[string]interface{})
	json.Unmarshal([]byte(s), &out)
	if idFix {
		if _, ok := out["id"]; ok {
			out["_id"] = out["id"]
			delete(out, "id")
		}
	}
	return out
}

func ProtoToMapV2(pb proto.Message, idFix bool) map[string]interface{} {
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: true,
	}

	s, _ := marshaler.MarshalToString(pb)
	out := make(map[string]interface{})
	json.Unmarshal([]byte(s), &out)
	if idFix {
		if _, ok := out["id"]; ok {
			out["_id"] = out["id"]
			delete(out, "id")
		}
	}
	return out
}
func GetOkLinkXAPIKey() string {
	t := time.Now().UnixNano() / int64(time.Millisecond)
	t1 := strconv.Itoa(int(1*t + 1111111111111))
	r := strconv.Itoa(rand.Intn(10))
	n := strconv.Itoa(rand.Intn(10))
	o := strconv.Itoa(rand.Intn(10))
	newT := t1 + r + n + o

	apiKey := "a2c903cc-b31e-4547-9299-b6d07b7631ab"
	key1 := apiKey[:8]
	key2 := apiKey[8:]
	newKey := key2 + key1

	res := newKey + "|" + newT
	// 对数据进行 base64 编码
	xAPIKey := base64.StdEncoding.EncodeToString([]byte(res))
	return xAPIKey
}
