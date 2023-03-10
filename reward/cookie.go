package reward

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Cookie struct {
	Cookies   []*http.Cookie
	cookieStr string
	cookieL   []string
}

// 判断cookies.txt是否存在
func (c *Cookie) isExist() bool {
	if !PathExists("./cookie/cookies.txt") {
		return false
	}
	return true
}

func (c *Cookie) initPath() {
	err := os.MkdirAll("./cookie", 0766)
	if err != nil {
		log.Fatalln("创建cookie目录失败:", err)
	}
	file, err := os.OpenFile("./cookie/cookies.txt", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Fatalln("创建cookie.txt失败:", err)
	}
	defer file.Close()
	log.Println("请将cookies添加至/cookie/cookies.txt")
	time.Sleep(time.Second * 5)
	os.Exit(114514)
}

func (c *Cookie) txt2Cookies() {

	for _, v := range c.cookieL {
		tmpL := ownSplit(v, "=")
		c.Cookies = append(c.Cookies, &(http.Cookie{Name: url.QueryEscape(tmpL[0]), Value: url.QueryEscape(tmpL[1])}))
	}
}

// cookie切片规则
func ownSplit(preStr string, pattern string) (preL []string) {
	firstend := -1
	for i, v := range preStr {
		if string(v) == pattern {
			firstend = i
			break
		}
	}
	preL = make([]string, 10)
	preL[0] = preStr[0:firstend]
	preL[1] = preStr[firstend+1:]
	return
}

func (c *Cookie) isEmpty() bool {
	// Cookies.txt to Cookies
	f, _ := os.OpenFile("./cookie/cookies.txt", os.O_RDWR, 0777)
	cookieB, _ := ioutil.ReadAll(f)
	cookieS := strings.TrimSpace(string(cookieB))
	c.cookieStr = cookieS
	cookieL := make([]string, 10)
	cookieL = strings.Split(cookieS, "; ")
	c.cookieL = cookieL
	// cookie是否没配置
	if c.cookieStr == "" || len(c.cookieL) == 0 {
		return true
	}
	return false
}

func (c *Cookie) Handler() {
	if !c.isExist() {
		c.initPath()
	}
	if c.isEmpty() {
		log.Println("cookie为空,请配置cookie")
		time.Sleep(time.Second * 5)
		os.Exit(400)
	}
	c.txt2Cookies()
}
