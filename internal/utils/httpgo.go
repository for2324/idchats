package utils

import (
	"Open_IM/pkg/common/config"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/guonaihong/gout"
	"io/ioutil"
	"net/http"
)

func PostUrl(url string, data interface{}, headers map[string]string) (string, error) {
	bodyByte, err := json.Marshal(data)
	if err != nil {
		return "", nil
	}
	//fmt.Println("ThirdPart Post:", string(bodyByte))
	reader := bytes.NewReader(bodyByte)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return "", nil
	}

	request.Header.Set("Content-type", "application/json;charset=UTF-8")
	//request.Header.Set("Content-type", "application/json")
	if headers != nil && len(headers) != 0 {
		for key := range headers {
			request.Header.Set(key, headers[key])
		}
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", nil
	}
	defer response.Body.Close()
	return string(result), nil
}
func GetUrl(domain string, query, headers map[string]string) (resultString string, err error) {
	dataForma := gout.GET(domain).SetHeader(headers)
	if !config.Config.IsPublicEnv {
		dataForma = dataForma.SetProxy("http://proxy.idchats.com:7890")
	}
	dataForma = dataForma.Debug(!config.Config.IsPublicEnv)
	var code int
	err = dataForma.BindBody(&resultString).Code(&code).Do()
	return
}

func ObjectToJson(src interface{}) (string, error) {
	if result, err := json.Marshal(src); err != nil {
		return "", errors.New("Json Str Parse Err: " + err.Error())
	} else {
		return string(result), nil
	}
}

func JsonToObject(src string, target interface{}) error {
	if err := json.Unmarshal([]byte(src), target); err != nil {
		return errors.New("Json Str Parse Err: " + err.Error())
	}
	return nil
}

func JsonToAny(src interface{}, target interface{}) error {
	if src == nil || target == nil {
		return errors.New("Param is empty")
	}
	str, err := ObjectToJson(src)
	if err != nil {
		return err
	}
	if err := JsonToObject(str, target); err != nil {
		return err
	}
	return nil
}
