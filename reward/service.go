package reward

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/exec"
	"sync"
)

type WebUI struct {
	path  string
	Conf  *Conf // 配置相关
	start StartFunc
}

type StartFunc func()

func (w WebUI) StartWebUI() {
	// 本地HTML文件的路径
	filePath := "./webui/webui.html"

	// 执行系统命令打开浏览器
	//cmd := exec.Command("xdg-open", filePath) // Linux 系统
	// cmd := exec.Command("open", filePath) // macOS 系统
	cmd := exec.Command("cmd", "/c", "start", filePath) // Windows 系统

	err := cmd.Run()
	if err != nil {
		// 处理错误
		panic(err)
	}
}

var conf *Conf

func (w *WebUI) ServiceWebUI(wg *sync.WaitGroup) {
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
	r.OPTIONS("/start", startHandler(w.start))

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
	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}

func (w *WebUI) SetStart(startFunc StartFunc) {
	w.start = startFunc
}
