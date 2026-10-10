package system

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitGalleryRouter(router *gin.RouterGroup) {
	GalleryApi := api_v1.ApiGroupApp.ApiSystem.GalleryApi

	// 图库为公开数据：索引、分页列表与图片均不校验登录
	// （图片仅能访问固定的 HD-Icons 仓库且经过白名单校验）
	router.GET("/system/gallery/hdicons/index", GalleryApi.HdIconsIndex)
	router.GET("/system/gallery/hdicons/list", GalleryApi.HdIconsList)
	router.GET("/system/gallery/hdicons/image", GalleryApi.HdIconsImage)

	// 手动同步属于系统操作，需要登录
	private := router.Group("", middleware.LoginInterceptor)
	{
		private.GET("/system/gallery/hdicons/syncStatus", GalleryApi.HdIconsSyncStatus)
		private.POST("/system/gallery/hdicons/sync", GalleryApi.HdIconsSync)
	}
}
