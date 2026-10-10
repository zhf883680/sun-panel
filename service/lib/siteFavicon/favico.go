package siteFavicon

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sun-panel/lib/cmn"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// 浏览器默认图标地址：页面没有显式声明图标时回退到这里
const defaultFaviconPath = "/favicon.ico"

func IsHTTPURL(url string) bool {
	httpPattern := `^(http://|https://|//)`
	match, err := regexp.MatchString(httpPattern, url)
	if err != nil {
		return false
	}
	return match
}

// 图标 rel 优先级：数值越小越优先。
// 例如 "icon" / "shortcut icon" 优先于 "apple-touch-icon"、"mask-icon"，
// 避免把 SPA 里常见的 apple-touch-icon 当成首选项导致下载失败。
func iconRelPriority(rel string) int {
	rel = strings.ToLower(strings.TrimSpace(rel))
	switch {
	case rel == "icon" || rel == "shortcut icon":
		return 0
	case strings.Contains(rel, "apple-touch-icon"):
		return 2
	case strings.Contains(rel, "mask-icon"):
		return 3
	case strings.Contains(rel, "icon"):
		return 1
	default:
		return 99
	}
}

// 以页面地址（并考虑 <base href>）为基准，把 href 解析成绝对地址
func resolveIconURL(base *url.URL, href string) string {
	href = strings.TrimSpace(href)
	if base == nil || href == "" {
		return href
	}
	if strings.HasPrefix(href, "//") {
		if base.Scheme == "" {
			return href
		}
		return base.Scheme + ":" + href
	}
	ref, err := url.Parse(href)
	if err != nil {
		return href
	}
	return base.ResolveReference(ref).String()
}

// 计算文档的基准地址（支持 <base href>）
func documentBaseURL(pageURL string, doc *goquery.Document) *url.URL {
	base, err := url.Parse(pageURL)
	if err != nil {
		return nil
	}
	if doc != nil {
		if href, ok := doc.Find("base[href]").First().Attr("href"); ok {
			if ref, err := url.Parse(strings.TrimSpace(href)); err == nil {
				base = base.ResolveReference(ref)
			}
		}
	}
	return base
}

// 收集页面声明的图标地址，并按优先级排序
func collectIconLinks(doc *goquery.Document, base *url.URL) []string {
	type iconLink struct {
		href string
		rank int
	}

	links := make([]iconLink, 0)
	if doc != nil {
		doc.Find("link").Each(func(i int, s *goquery.Selection) {
			rel, _ := s.Attr("rel")
			href, _ := s.Attr("href")
			href = strings.TrimSpace(href)
			if href == "" {
				return
			}
			rank := iconRelPriority(rel)
			if rank >= 99 {
				return
			}
			links = append(links, iconLink{href: href, rank: rank})
		})
	}
	sort.SliceStable(links, func(i, j int) bool { return links[i].rank < links[j].rank })

	icons := make([]string, 0, len(links))
	for _, l := range links {
		icons = append(icons, resolveIconURL(base, l.href))
	}
	return icons
}

// 页面未声明图标时的兜底地址：站点根目录下的 /favicon.ico
func defaultFaviconURL(base *url.URL) string {
	if base == nil || base.Host == "" {
		return ""
	}
	return base.ResolveReference(&url.URL{Path: defaultFaviconPath}).String()
}

func GetOneFaviconURL(urlStr string) (string, error) {
	iconURLs, err := getFaviconURL(urlStr)
	if err != nil {
		return "", err
	}
	if len(iconURLs) == 0 {
		return "", fmt.Errorf("not found ico")
	}
	return iconURLs[0], nil
}

// 从网页上的link标签中获取图标的完整地址，同时返回网页标题
// 即使图标获取失败，只要页面能访问，也会尽可能返回标题
func GetOneFaviconURLAndTitle(urlStr string) (iconUrl string, title string, err error) {
	doc, err := getDocument(urlStr)
	if err != nil {
		return "", "", err
	}

	title = strings.TrimSpace(doc.Find("title").First().Text())

	base := documentBaseURL(urlStr, doc)
	iconURLs := collectIconLinks(doc, base)
	if len(iconURLs) == 0 {
		// 页面没有显式声明图标时，回退到浏览器默认的 /favicon.ico
		if fallback := defaultFaviconURL(base); fallback != "" {
			return fallback, title, nil
		}
		return "", title, errors.New("favicon not found on the page")
	}

	return iconURLs[0], title, nil
}

// 获取远程文件的大小
func GetRemoteFileSize(url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP request failed, status code: %d", resp.StatusCode)
	}

	// 获取Content-Length字段，即文件大小
	size := resp.ContentLength
	return size, nil
}

// 下载图片
func DownloadImage(imgURL, savePath string, maxSize int64) (*os.File, error) {
	// 获取远程文件大小。部分服务器不支持 HEAD，这里不作为失败条件，仅用于大小限制。
	if fileSize, err := GetRemoteFileSize(imgURL); err == nil && fileSize > maxSize {
		return nil, fmt.Errorf("文件太大，不下载。大小：%d字节", fileSize)
	}

	// 发送HTTP GET请求获取图片数据
	response, err := http.Get(imgURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// 检查HTTP响应状态
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed, status code: %d", response.StatusCode)
	}

	// 从 URL 的 path 部分推断文件名与扩展名，避免把 query 参数当成扩展名
	urlPath := imgURL
	if u, perr := url.Parse(imgURL); perr == nil {
		urlPath = u.Path
	}
	urlFileName := path.Base(urlPath)
	fileExt := strings.ToLower(path.Ext(urlPath))
	fileName := cmn.Md5(fmt.Sprintf("%s%s", urlFileName, time.Now().String())) + fileExt

	// 使用单斜杠拼接，避免出现 "//" 导致静态路由重定向
	destination := strings.TrimRight(savePath, "/") + "/" + fileName

	// 创建本地文件用于保存图片
	file, err := os.Create(destination)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 将图片数据写入本地文件
	if _, err = io.Copy(file, response.Body); err != nil {
		return nil, err
	}
	return file, nil
}

func GetOneFaviconURLAndUpload(urlStr string) (string, bool) {
	iconURLs, err := getFaviconURL(urlStr)
	if err != nil || len(iconURLs) == 0 {
		return "", false
	}
	return iconURLs[0], true
}

func getFaviconURL(urlStr string) ([]string, error) {
	doc, err := getDocument(urlStr)
	if err != nil {
		return nil, err
	}

	base := documentBaseURL(urlStr, doc)
	icons := collectIconLinks(doc, base)
	if len(icons) == 0 {
		// 没有显式声明图标时，回退到 /favicon.ico
		if fallback := defaultFaviconURL(base); fallback != "" {
			return []string{fallback}, nil
		}
		return icons, errors.New("favicon not found on the page")
	}

	return icons, nil
}

// 请求页面并解析为 goquery 文档
func getDocument(url string) (*goquery.Document, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置User-Agent头字段，模拟浏览器请求
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("HTTP request failed with status code " + strconv.Itoa(resp.StatusCode))
	}

	return goquery.NewDocumentFromReader(resp.Body)
}
