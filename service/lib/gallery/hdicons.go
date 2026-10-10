// Package gallery 负责在线图库（HD-Icons）的本地化：
// 首次启动/过期时把索引与图片全量拉取到本地磁盘，之后每天增量更新一次。
// 前端只访问本地数据，支持服务端分页与搜索，不再逐张回源 GitHub。
package gallery

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"sun-panel/global"
)

// HD-Icons 数据源
const (
	hdIconsIndexURL    = "https://raw.githubusercontent.com/xushier/HD-Icons/main/icons.json"
	hdIconsIndexURLAlt = "https://raw.githubusercontent.com/xushier/HD-Icons/master/icons.json"
	hdIconsRawBase     = "https://raw.githubusercontent.com/xushier/HD-Icons/main/"
	hdIconsRawBaseAlt  = "https://raw.githubusercontent.com/xushier/HD-Icons/master/"
)

const (
	// SyncInterval 本地数据超过该时长后重新同步（每天更新）
	SyncInterval = 24 * time.Hour
	// checkInterval 后台检查间隔（到点后按 SyncInterval 判断是否需要同步）
	checkInterval = 30 * time.Minute
	// downloadWorkers 图片下载并发数
	downloadWorkers = 8
	// httpTimeout 单次请求超时
	httpTimeout = 60 * time.Second
)

// HdIconsStyles 支持的图标风格目录
var HdIconsStyles = []string{"svg", "circle", "border-radius"}

// 文件名安全校验，防止路径穿越
var fileRegexp = regexp.MustCompile(`^[A-Za-z0-9_\-]+\.(png|svg|jpg|jpeg|webp|gif|ico)$`)

// Item 在线图库单个图标
type Item struct {
	Name  string `json:"name"`
	Style string `json:"style"`
	File  string `json:"file"`
}

// Manifest 本地索引清单（图片与索引同目录，可整体删除后重建）
type Manifest struct {
	Name       string    `json:"name"`
	UpdateAt   string    `json:"updateAt"`
	SyncedAt   time.Time `json:"syncedAt"`
	TotalCount int       `json:"totalCount"`
	List       []Item    `json:"list"`
}

// Status 在线图库同步状态
type Status struct {
	Syncing    bool      `json:"syncing"`
	SyncedAt   time.Time `json:"syncedAt"`
	UpdateAt   string    `json:"updateAt"`
	TotalCount int       `json:"totalCount"`
	LastError  string    `json:"lastError"`
}

var httpClient = &http.Client{Timeout: httpTimeout}

var (
	manifestMu    sync.RWMutex
	manifestCache *Manifest

	syncMu    sync.Mutex
	syncing   bool
	lastError string
)

// CacheDir 在线图库本地缓存目录：source_path/hdicon（默认 conf/uploads/hdicon）
// 与上传文件同目录，docker 下 conf 为挂载卷可持久化；
// 同时会被 /uploads 静态路由直接提供（/uploads/hdicon/...），浏览器可正常协商缓存。
func CacheDir() string {
	p := global.Config.GetValueString("base", "source_path")
	if p == "" {
		p = "./conf/uploads"
	}
	return filepath.Join(p, "hdicon")
}

// ImageFile 单张图片的本地路径
func ImageFile(style, file string) string {
	return filepath.Join(CacheDir(), style, file)
}

func manifestFile() string {
	return filepath.Join(CacheDir(), "manifest.json")
}

// IsSupportedStyle 校验风格是否合法
func IsSupportedStyle(style string) bool {
	for _, v := range HdIconsStyles {
		if v == style {
			return true
		}
	}
	return false
}

// IsValidFile 校验文件名是否合法（防路径穿越）
func IsValidFile(file string) bool {
	return fileRegexp.MatchString(file)
}

// ContentTypeByFileName 按扩展名推断 Content-Type
func ContentTypeByFileName(file string) string {
	switch strings.ToLower(path.Ext(file)) {
	case ".svg":
		return "image/svg+xml"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	case ".ico":
		return "image/x-icon"
	}
	return "application/octet-stream"
}

