package system

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
)

type About struct {
}

// 进程启动时间
var processStartTime = time.Now()

func (a *About) Get(c *gin.Context) {
	version := cmn.GetSysVersionInfo()
	apiReturn.SuccessData(c, gin.H{
		"versionName": version.Version,
		"versionCode": version.Version_code,
	})
}

// 运行信息
func (a *About) RunningInfo(c *gin.Context) {
	version := cmn.GetSysVersionInfo()

	var userCount, groupCount, itemCount, fileCount int64
	var fileSize int64
	global.Db.Model(&models.User{}).Count(&userCount)
	global.Db.Model(&models.ItemIconGroup{}).Count(&groupCount)
	global.Db.Model(&models.ItemIcon{}).Count(&itemCount)
	global.Db.Model(&models.File{}).Count(&fileCount)
	global.Db.Model(&models.File{}).Select("COALESCE(SUM(size),0)").Scan(&fileSize)

	runTime := time.Since(processStartTime)
	apiReturn.SuccessData(c, gin.H{
		"versionName": version.Version,
		"versionCode": version.Version_code,
		"goVersion":   runtime.Version(),
		"os":          runtime.GOOS,
		"arch":        runtime.GOARCH,
		"isDocker":    global.ISDOCKER == "docker",
		"database":    global.DB_DRIVER,
		"runTime":     formatDuration(runTime),
		"runSeconds":  int64(runTime.Seconds()),
		"startTime":   processStartTime.Format("2006-01-02 15:04:05"),
		"userCount":   userCount,
		"groupCount":  groupCount,
		"itemCount":   itemCount,
		"fileCount":   fileCount,
		"fileSize":    fileSize,
	})
}

// 检查新版本（查询 GitHub 最新 release）
func (a *About) CheckNewVersion(c *gin.Context) {
	req := struct {
		IsCheckBeta bool `json:"isCheckBeta"`
	}{}
	_ = c.ShouldBindJSON(&req)

	version := cmn.GetSysVersionInfo()
	hasNewVersion := false
	newVersion := ""
	releaseNote := ""

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/zhf883680/sun-panel/releases?per_page=15")
	if err == nil {
		defer resp.Body.Close()
		var releases []struct {
			TagName    string `json:"tag_name"`
			Body       string `json:"body"`
			Prerelease bool   `json:"prerelease"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&releases); err == nil {
			for _, r := range releases {
				if r.Prerelease && !req.IsCheckBeta {
					continue
				}
				tag := strings.TrimPrefix(r.TagName, "v")
				if tag == "" {
					continue
				}
				if compareVersion(tag, version.Version) > 0 {
					hasNewVersion = true
					newVersion = tag
					releaseNote = r.Body
				}
				break
			}
		}
	}

	apiReturn.SuccessData(c, gin.H{
		"hasNewVersion": hasNewVersion,
		"newVersion":    newVersion,
		"nowVersion":    version.Version,
		"releaseNote":   releaseNote,
		"downloadUrl":   "https://github.com/zhf883680/sun-panel/releases",
	})
}

// 简单版本号比较：0 相等，1 前者大，-1 后者大
func compareVersion(a, b string) int {
	pa := strings.Split(strings.TrimSpace(a), ".")
	pb := strings.Split(strings.TrimSpace(b), ".")
	for i := 0; i < len(pa) || i < len(pb); i++ {
		na, nb := 0, 0
		if i < len(pa) {
			na, _ = strconv.Atoi(strings.TrimFunc(pa[i], func(r rune) bool { return r < '0' || r > '9' }))
		}
		if i < len(pb) {
			nb, _ = strconv.Atoi(strings.TrimFunc(pb[i], func(r rune) bool { return r < '0' || r > '9' }))
		}
		if na > nb {
			return 1
		}
		if na < nb {
			return -1
		}
	}
	return 0
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	parts := []string{}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d天", days))
	}
	if hours > 0 || days > 0 {
		parts = append(parts, fmt.Sprintf("%d小时", hours))
	}
	parts = append(parts, fmt.Sprintf("%d分钟", minutes))
	return strings.Join(parts, " ")
}
