package system

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type FileApi struct{}

func (a *FileApi) UploadImg(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	configUpload := global.Config.GetValueString("base", "source_path")
	f, err := c.FormFile("imgfile")
	if err != nil {
		apiReturn.ErrorByCode(c, 1300)
		return
	} else {
		fileExt := strings.ToLower(path.Ext(f.Filename))
		agreeExts := []string{
			".png",
			".jpg",
			".gif",
			".jpeg",
			".webp",
			".svg",
			".ico",
		}

		if !cmn.InArray(agreeExts, fileExt) {
			apiReturn.ErrorByCode(c, 1301)
			return
		}
		fileName := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
		fildDir := fmt.Sprintf("%s/%d/%d/%d/", configUpload, time.Now().Year(), time.Now().Month(), time.Now().Day())
		isExist, _ := cmn.PathExists(fildDir)
		if !isExist {
			os.MkdirAll(fildDir, os.ModePerm)
		}
		filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
		c.SaveUploadedFile(f, filepath)

		// 图库分类：前端传 type（v1.8.1：0未知/1壁纸/2图标）或 category（icon/wallpaper/unknown）
		category := c.PostForm("category")
		fileType := cmn.StrToInt(c.PostForm("type"))
		if category == "" {
			switch fileType {
			case 1:
				category = "wallpaper"
			case 2:
				category = "icon"
			}
		}

		// 像数据库添加记录
		if category == "" {
			category = "icon"
		}
		mFile := models.File{}
		mFile.AddFile(userInfo.ID, f.Filename, fileExt, filepath, f.Size, category)
		apiReturn.SuccessData(c, gin.H{
			"imageUrl": filepath[1:],
			"category": category,
			"size":     f.Size,
		})
	}
}

func (a *FileApi) UploadFiles(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	configUpload := global.Config.GetValueString("base", "source_path")

	form, err := c.MultipartForm()
	if err != nil {
		apiReturn.ErrorByCode(c, 1300)
		return
	}
	files := form.File["files[]"]
	category := c.PostForm("category")
	errFiles := []string{}
	succMap := map[string]string{}
	for _, f := range files {
		fileExt := strings.ToLower(path.Ext(f.Filename))
		fileName := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
		fildDir := fmt.Sprintf("%s/%d/%d/%d/", configUpload, time.Now().Year(), time.Now().Month(), time.Now().Day())
		isExist, _ := cmn.PathExists(fildDir)
		if !isExist {
			os.MkdirAll(fildDir, os.ModePerm)
		}
		filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
		if c.SaveUploadedFile(f, filepath) != nil {
			errFiles = append(errFiles, f.Filename)
		} else {
			// 成功
			// 像数据库添加记录
			mFile := models.File{}
			mFile.AddFile(userInfo.ID, f.Filename, fileExt, filepath, f.Size, category)
			succMap[f.Filename] = filepath[1:]
		}
	}

	apiReturn.SuccessData(c, gin.H{
		"succMap":  succMap,
		"errFiles": errFiles,
	})
}

func (a *FileApi) GetList(c *gin.Context) {
	type GetListReq struct {
		Category string `json:"category"`
	}
	req := GetListReq{}
	// 参数可选，解析失败时按全部处理
	_ = c.ShouldBindBodyWith(&req, binding.JSON)

	list := []models.File{}
	userInfo, _ := base.GetCurrentUserInfo(c)
	var count int64
	buildQuery := func() *gorm.DB {
		query := global.Db.Model(&models.File{}).Where("user_id=?", userInfo.ID)
		switch req.Category {
		case "", "all":
			// 不限制分类
		case "wallpaper":
			query = query.Where("category=?", "wallpaper")
		default:
			// 图标：v1.8.1 之前上传的图片没有分类（NULL/空），默认按图标处理，避免升级后图库为空
			query = query.Where("(category=? OR category IS NULL OR category='')", "icon")
		}
		return query
	}
	if err := buildQuery().Count(&count).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	if err := buildQuery().Order("created_at desc").Find(&list).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	data := []map[string]interface{}{}
	for _, v := range list {
		data = append(data, map[string]interface{}{
			"src":             normalizeFileSrc(v.Src),
			"fileName":        v.FileName,
			"id":              v.ID,
			"createTime":      v.CreatedAt,
			"updateTime":      v.UpdatedAt,
			"path":            v.Src,
			"ext":             v.Ext,
			"category":        v.Category,
			"size":            v.Size,
			"type":            v.Type,
			"isPublicGallery": v.IsPublicGallery,
		})
	}
	apiReturn.SuccessListData(c, data, count)
}

// 统一图片访问地址（均为基于站点根目录的绝对路径）：
// "./uploads/xxx"（v1.3 全路径）与 "./conf/uploads/xxx"（按 source_path 的静态目录）去掉开头的点；
// "/2026/9/24/xxx" 这类相对上传根目录的历史数据补上 /uploads 前缀。
func normalizeFileSrc(src string) string {
	if strings.HasPrefix(src, ".") {
		return src[1:]
	}
	if strings.HasPrefix(src, "/") {
		return "/uploads" + src
	}
	return src
}

// 设置图库分类（icon / wallpaper / unknown）
func (a *FileApi) UpdateCategory(c *gin.Context) {
	type UpdateCategoryReq struct {
		Id       uint   `json:"id"`
		Category string `json:"category"`
	}
	req := UpdateCategoryReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	userInfo, _ := base.GetCurrentUserInfo(c)

	if err := global.Db.Model(&models.File{}).Where("id=? AND user_id=?", req.Id, userInfo.ID).
		Updates(map[string]interface{}{
			"category": req.Category,
			"type":     models.CategoryToType(req.Category),
		}).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

// 与 v1.8.1 对齐的文件类型设置：type 0未知 1壁纸 2图标
func (a *FileApi) UpdateFileType(c *gin.Context) {
	type UpdateFileTypeReq struct {
		Id   uint `json:"id"`
		Type int  `json:"type"`
	}
	req := UpdateFileTypeReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	category := "unknown"
	switch req.Type {
	case 1:
		category = "wallpaper"
	case 2:
		category = "icon"
	}
	userInfo, _ := base.GetCurrentUserInfo(c)
	if err := global.Db.Model(&models.File{}).Where("id=? AND user_id=?", req.Id, userInfo.ID).
		Updates(map[string]interface{}{
			"category": category,
			"type":     req.Type,
		}).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}

func (a *FileApi) Deletes(c *gin.Context) {
	req := commonApiStructs.RequestDeleteIds[uint]{}
	userInfo, _ := base.GetCurrentUserInfo(c)
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	global.Db.Transaction(func(tx *gorm.DB) error {
		files := []models.File{}

		if err := tx.Order("created_at desc").Find(&files, "user_id=? AND id in ?", userInfo.ID, req.Ids).Error; err != nil {
			return err
		}

		for _, v := range files {
			os.Remove(v.Src)
		}

		if err := tx.Order("created_at desc").Delete(&files, "user_id=? AND id in ?", userInfo.ID, req.Ids).Error; err != nil {
			return err
		}

		return nil
	})

	apiReturn.Success(c)

}