// LoadManifest 读取本地索引清单，不存在或损坏时返回 nil
func LoadManifest() *Manifest {
	manifestMu.RLock()
	if manifestCache != nil {
		m := manifestCache
		manifestMu.RUnlock()
		return m
	}
	manifestMu.RUnlock()

	manifestMu.Lock()
	defer manifestMu.Unlock()
	if manifestCache != nil {
		return manifestCache
	}
	b, err := os.ReadFile(manifestFile())
	if err != nil || len(b) == 0 {
		return nil
	}
	m := &Manifest{}
	if err := json.Unmarshal(b, m); err != nil {
		return nil
	}
	if m.List == nil {
		m.List = []Item{}
	}
	manifestCache = m
	return m
}

// GetStatus 返回当前同步状态
func GetStatus() Status {
	manifestMu.RLock()
	m := manifestCache
	manifestMu.RUnlock()
	if m == nil {
		m = LoadManifest()
	}

	syncMu.Lock()
	defer syncMu.Unlock()
	st := Status{Syncing: syncing, LastError: lastError}
	if m != nil {
		st.SyncedAt = m.SyncedAt
		st.UpdateAt = m.UpdateAt
		st.TotalCount = m.TotalCount
	}
	return st
}

// LocalImage 读取已缓存的图片，未命中返回 nil
func LocalImage(style, file string) []byte {
	b, err := os.ReadFile(ImageFile(style, file))
	if err != nil || len(b) == 0 {
		return nil
	}
	return b
}

// EnsureImage 返回图片内容：优先本地缓存，未命中则回源下载并落盘
func EnsureImage(style, file string) ([]byte, error) {
	if b := LocalImage(style, file); b != nil {
		return b, nil
	}
	if !IsSupportedStyle(style) || !IsValidFile(file) {
		return nil, fmt.Errorf("invalid image request")
	}
	item := Item{Style: style, File: file}
	if err := downloadImage(item); err != nil {
		return nil, err
	}
	b := LocalImage(style, file)
	if b == nil {
		return nil, fmt.Errorf("image not available")
	}
	return b, nil
}

func logf(format string, args ...interface{}) {
	if global.Logger != nil {
		global.Logger.Infof(format, args...)
	}
}

func logErrorf(format string, args ...interface{}) {
	if global.Logger != nil {
		global.Logger.Errorf(format, args...)
	}
}

func fetch(u string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Sun-Panel-HdIcons-Sync")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed, status code: %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func fetchAny(urls []string) ([]byte, error) {
	var err error
	for _, u := range urls {
		var body []byte
		body, err = fetch(u)
		if err == nil && len(body) > 0 {
			return body, nil
		}
	}
	if err == nil {
		err = fmt.Errorf("empty response")
	}
	return nil, err
}

// FetchIndex 拉取并解析远端索引
func FetchIndex() (*Manifest, error) {
	body, err := fetchAny([]string{hdIconsIndexURL, hdIconsIndexURLAlt})
	if err != nil {
		return nil, err
	}

	type rawItem struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	type rawIndex struct {
		Name     string    `json:"name"`
		UpdateAt string    `json:"update_at"`
		Icons    []rawItem `json:"icons"`
	}
	raw := rawIndex{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	m := &Manifest{
		Name:     raw.Name,
		UpdateAt: raw.UpdateAt,
		List:     make([]Item, 0, len(raw.Icons)),
	}
	for _, v := range raw.Icons {
		item := Item{Name: v.Name}
		if u, err := url.Parse(v.Url); err == nil {
			item.File = path.Base(u.Path)
			item.Style = path.Base(path.Dir(u.Path))
		}
		if item.File == "" {
			item.File = v.Name + ".png"
		}
		if !IsSupportedStyle(item.Style) {
			item.Style = "border-radius"
		}
		if !IsValidFile(item.File) {
			continue
		}
		m.List = append(m.List, item)
	}
	m.TotalCount = len(m.List)
	return m, nil
}

func writeFileAtomic(filePath string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}
	tmp := filePath + ".tmp"
	if err := os.WriteFile(tmp, data, 0666); err != nil {
		return err
	}
	return os.Rename(tmp, filePath)
}

