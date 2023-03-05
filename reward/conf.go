package reward

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Conf struct {
	ProxyOn  bool     `yaml:"proxy_on"`
	Proxy    string   `yaml:"proxy"`
	KeyWords []string `yaml:"key_words"`
}

var (
	InfoHelp = "# 配置文件\n# proxy_on 是否开启代理模式 （true or false)\n# proxy 代理地址 \n# key_words 发起请求的关键词\n"
	// 默认配置
	DefaultConf = Conf{
		ProxyOn:  false,
		Proxy:    "http://127.0.0.1:7890",
		KeyWords: []string{"a", "AB", "bingNew", "iaimi", "爱美", "孤独摇滚", "b站", "bilibili", "大佐de苦de手", "柚子社"},
	}
)

// 判断conf是否存在
func (c *Conf) IsExist() bool {
	if !PathExists("./conf/conf.yaml") {
		return false
	}
	return true
}

// 初始化conf
func (c *Conf) InitConfDefault() {
	log.Println("生成默认conf文件中")
	log.Println("[默认不开启代理]")
	err := os.MkdirAll("./conf", 0766)
	if err != nil {
		log.Fatalln("创建conf目录失败:", err)
	}
	file, err := os.OpenFile("./conf/conf.yaml", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Fatalln("创建conf.yaml失败:", err)
	}
	defer file.Close()
	file.Write([]byte(InfoHelp))
	// 初始化conf.yaml
	dataStr, _ := yaml.Marshal(&DefaultConf)
	file.Write(dataStr)
}

// 加载配置
func (c *Conf) ReadConf() error {
	err := ReadYaml(&c, "./conf/conf.yaml")
	if err != nil {
		return err
	}
	return nil
}

// 判断路径是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 读取yaml
func ReadYaml(_type interface{}, path string) (err error) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("读取Error path: "+path, err)
	}
	err = yaml.Unmarshal(file, _type)
	if err != nil {
		return err
	}
	return nil
}

// 处理conf
func (c *Conf) Handler() {
	exist := c.IsExist()
	if !exist {
		c.InitConfDefault()
	}
	err := c.ReadConf()
	if err != nil {
		log.Fatalln(err)
	}
}
