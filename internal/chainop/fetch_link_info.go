package chainop

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	tls "github.com/refraction-networking/utls"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"k8s.io/utils/strings/slices"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func checkIsEffectiveHostname(hostName string) bool {
	if hostName == "" {
		return false
	}
	// 如果主机名是数字 IP，则认为不满足条件
	isIP := net.ParseIP(hostName) != nil
	if isIP {
		fmt.Println("主机名不能是数字 IP")
		return false
	}

	// 正则表达式匹配非本地 IP 和 localhost
	isLocalIP := regexp.MustCompile(`^(127\.0\.0\.1|::1|0\.0\.0\.0)$`).MatchString(hostName)
	if isLocalIP || hostName == "localhost" {
		fmt.Println("主机名不能是本地 IP 或 localhost")
		return false
	}

	fmt.Println("主机名符合条件")
	return true
}

var confighostName = []string{"pro-api.coinmarketcap.com",
	"app.geckoterminal.com",
	"api.fgsasd.org",
	"api.ooesafph5.com",
	"www.dextools.io",
	"id.dexscreener.com",
	"api.dexview.com"}

func checkIsConfigHostName(str string) bool {
	return slices.Contains(confighostName, str)
}
func RedirectUrlLink(c *gin.Context) {
	urlValue := c.Query("url")
	if urlValue == "" {
		c.String(http.StatusOK, "")
		return
	}
	if !strings.HasPrefix(urlValue, "http://") && !strings.HasPrefix(urlValue, "https://") {
		c.String(http.StatusOK, "")
		return
	}
	urlHostName := utils.GetHostnameFromUrl(urlValue)
	fmt.Println("request ip:", c.ClientIP(), "request host name:", urlHostName)

	if !checkIsEffectiveHostname(urlHostName) {
		c.String(http.StatusOK, "")
		return
	}
	if !checkIsConfigHostName(urlHostName) {
		c.String(http.StatusOK, "")
		return
	}
	if _, err := url.Parse(urlValue); err != nil {
		c.String(http.StatusOK, "")
		return
	} else {
		// 获取请求头
		header := c.Request.Header
		headers := make(map[string]string)
		for k, v := range header {
			headers[k] = strings.Join(v, ", ")
		}
		var bodyByte []byte
		var err error
		var contentEncode = ""
		if urlHostName == "pro-api.coinmarketcap.com" {
			headers["X-Cmc_pro_api_key"] = "1ea4909a-388b-4a05-ba31-b1802bd7fcba"
		}
		if !config.Config.IsPublicEnv {
			if urlHostName == "api.dexview.com" {
				headers = map[string]string{
					"Secret": headers["Secret"],
				}
				bodyByte, resultCode, _ := utils.HttpGetWithHeaderWithGout(urlValue, headers, "http://proxy.idchats.com:7890")
				if err != nil {
					fmt.Println(err.Error())
				}
				c.String(resultCode, bodyByte)
				return
			} else {
				proxyAddress, _ := url.Parse("http://proxy.idchats.com:7890")
				bodyByte, err = utils.HttpGetWithHeaderWithProxy(urlValue, headers, http.ProxyURL(proxyAddress))
			}

		} else {
			if urlHostName == "app.geckoterminal.com" {
				// 解析 URL
				u, _ := url.Parse(urlValue)
				baseURL := u.Scheme + "://" + u.Host + u.Path
				queryParams := u.Query()
				paramsMap := make(map[string]string)
				for key, values := range queryParams {
					if len(values) > 0 {
						paramsMap[key] = values[0]
					}
				}
				postDataToChinaServer := &TGeckoterminal{
					Url:     baseURL,
					Methods: "GET",
					Code:    "UTF-8",
					Header:  map[string]string{"content-type": "application/x-www-form-urlencoded"},
					Parm:    paramsMap,
					Cookie:  "",
					Proxy:   "",
				}
				jsonByteString, _ := json.Marshal(postDataToChinaServer)
				bodyByte, err = utils.HttpPost(
					"https://api.fly63.com/home/static/php/http/api.php", "",
					map[string]string{"Content-Type": "application/json"},
					jsonByteString)
				type gecoProxyResultdata struct {
					Code int    `json:"code"`
					Data string `json:"data"`
				}
				var gecoResultdata gecoProxyResultdata
				json.Unmarshal(bodyByte, &gecoResultdata)
				fmt.Println(gecoResultdata.Data)
				c.Writer.Write([]byte(gecoResultdata.Data))
				return
			} else if urlHostName == "api.dexview.com" {
				headers = map[string]string{
					"Secret": "5ff3a258-2700-11ed-a261-0242ac120002",
				}
				bodyByte, err = utils.HttpGetWithHeader(urlValue, headers)
				if len(bodyByte) > 0 {
					fmt.Println(string(bodyByte))
				}
			} else {
				fmt.Printf("%#+v", headers)
				bodyByte, err = utils.HttpGetWithHeader(urlValue, headers)
				if len(bodyByte) > 0 {
					fmt.Println(string(bodyByte))
				}
			}
		}
		if err == nil {
			if contentEncode != "" {
				c.Header("Content-Encoding", contentEncode)
				c.Writer.Write(bodyByte)
			} else {
				c.String(http.StatusOK, utils.Bytes2string(bodyByte))
			}
			return
		} else {
			fmt.Println(err.Error())
			c.String(http.StatusOK, "")
			return
		}
	}
}

