package cmd

import (
	"log"
	"os"
	"time"
)

// 初始化.bat文件
func InitrunBat() {
	if !IsExistBat() {
		log.Println("生成run.bat")
		log.Println("生成run-web.bat")
		log.Println("强烈建议使用run.bat运行本程序")
		geneNewBat()
		time.Sleep(8 * time.Second)
		os.Exit(-2)
	}
}

func geneNewBat() {
	frun, err := os.OpenFile("./run.bat", os.O_CREATE|os.O_RDWR, 0777)
	fweb, err := os.OpenFile("./run-web.bat", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Fatalln("GeneDefaultBat ERROR", err)
	}
	defer frun.Close()
	defer fweb.Close()
	cmdRun := "start cmd /K rewards-helper.exe"
	cmdWeb := "start cmd /K rewards-helper.exe webui"
	frun.Write([]byte(cmdRun))
	fweb.Write([]byte(cmdWeb))
}

// 判断run.bat是否存在
func IsExistBat() bool {
	if !PathExists("./run.bat") {
		return false
	}
	return true
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
