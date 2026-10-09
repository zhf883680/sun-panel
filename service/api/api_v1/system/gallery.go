package system

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/lib/cache"

	"github.com/gin-gonic/gin"
)

type GalleryApi struct{}

// HD-Icons 图库信息
const (
	hdIconsIndexUrl    = "https://raw.githubusercontent.com/xushier/HD-Icons/main/icons.json"
	hdIconsIndexUrlAlt = "https://raw.githubusercontent.com/xushier/HD-Icons/master/icons.json"
	hdIconsRawBase     = "https://raw.githubusercontent.com/xushier/HD-Icons/main/"
	hdIconsRawBaseAlt  = "https://raw.githubusercontent.com/xushier/HD-Icons/master/"
)

// 支持的图标风格目录
var hdIconsStyles = []string{"svg", "circle", "border-radius"}

// 索引缓存（1小时）
var hdIconsIndexCache = cache.NewGoCache[hdIconsIndexResp](time.Hour, 10*time.Minute)

// 文件名安全校验，防止路径穿越
var hdIconsFileRegexp = regexp.MustCompile(`^[A-Za-z0-9_\-]+\.(png|svg|jpg|jpeg|webp|gif|ico)$`)

type hdIconRawItem struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type hdIconsRawIndex struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	TotalCount  int             `json:"total_count"`
	UpdateAt    string          `json:"update_at"`
	Icons       []hdIconRawItem `json:"icons"`
}

type hdIconItem struct {
	Name  string `json:"name"`
	Style string `json:"style"`
	File  string `json:"file"`
	Url   string `json:"url"`
}

type hdIconsIndexResp struct {
	List       []hdIconItem `json:"list"`
	Name       string       `json:"name"`
	TotalCount int          `json:"totalCount"`
	UpdateAt   string       `json:"updateAt"`
}

func isSupportedHdIconsStyle(style string) bool {
	for _, v := range hdIconsStyles {
		if v == style {
			return true
		}
	}
	return false
}

// 拉取HD-Icons索引
func fetchHdIconsIndex() (hdIconsIndexResp, error) {
	resp := hdIconsIndexResp{List: []hdIconItem{}, TotalCount: 0}

	var body []byte
	var err error
	client := &http.Client{Timeout: 30 * time.Second}
	for _, u := range []string{hdIconsIndexUrl, hdIconsIndexUrlAlt} {
		var r *http.Response
		r, err = client.Get(u)
		if err != nil {
			continue
		}
		func() {
			defer r.Body.Close()
			if r.StatusCode != http.StatusOK {
				err = fmt.Errorf("HTTP request failed, status code: %d", r.StatusCode)
				return
			}
			body, err = io.ReadAll(r.Body)
		}()
		if err == nil && len(body) > 0 {
			break
		}
	}
	if err != nil {
		return resp, err
	}
	if len(body) == 0 {
		return resp, fmt.Errorf("empty index data")
	}

	raw := hdIconsRawIndex{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return resp, err
	}
	resp.Name = raw.Name
	resp.UpdateAt = raw.UpdateAt

	for _, v := range raw.Icons {
		item := hdIconItem{Name: v.Name, Url: v.Url}
		if u, err := url.Parse(v.Url); err == nil {
			item.File = path.Base(u.Path)
			item.Style = path.Base(path.Dir(u.Path))
		}
		if item.File == "" {
			item.File = v.Name + ".png"
		}
		if !isSupportedHdIconsStyle(item.Style) {
			item.Style = "border-radius"
		}
		resp.List = append(resp.List, item)
	}
	resp.TotalCount = len(resp.List)
	return resp, nil
}

// 获取HD-Icons索引（带缓存）
func (a *GalleryApi) HdIconsIndex(c *gin.Context) {
	if v, ok := hdIconsIndexCache.Get("index"); ok {
		apiReturn.SuccessData(c, v)
		return
	}

	resp, err := fetchHdIconsIndex()
	if err != nil {
		global.Logger.Error("[gallery] fetch hdicons index error:", err.Error())
		apiReturn.Error(c, "Failed to fetch the online gallery index: "+err.Error())
		return
	}

	hdIconsIndexCache.Set("index", resp, time.Hour)
	apiReturn.SuccessData(c, resp)
}

// 代理HD-Icons图片，解决跨域/速度问题
func (a *GalleryApi) HdIconsImage(c *gin.Context) {
	style := c.Query("style")
	file := c.Query("file")

	if !isSupportedHdIconsStyle(style) || !hdIconsFileRegexp.MatchString(file) {
		c.Status(http.StatusBadRequest)
		return
	}

	client := &http.Client{Timeout: 30 * time.Second}
	var resp *http.Response
	var err error
	for _, base := range []string{hdIconsRawBase, hdIconsRawBaseAlt} {
		resp, err = client.Get(base + style + "/" + file)
		if err != nil {
			continue
		}
		if resp.StatusCode == http.StatusOK {
			break
		}
		resp.Body.Close()
		err = fmt.Errorf("HTTP request failed, status code: %d", resp.StatusCode)
	}
	if err != nil || resp == nil {
		global.Logger.Error("[gallery] fetch hdicons image error:", err)
		c.Status(http.StatusNotFound)
		return
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" || strings.HasPrefix(contentType, "text/") {
		if strings.HasSuffix(file, ".svg") {
			contentType = "image/svg+xml"
		} else {
			contentType = "image/png"
		}
	}

	c.Header("Cache-Control", "public, max-age=86400")
	c.DataFromReader(http.StatusOK, resp.ContentLength, contentType, resp.Body, nil)
}