// {"url":"https://app.geckoterminal.com/api/p1/search","methods":"GET","code":"UTF-8","header":{"content-type":"application/x-www-form-urlencoded"},"parm":{"query":"wbnb"},"cookie":"","proxy":""}
type TGeckoterminal struct {
	Url     string            `json:"url"`
	Methods string            `json:"methods"`
	Code    string            `json:"code"`
	Header  map[string]string `json:"header"`
	Parm    map[string]string `json:"parm"`
	Cookie  string            `json:"cookie"`
	Proxy   string            `json:"proxy"`
}

func RedirectUrlLinkV2(c *gin.Context) {
	urlValue := c.Query("url")
	if urlValue == "" {
		c.String(http.StatusOK, "")
		return
	}
	if !strings.HasPrefix(urlValue, "http://") && !strings.HasPrefix(urlValue, "https://") {
		c.String(http.StatusOK, "")
		return
	}
	urlHostName := utils.GetHostnameFromUrl(urlValue)
	fmt.Println("request ip:", c.ClientIP(), "request host name:", urlHostName)
	if !checkIsEffectiveHostname(urlHostName) {
		c.String(http.StatusOK, "")
		return
	}
	if !checkIsConfigHostName(urlHostName) {
		c.String(http.StatusOK, "")
		return
	}
	if _, err := url.Parse(urlValue); err != nil {
		c.String(http.StatusOK, "")
		return
	} else {
		// 获取请求头
		header := c.Request.Header
		headers := make(map[string]string)
		for k, v := range header {
			headers[k] = strings.Join(v, ", ")
		}
		var bodyByte []byte
		var err error
		if urlHostName == "pro-api.coinmarketcap.com" {
			headers["X-Cmc_pro_api_key"] = "1ea4909a-388b-4a05-ba31-b1802bd7fcba"
		}
		newTransport := &http.Transport{
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {

				dialer := net.Dialer{}
				conn, err := dialer.DialContext(ctx, network, addr)
				if err != nil {
					fmt.Println("dialer.DialContext")
					return nil, err
				}
				//conn, err := netx.DialContextWithoutProxy(ctx, network, addr)
				//if err != nil {
				//	println("Error creating connection #123", err)
				//	return nil, err
				//}
				host, _, err := net.SplitHostPort(addr)
				if err != nil {
					fmt.Println("net.SplitHostPort")
					return nil, err
				}
				config := &tls.Config{ServerName: host}
				uconn := tls.UClient(conn, config, tls.HelloCustom)
				if err := uconn.ApplyPreset(NewSpec()); err != nil {
					return nil, err
				}
				if err := uconn.Handshake(); err != nil {
					return nil, err
				}
				return uconn, nil
			},
		}
		if !config.Config.IsPublicEnv {
			//	proxyAddress, _ := url.Parse("http://proxy.idchats.com:7890")
			//newTransport.Proxy = http.ProxyURL(proxyAddress)
		}
		req, err := http.NewRequest("GET", urlValue, nil)
		if err != nil {
			panic(err)
		}
		req.Header = http.Header{
			"Referer":         {"https://" + urlHostName + "/"},
			"Accept":          {"application/json"},
			"Accept-Language": {"en-US,en;q=0.5"},
			"Cache-Control":   {"no-cache"},
			"Pragma":          {"no-cache"},
			"User-Agent":      {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"},
		}
		res, err := newTransport.RoundTrip(req)
		if err != nil {
			c.String(http.StatusForbidden, "")
			return
		}
		bodyByte = bodyByte
		defer res.Body.Close()
		resp, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		c.Status(res.StatusCode)
		c.Header("encode-type", res.Header.Get("encode-type"))
		c.Writer.Write(resp)
	}
}

type LinkInfo struct {
	Title   string            `json:"title"`
	Favicon string            `json:"favicon"`
	OG      map[string]string `json:"og"`
}

type FetchParams struct {
	Url string `json:"url" binding:"required"`
}

func FetchLink(c *gin.Context) {
	fetchParams := FetchParams{}
	err := c.Bind(&fetchParams)
	if err != nil {
		c.JSON(200, gin.H{
			"data": gin.H{},
			"msg":  err.Error(),
			"code": 500,
		})
		logrus.Error(err)
		return
	}
	info, fetchErr := fetchLinkInfo(fetchParams.Url)
	if fetchErr != nil {
		c.JSON(200, gin.H{
			"data": gin.H{},
			"msg":  fetchErr.Error(),
			"code": 500,
		})
		logrus.Error(err)
		return
	}
	c.JSON(200, gin.H{
		"data": info,
		"msg":  "",
		"code": 200,
	})
}

// get a web html content og:image og:title: og:descript
func databbbfun() {

}

func fetchLinkInfo(urlstr string) (LinkInfo, error) {
	var linkInfo LinkInfo

	client := http.Client{}
	if config.Config.OpenNetProxy.OpenFlag {
		proxyAddress, _ := url.Parse("http://proxy.idchats.com:7890")
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyAddress),
		}
	}
	req, err := http.NewRequest("GET", urlstr, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)")
	req.Header.Set("Connection", "close")
	req.Header.Set("Range", "bytes=0-524288")
	// 发起 GET 请求
	resp, err := client.Do(req)
	if err != nil {
		return linkInfo, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return linkInfo, err
	}

	// 解析 HTML 文档
	doc, err := html.Parse(strings.NewReader(string(bodyBytes)))
	if err != nil {
		return linkInfo, err
	}

	// 递归查找 <title> 标签
	var findTitle func(*html.Node)
	findTitle = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil && len(n.FirstChild.Data) > 0 {
			linkInfo.Title = n.FirstChild.Data
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findTitle(c)
		}
	}
	findTitle(doc)

	// 查找 <link rel="shortcut icon"> 标签
	var findFavicon func(*html.Node)
	findFavicon = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "link" {
			rel := ""
			href := ""
			for _, attr := range n.Attr {
				if attr.Key == "rel" {
					rel = strings.ToLower(attr.Val)
				}
				if attr.Key == "href" {
					href = attr.Val
				}
			}
			index := strings.Index(rel, "icon")
			if index > -1 && len(href) > 0 {
				if strings.HasPrefix(href, "//") {
					linkInfo.Favicon = fmt.Sprintf("https:%s", href)
				} else if !strings.HasPrefix(href, "http") {
					linkInfo.Favicon = fmt.Sprintf("%s://%s%s", resp.Request.URL.Scheme, resp.Request.URL.Host, href)
				} else {
					linkInfo.Favicon = href
				}
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findFavicon(c)
		}
	}
	findFavicon(doc)

	// 查找 <meta property="og:xxx" content="yyy"> 标签
	ogTags := map[string]string{}
	var findOG func(*html.Node)
	findOG = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			property := ""
			content := ""
			for _, attr := range n.Attr {
				if attr.Key == "property" {
					property = strings.ToLower(attr.Val)
				}
				if attr.Key == "content" {
					content = attr.Val
				}
			}
			if len(property) > 0 && len(content) > 0 {
				if strings.HasPrefix(property, "og:") {
					ogTags[property] = content
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findOG(c)
		}
	}
	findOG(doc)
	if len(ogTags) > 0 {
		linkInfo.OG = ogTags
	}

	return linkInfo, nil
}
