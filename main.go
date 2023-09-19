package main

import (
	"fmt"
	"log"
	"os"
	"rewards-helper/reward"
	"rewards-helper/reward/cmd"
	"sync"
	"time"
)

var (
	ui = &reward.WebUI{}
)

func main() {
	cmd.InitrunBat()
	var wg sync.WaitGroup

	ViewUrl := "https://rewards.bing.com/"
	conn := reward.New(ViewUrl, ui)
	if len(os.Args) <= 1 {
		start(conn)
	} else if os.Args[1] == "webui" {
		ui.StartWebUI(conn, &wg)
		// 传递启动函数
		ui.SetStart(func() {
			start(conn)
		})
		// TODO web服务
		wg.Add(1)
		go ui.ServiceWebUI(&wg)
		wg.Wait()
	}

	//// 测试使用
	//ui.StartWebUI(conn, &wg)
	//// 传递启动函数
	//ui.SetStart(func() {
	//	start(conn)
	//})
	//go ui.ServiceWebUI(&wg)
	//wg.Wait()
}

func start(conn *reward.Conn) {
	err := conn.View.Handler(conn)
	if err != nil {
		log.Fatalln(err)
	}
	if conn.Conf.ProxyOn {
		fmt.Println("[Info]当前处于代理模式!!!")
	}
	fmt.Println("[Info]开始获取积分")
	fmt.Println("当前可用分数: ", conn.View.Infov.AvailablePoints)
	fmt.Println("今日可获取最大分数: ", conn.View.Infov.DailyPoints.PointProgressMax)
	fmt.Println("今日分数: ", conn.View.Infov.DailyPoints.PointProgress)

	// 初始化任务管理器
	manager := conn.NewManager()
	params := reward.Params{
		Conn:   conn,
		UrlGet: reward.UrlGet(conn.Conf.SearchUrl),
		//UaPc:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edg/110.0.1587.63",
		UaPc: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
		UaMb: "Mozilla/5.0 (Linux; Android 11; PEAT00) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36 EdgA/110.0.1587.54",
	}
	// init任务管理器处理
	manager.Handler(params)
	// goroutine
	go manager.AddTask(conn.Get.Handler)
	go manager.StartTask()
	func() {
		statusPc, statusMb := 0, 0
		for i := range manager.DoneIndex {
			fmt.Println("Task: ", i)
			conn.View.Handler(conn)
			// 拼接显示信息
			dp := conn.View.Infov.DailyPoints
			viewinfo := fmt.Sprintf("%d/%d", dp.PointProgress, dp.PointProgressMax)
			ui.ViewInfo = &viewinfo
			// 命令行相关输出
			pcSearch := conn.View.Infov.PcSearch
			mobiSearch := conn.View.Infov.MobiSearch
			if statusPc == 0 && pcSearch.PointProgress == pcSearch.PointMax {
				log.Println("Pc分数刷取完毕")
				statusPc = 1
			}
			if statusMb == 0 && mobiSearch.PointProgress == mobiSearch.PointMax {
				log.Println("手机分数刷取完毕")
				statusMb = 1
			}
		}
		fmt.Println("获取积分完毕！！")
		conn.View.Handler(conn)
		fmt.Println("当前可用分数: ", conn.View.Infov.AvailablePoints)
		fmt.Println("今日可获取最大分数: ", conn.View.Infov.DailyPoints.PointProgressMax)
		fmt.Println("今日分数: ", conn.View.Infov.DailyPoints.PointProgress)
	}()

	// 阻塞更换为执行完毕后休眠5s
	time.Sleep(time.Second * 5)
}
