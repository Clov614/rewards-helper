package reward

import (
	"github.com/levigross/grequests"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

type UrlGet string
type UaPc string
type UaMb string
type TypeUa string

type Get struct {
	Url  UrlGet
	Info Infog
	//ProxyUrl *url.URL
	//client   *http.Client
	UApc UaPc
	UAmb UaMb
	RO   grequests.RequestOptions
}

// 请求后返回的信息

type Infog struct {
	Pc string
	mb string
}

// 发起请求 _type string: UA头的类型 (pc mb)表示电脑或手机
func (g Get) do(c *Conn, _type string) *grequests.Response {
	// 判断是否开启代理
	if c.Conf.ProxyOn {
		if c.Conf.Proxy == "" {
			panic("Conf.Proxy Empty!!")
		}
		proxyURL, err := url.Parse(c.Conf.Proxy)
		if err != nil {
			panic(err)
		}
		g.RO = grequests.RequestOptions{
			Proxies: map[string]*url.URL{proxyURL.Scheme: proxyURL},
			Cookies: c.Cookie.Cookies,
		}

	} else {
		g.RO = grequests.RequestOptions{
			Cookies: c.Cookie.Cookies,
		}
	}
	// 组合搜索url
	if len(c.Conf.KeyWords) == 0 {
		log.Fatalln("c.Conf.KeyWords == 0", "请配置KeyWords或删除conf.yaml文件(重置配置)")
	}
	rand.Seed(time.Now().Unix()) // 设置随机数种子
	keyword := c.Conf.KeyWords[rand.Intn(len(c.Conf.KeyWords))] + strconv.Itoa(rand.Intn(10000))
	url := string(g.Url) + "?q=" + url.QueryEscape(keyword)

	//// new req
	//req, err := http.NewRequest("GET", url, nil)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	// 设置UA头
	//if _type == "pc" {
	//	req.Header.Set("User-Agent", string(g.UApc))
	//} else {
	//	req.Header.Set("User-Agent", string(g.UAmb))
	//}
	if _type == "mb" {
		g.RO.UserAgent = string(g.UAmb)
	}

	// 向req添加Cookies
	//for _, v := range c.Cookie.Cookies {
	//	req.AddCookie(v)
	//}
	// do
	//resp, err := g.client.Do(req)
	resp, err := grequests.Get(url, &g.RO)
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}

// _type 为请求的UA类型: pc or mb(mobilePhone)
func (g *Get) Handler(c *Conn, searchUrl UrlGet, UApc UaPc, UAmb UaMb, _type TypeUa) {
	g.Url = searchUrl
	g.UApc = UApc
	g.UAmb = UAmb
	resp := g.do(c, string(_type))
	defer resp.Close()
	if resp.StatusCode == 200 {
		log.Println("<"+_type+"> ", "200 OK")
		log.Println("当前分数:", c.View.Infov.AvailablePoints)
	} else {
		log.Println("bad response", "code: ", resp.StatusCode)
	}
}
