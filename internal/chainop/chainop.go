package chainop

import (
	"Open_IM/pkg/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type TResultData struct {
	Data struct {
		Domains []struct {
			Name     string `json:"name"`
			Typename string `json:"__typename"`
			Id       string `json:"id"`
		} `json:"domains"`
	} `json:"data"`
}

type EnsAddressReq struct {
	Address string
}

var headers = map[string]string{
	"Content-Type":  "application/json",
	"Authorization": "",
}

func GraphSqlEnsDomainTest(c *gin.Context) {
	chainid := "5"
	ip := c.Request.Header.Get("X-Forward-For")
	if ip == "" {
		ip = c.ClientIP()
	}
	log.Println("request :", ip, chainid, "0x1c8ec996420db47c0859a5aaed0148fd23426bbf")
	ensUrl := "https://api.thegraph.com/subgraphs/name/ensdomains/ens"
	switch chainid {
	case "5":
		ensUrl = "https://api.thegraph.com/subgraphs/name/ensdomains/ensgoerli"

	}
	postData := `{"operationName":"getNamesFromSubgraph","variables":{"address":"0x1c8ec996420db47c0859a5aaed0148fd23426bbf"},"query":"query getNamesFromSubgraph($address: String!) {\n  domains(first: 1000, where: {resolvedAddress: $address}) {\n    name\n    __typename\n    id\n  }\n}\n"}`
	resbody, eerr := utils.HttpPost(ensUrl, "", map[string]string{"Content-Type": "application/json"}, []byte(postData))
	var dataGraph TResultData
	json.Unmarshal([]byte(resbody), &dataGraph)
	if eerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 450, "errMsg": "do not scan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": dataGraph})
	return
}
func GraphSqlEnsDomain(c *gin.Context) {
	var ask EnsAddressReq
	if err := c.ShouldBind(&ask); err != nil || ask.Address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "parame id error"})
		return
	}
	chainid := c.Request.Header.Get("chainId")
	ip := c.Request.Header.Get("X-Forward-For")
	if ip == "" {
		ip = c.ClientIP()
	}
	log.Println("request :", ip, chainid, ask.Address)
	ensUrl := "https://api.thegraph.com/subgraphs/name/ensdomains/ens"
	switch chainid {
	case "5":
		ensUrl = "https://api.thegraph.com/subgraphs/name/ensdomains/ensgoerli"

	}

	postData := `{"operationName":"getNamesFromSubgraph","variables":{"address":"` + strings.ToLower(ask.Address) + `"},"query":"query getNamesFromSubgraph($address: String!) {\n  domains(first: 1000, where: {resolvedAddress: $address}) {\n    name\n    __typename\n    id\n  }\n}\n"}`
	resbody, eerr := utils.HttpPost(ensUrl, "", map[string]string{"Content-Type": "application/json"}, []byte(postData))
	var dataGraph TResultData
	json.Unmarshal(resbody, &dataGraph)
	if eerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 450, "errMsg": "do not scan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": dataGraph})
	return
}

type EnsDomainName struct {
	EnsName string
	Address string
}

type TDomain struct {
	Data struct {
		Domains []struct {
			Owner struct {
				Id string `json:"id"`
			} `json:"owner"`
			Name     string `json:"name"`
			Resolver struct {
				Texts     interface{} `json:"texts"`
				CoinTypes []string    `json:"coinTypes"`
			} `json:"resolver"`
		} `json:"domains"`
	} `json:"data"`
}

func GraphSqlEnsDomainByName(c *gin.Context) {
	var ask EnsDomainName
	if err := c.ShouldBind(&ask); err != nil || ask.EnsName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "parame id error"})
		return
	}
	chainid := c.Request.Header.Get("chainId")
	ip := c.Request.Header.Get("X-Forward-For")
	if ip == "" {
		ip = c.ClientIP()
	}
	log.Println("request :", ip, chainid, ask.EnsName, ask.Address)
	ensUrl := "https://api.thegraph.com/subgraphs/name/ensdomains/ens"
	switch chainid {
	case "5":
		ensUrl = "https://api.thegraph.com/subgraphs/name/ensdomains/ensgoerli"

	}
	postData := `{"query":"{\n  domains(where: {name: \"` + ask.EnsName + `\"}) {\n    owner {\n      id\n    }\n    name\n    resolver {\n      texts\n      coinTypes\n    }\n  }\n}","variables":null,"extensions":{"headers":null}}`
	resbody, eerr := utils.HttpPost(ensUrl, "", map[string]string{"Content-Type": "application/json"}, []byte(postData))
	var dataGraph TDomain
	json.Unmarshal([]byte(resbody), &dataGraph)
	if eerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 450, "errMsg": "do not scan"})
		return
	}
	if len(dataGraph.Data.Domains) >= 1 && dataGraph.Data.Domains[0].Owner.Id == ask.Address {
		c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "ok"})
		return
	}
	log.Println("查询域名：", " 不匹配内容地址")

	c.JSON(http.StatusOK, gin.H{"errCode": 1, "errMsg": "域名和拥有者不匹配"})
	return
}
