package system

import (
	"sun-panel/api/api_v1"

	"github.com/gin-gonic/gin"
)

func InitGalleryRouter(router *gin.RouterGroup) {
	GalleryApi := api_v1.ApiGroupApp.ApiSystem.GalleryApi

	// 与 v1.8.1 保持一致：图库为公开数据，索引与图片代理均不校验登录
	// （图片代理只能访问固定的 HD-Icons 仓库且经过白名单校验）
	router.GET("/system/gallery/hdicons/index", GalleryApi.HdIconsIndex)
	router.GET("/system/gallery/hdicons/image", GalleryApi.HdIconsImage)
}
