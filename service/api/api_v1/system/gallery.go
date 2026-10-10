package system

import (
	"net/http"
	"strconv"
	"strings"

	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/lib/gallery"

	"github.com/gin-gonic/gin"
)

type GalleryApi struct{}

// HdIconsIndex 返回完整索引（兼容旧接口）。
// 数据来自本地清单，未同步完成时回源一次，避免首次启动空列表。
func (a *GalleryApi) HdIconsIndex(c *gin.Context) {
	m := gallery.LoadManifest()
	if m == nil {
		// 首次启动同步尚未完成：触发后台同步并返回空列表，前端稍后重试
		gallery.TriggerSync()
		apiReturn.SuccessData(c, gin.H{
			"name":       "",
			"updateAt":   "",
			"syncedAt":   nil,
			"totalCount": 0,
			"list":       []gallery.Item{},
			"syncing":    true,
		})
		return
	}
	apiReturn.SuccessData(c, gin.H{
		"name":       m.Name,
		"updateAt":   m.UpdateAt,
		"syncedAt":   m.SyncedAt,
		"totalCount": m.TotalCount,
		"list":       m.List,
	})
}

// HdIconsList 本地分页 + 搜索 + 风格筛选，前端滚动加载更多。
// GET /system/gallery/hdicons/list?page=1&pageSize=60&keyword=&style=
func (a *GalleryApi) HdIconsList(c *gin.Context) {
	page := parsePositiveInt(c.Query("page"), 1)
	pageSize := parsePositiveInt(c.Query("pageSize"), 60)
	if pageSize > 200 {
		pageSize = 200
	}
	keyword := strings.ToLower(strings.TrimSpace(c.Query("keyword")))
	style := strings.TrimSpace(c.Query("style"))

	m := gallery.LoadManifest()
	status := gallery.GetStatus()
	if m == nil {
		// 首次启动尚未完成同步，返回空列表并标记 syncing，前端据此提示与重试
		apiReturn.SuccessData(c, gin.H{
			"list":       []gallery.Item{},
			"total":      0,
			"page":       page,
			"pageSize":   pageSize,
			"syncing":    status.Syncing,
			"syncedAt":   status.SyncedAt,
			"updateAt":   status.UpdateAt,
			"totalCount": 0,
		})
		return
	}

	filtered := make([]gallery.Item, 0, len(m.List))
	for _, item := range m.List {
		if style != "" && style != "all" && item.Style != style {
			continue
		}
		if keyword != "" && !strings.Contains(strings.ToLower(item.Name), keyword) {
			continue
		}
		filtered = append(filtered, item)
	}

	total := len(filtered)
	start := (page - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	apiReturn.SuccessData(c, gin.H{
		"list":       filtered[start:end],
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
		"syncing":    status.Syncing,
		"syncedAt":   status.SyncedAt,
		"updateAt":   m.UpdateAt,
		"totalCount": m.TotalCount,
	})
}

// HdIconsSyncStatus 查询当前同步状态（不触发同步）
func (a *GalleryApi) HdIconsSyncStatus(c *gin.Context) {
	apiReturn.SuccessData(c, gallery.GetStatus())
}

// HdIconsSync 手动触发一次同步（后台执行，立即返回）
func (a *GalleryApi) HdIconsSync(c *gin.Context) {
	started := gallery.TriggerSync()
	apiReturn.SuccessData(c, gin.H{"started": started})
}

// HdIconsImage 提供本地图片，未缓存时回源一次并落盘
func (a *GalleryApi) HdIconsImage(c *gin.Context) {
	style := c.Query("style")
	file := c.Query("file")

	if !gallery.IsSupportedStyle(style) || !gallery.IsValidFile(file) {
		c.Status(http.StatusBadRequest)
		return
	}

	body, err := gallery.EnsureImage(style, file)
	if err != nil {
		global.Logger.Error("[gallery] hdicons image error: ", err.Error())
		c.Status(http.StatusNotFound)
		return
	}

	c.Header("Cache-Control", "public, max-age=86400")
	c.Data(http.StatusOK, gallery.ContentTypeByFileName(file), body)
}

func parsePositiveInt(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return fallback
	}
	return v
}
