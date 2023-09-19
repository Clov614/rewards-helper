package reward

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"rewards-helper/reward/html"
	"runtime"
	"sync"
)

type WebUI struct {
	path     string
	Conf     *Conf // 配置相关
	start    StartFunc
	ViewInfo *string // 分数输出
}

type StartFunc func()

var conn *Conn
var conf *Conf

const port = 8080

func (w WebUI) StartWebUI(c *Conn, wg *sync.WaitGroup) {
	// 载入全局连接
	conn = c
	//// 本地HTML文件的路径
	//filePath := "./webui/webui.html"
	//
	//// 执行系统命令打开浏览器
	////cmd := exec.Command("xdg-open", filePath) // Linux 系统
	//// cmd := exec.Command("open", filePath) // macOS 系统
	//cmd := exec.Command("cmd", "/c", "start", filePath) // Windows 系统
	//
	//err := cmd.Run()
	//if err != nil {
	//	// 处理错误
	//	panic(err)
	//}
	go w.StartWebPage(wg)

	// 打开默认浏览器
	err := openURL(fmt.Sprintf("http://localhost:%d/", port))
	if err != nil {
		log.Fatalln("openURL error:", err)
	}
}

func (w *WebUI) ServiceWebUI(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	// Initialize Gin router
	r := gin.Default()

	// 使用CORS中间件处理跨域请求
	r.Use(cors.Default())
	conf = w.Conf
	// Define an endpoint to handle settings retrieval
	r.GET("/settings", getConfHandler)

	// Define an endpoint to handle settings update
	r.POST("/settings", updateConfHandler)

	// 启动脚本
	r.POST("/start", startHandler(w.start))

	// Run the server on port 8099
	r.Run(":8099")

}

func startHandler(startFunc StartFunc) gin.HandlerFunc {

	return func(context *gin.Context) {
		startFunc() // 启动服务
	}
}

func getConfHandler(c *gin.Context) {
	// Return the current settings as JSON response
	c.JSON(http.StatusOK, conf)
}

func updateConfHandler(c *gin.Context) {
	var newConf Conf

	// Bind the request body to the newSettings variable
	if err := c.BindJSON(&newConf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update the current settings with the newSettings
	*conf = newConf
	err := conf.WriteConf()
	if err != nil {
		log.Fatalln("writeConf error: ", err)
	}
	// 更新cookies.txt文件
	conn.Cookie.UpdateCookies(conf.Cookies)

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}

func (w *WebUI) SetStart(startFunc StartFunc) {
	w.start = startFunc
}

// 更新websorcket
var upgrader = websocket.Upgrader{}

type TemplateData struct {
	Host string
	Data []string
}

func (wu *WebUI) StartWebPage(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	// 设置HTTP处理函数
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 读取HTML文件内容
		tmpl, err := template.New("index").Parse(html.HtmlTemplate)
		if err != nil {
			http.Error(w, "Error processing template", http.StatusInternalServerError)
			return
		}
		// 渲染模板并将结果写入ResponseWriter
		err = tmpl.Execute(w, r.Host)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}

	})

	http.HandleFunc("/getInfo", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Error upgrading connection:", err)
			return
		}
		defer conn.Close()

		for {
			// 将数据发送给客户端
			err := conn.WriteMessage(websocket.TextMessage, []byte(*wu.ViewInfo))
			if err != nil {
				fmt.Println("Error sending data:", err)
				return
			}
		}
	})

	// 启动HTTP服务器并监听指定端口
	port := 8080
	fmt.Printf("Server running on port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Error starting the server: %s\n", err)
	}
}

func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	default:
		return fmt.Errorf("unsupported platform")
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
