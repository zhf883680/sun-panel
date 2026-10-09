package panel

import (
	"sun-panel/api/api_v1"
	"sun-panel/api/api_v1/middleware"

	"github.com/gin-gonic/gin"
)

func InitGlobalSetting(router *gin.RouterGroup) {
	api := api_v1.ApiGroupApp.ApiPanel.GlobalSettingApi
	r := router.Group("", middleware.LoginInterceptor)
	{
		r.POST("/panel/globalSetting/getSiteStting", api.GetSiteStting)
		r.POST("/panel/globalSetting/saveSiteStting", api.SaveSiteStting)
		r.POST("/panel/globalSetting/uploadFaviconImage", api.UploadFaviconImage)
		r.POST("/panel/globalSetting/getCustomStylePath", api.GetCustomStylePath)
		r.POST("/panel/globalSetting/getCustomJsAndCssCode", api.GetCustomJsAndCssCode)
		r.POST("/panel/globalSetting/saveCustomJsAndCssCode", api.SaveCustomJsAndCssCode)
	}
}
