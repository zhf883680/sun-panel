package openness

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/lib/cmn/systemSetting"

	"github.com/gin-gonic/gin"
)

type Openness struct {
}

func (a *Openness) LoginConfig(c *gin.Context) {
	cfg := systemSetting.ApplicationSetting{}
	if err := global.SystemSetting.GetValueByInterface(systemSetting.SYSTEM_APPLICATION, &cfg); err != nil {
		apiReturn.Error(c, "配置查询失败："+err.Error())
		return
	}
	apiReturn.SuccessData(c, gin.H{
		"loginCaptcha": cfg.LoginCaptcha,
		"register":     cfg.Register,
	})
}

func (a *Openness) GetDisclaimer(c *gin.Context) {
	if content, err := global.SystemSetting.GetValueString(systemSetting.DISCLAIMER); err != nil {
		global.SystemSetting.Set(systemSetting.DISCLAIMER, "")
		apiReturn.SuccessData(c, "")
		return
	} else {
		apiReturn.SuccessData(c, content)
	}
}

func (a *Openness) GetAboutDescription(c *gin.Context) {
	if content, err := global.SystemSetting.GetValueString(systemSetting.WEB_ABOUT_DESCRIPTION); err != nil {
		global.SystemSetting.Set(systemSetting.WEB_ABOUT_DESCRIPTION, "")
		apiReturn.SuccessData(c, "")
		return
	} else {
		apiReturn.SuccessData(c, content)
	}
}

// 站点设置（登录页需要，无需鉴权）
func (a *Openness) GetSiteSetting(c *gin.Context) {
	siteSetting := systemSetting.DefaultSiteSetting()
	_ = global.SystemSetting.GetValueByInterface(systemSetting.SITE_SETTING, &siteSetting)
	apiReturn.SuccessData(c, siteSetting)
}
