package panel

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/cmn/systemSetting"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 全局设置（站点设置、自定义JS/CSS）
type GlobalSettingApi struct {
}

// 获取站点设置
func (a *GlobalSettingApi) GetSiteStting(c *gin.Context) {
	siteSetting := systemSetting.DefaultSiteSetting()
	_ = global.SystemSetting.GetValueByInterface(systemSetting.SITE_SETTING, &siteSetting)
	apiReturn.SuccessData(c, siteSetting)
}

// 保存站点设置
func (a *GlobalSettingApi) SaveSiteStting(c *gin.Context) {
	req := systemSetting.SiteSetting{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	if err := global.SystemSetting.Set(systemSetting.SITE_SETTING, req); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}

// 上传站点图标（favicon）
func (a *GlobalSettingApi) UploadFaviconImage(c *gin.Context) {
	configUpload := global.Config.GetValueString("base", "source_path")
	f, err := c.FormFile("imgfile")
	if err != nil {
		apiReturn.ErrorByCode(c, 1300)
		return
	}
	fileExt := strings.ToLower(path.Ext(f.Filename))
	agreeExts := []string{".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".ico"}
	if !cmn.InArray(agreeExts, fileExt) {
		apiReturn.ErrorByCode(c, 1301)
		return
	}

	fileName := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
	fileDir := fmt.Sprintf("%s/custom/", configUpload)
	if isExist, _ := cmn.PathExists(fileDir); !isExist {
		os.MkdirAll(fileDir, os.ModePerm)
	}
	filePath := fmt.Sprintf("%s%s%s", fileDir, fileName, fileExt)
	if err := c.SaveUploadedFile(f, filePath); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.SuccessData(c, gin.H{
		"imageUrl": filePath[1:],
	})
}

// 自定义JS/CSS所在目录（站点访问路径）
func (a *GlobalSettingApi) GetCustomStylePath(c *gin.Context) {
	apiReturn.SuccessData(c, gin.H{
		"customPath": customStyleUrlPath(),
	})
}

// 读取自定义JS/CSS内容
func (a *GlobalSettingApi) GetCustomJsAndCssCode(c *gin.Context) {
	apiReturn.SuccessData(c, gin.H{
		"jsContent":  readFileString(customStyleDir(), "index.js"),
		"cssContent": readFileString(customStyleDir(), "index.css"),
	})
}

// 保存自定义JS/CSS内容
func (a *GlobalSettingApi) SaveCustomJsAndCssCode(c *gin.Context) {
	req := struct {
		JsContent  string `json:"jsContent"`
		CssContent string `json:"cssContent"`
	}{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	dir := customStyleDir()
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	if err := os.WriteFile(path.Join(dir, "index.js"), []byte(req.JsContent), 0666); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	if err := os.WriteFile(path.Join(dir, "index.css"), []byte(req.CssContent), 0666); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}

// 自定义样式目录（磁盘路径）
func customStyleDir() string {
	p := global.Config.GetValueString("base", "custom_style_path")
	if p == "" {
		webPath := global.Config.GetValueString("base", "web_path")
		if webPath == "" {
			webPath = "./web"
		}
		p = webPath + "/custom"
	}
	return p
}

// 自定义样式目录（站点访问路径）
func customStyleUrlPath() string {
	p := strings.TrimPrefix(customStyleDir(), ".")
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return p
}

func readFileString(dir, name string) string {
	b, err := os.ReadFile(path.Join(dir, name))
	if err != nil {
		return ""
	}
	return string(b)
}