func downloadImage(item Item) error {
	urls := []string{
		hdIconsRawBase + item.Style + "/" + item.File,
		hdIconsRawBaseAlt + item.Style + "/" + item.File,
	}
	body, err := fetchAny(urls)
	if err != nil {
		return err
	}
	return writeFileAtomic(ImageFile(item.Style, item.File), body)
}

func saveManifest(m *Manifest) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return writeFileAtomic(manifestFile(), b)
}

// Sync 全量同步（阻塞）：拉取索引并下载缺失图片，最后写入本地清单。
// 已存在的图片会跳过，因此重复执行是增量的。
func Sync() (*Manifest, error) {
	syncMu.Lock()
	if syncing {
		syncMu.Unlock()
		return nil, fmt.Errorf("sync already running")
	}
	syncing = true
	lastError = ""
	syncMu.Unlock()
	return runSync()
}

// TriggerSync 异步触发一次同步：若已在同步中返回 false，否则立即返回 true 并在后台执行
func TriggerSync() bool {
	syncMu.Lock()
	if syncing {
		syncMu.Unlock()
		return false
	}
	syncing = true
	lastError = ""
	syncMu.Unlock()
	go func() { _, _ = runSync() }()
	return true
}

func runSync() (*Manifest, error) {
	defer func() {
		syncMu.Lock()
		syncing = false
		syncMu.Unlock()
	}()

	m, err := FetchIndex()
	if err != nil {
		syncMu.Lock()
		lastError = err.Error()
		syncMu.Unlock()
		return nil, err
	}

	start := time.Now()
	downloaded, failed := downloadAll(m.List)
	logf("[gallery] hdicons sync done: total=%d downloaded=%d failed=%d cost=%s",
		len(m.List), downloaded, failed, time.Since(start).Round(time.Second))

	m.SyncedAt = time.Now()
	if err := saveManifest(m); err != nil {
		logErrorf("[gallery] save hdicons manifest error: %v", err)
		return nil, err
	}

	manifestMu.Lock()
	manifestCache = m
	manifestMu.Unlock()
	return m, nil
}

// downloadAll 并发下载缺失图片，返回新下载数量与失败数量
func downloadAll(list []Item) (downloaded, failed int) {
	jobs := make(chan Item)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < downloadWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range jobs {
				if fi, err := os.Stat(ImageFile(item.Style, item.File)); err == nil && fi.Size() > 0 {
					continue
				}
				if err := downloadImage(item); err != nil {
					mu.Lock()
					failed++
					mu.Unlock()
					continue
				}
				mu.Lock()
				downloaded++
				mu.Unlock()
			}
		}()
	}

	for _, item := range list {
		jobs <- item
	}
	close(jobs)
	wg.Wait()
	return downloaded, failed
}

// StartScheduler 启动后台同步：进程启动时先判断一次（首次启动会立即拉取），
// 之后按 checkInterval 轮询，本地数据超过 SyncInterval 时增量更新。
func StartScheduler() {
	go func() {
		scheduledSync()
		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()
		for range ticker.C {
			scheduledSync()
		}
	}()
}

func scheduledSync() {
	if GetStatus().Syncing {
		return
	}
	m := LoadManifest()
	if m != nil && !m.SyncedAt.IsZero() && time.Since(m.SyncedAt) < SyncInterval {
		return
	}
	if m == nil {
		logf("[gallery] hdicons local cache not found, start first sync")
	} else {
		logf("[gallery] hdicons local cache outdated, start sync")
	}
	if _, err := Sync(); err != nil {
		logErrorf("[gallery] hdicons sync error: %v", err)
	}
}
