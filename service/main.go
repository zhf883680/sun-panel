package main

import (
	"log"
	"sun-panel/global"
	"sun-panel/initialize"
	"sun-panel/lib/gallery"
	"sun-panel/router"
)

func main() {
	err := initialize.InitApp()
	if err != nil {
		log.Println("初始化错误:", err.Error())
		panic(err)
	}

	// 在线图库本地化：首次启动拉取全部图标并落盘，之后每天增量更新
	gallery.StartScheduler()

	httpPort := global.Config.GetValueStringOrDefault("base", "http_port")

	if err := router.InitRouters(":" + httpPort); err != nil {
		panic(err)
	}
}
