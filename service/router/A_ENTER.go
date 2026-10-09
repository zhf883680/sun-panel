package router

import (
	"strings"
	"sun-panel/global"
	// "sun-panel/router/admin"
	"sun-panel/router/openness"
	"sun-panel/router/panel"
	"sun-panel/router/system"

	"github.com/gin-gonic/gin"
)

// 初始化总路由
func InitRouters(addr string) error {
	router := gin.Default()
	rootRouter := router.Group("/")
	routerGroup := rootRouter.Group("api")

	// 接口
	system.Init(routerGroup)
	panel.Init(routerGroup)
	openness.Init(routerGroup)

	// WEB文件服务
	{
		webPath := global.Config.GetValueStringOrDefault("base", "web_path")
		if webPath == "" {
			webPath = "./web"
		}
		// 自定义css/js路径，默认使用 web/custom
		customStylePath := global.Config.GetValueStringOrDefault("base", "custom_style_path")
		if customStylePath == "" {
			customStylePath = webPath + "/custom"
		}
		router.StaticFile("/", webPath+"/index.html")
		router.Static("/assets", webPath+"/assets")
		router.Static("/custom", customStylePath)
		router.StaticFile("/favicon.ico", webPath+"/favicon.ico")
		router.StaticFile("/favicon.svg", webPath+"/favicon.svg")
	}

	// 上传的文件
	sourcePath := global.Config.GetValueStringOrDefault("base", "source_path")
	if sourcePath == "" {
		sourcePath = "./conf/uploads"
	}
	// 按 source_path 提供访问（如 /conf/uploads）
	router.Static(strings.TrimPrefix(sourcePath, "."), sourcePath)
	// 兼容历史数据中形如 /uploads/xxx 的图片地址
	router.Static("/uploads", sourcePath)

	global.Logger.Info("Sun-Panel is Started.  Listening and serving HTTP on ", addr)
	return router.Run(addr)
}
