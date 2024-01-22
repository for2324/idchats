package utils

import (
	"fmt"
	"testing"
)

func TestHttpGetWithHeader(t *testing.T) {
	got, _ := HttpGetWithHeader("https://api.biubiu.id/graph/redirecturl?url=https://app.geckoterminal.com/api/p1/eth/pools/0x9476f8fccefcee23481f9fa5bb5d8bdf1f145f5c?base_token=0", map[string]string{
		"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36",
	})
	//fmt.Println(err.error)
	fmt.Println(len(got))
	if len(got) > 0 {
		fmt.Println(string(got))
	}

}
