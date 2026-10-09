package models

// 默认面板配置（用户从未设置过面板样式时返回），与前端 defaultStatePanelConfig 保持一致
const DefaultPanelJson = `{"backgroundBlur":0,"backgroundMaskNumber":0,"clockShow":true,"clockShowSecond":false,"iconStyle":1,"iconTextColor":"#ffffff","iconTextIconHideTitle":false,"iconTextInfoHideDescription":false,"logoImageSrc":"","logoShow":true,"logoText":"Sun-Panel","marginBottom":10,"marginTop":10,"marginX":5,"maxWidth":1200,"maxWidthUnit":"px","netModeChangeButtonShow":true,"accountSwitcherShow":false,"searchBoxSearchIcon":false,"searchBoxShow":false}`

type UserConfig struct {
	UserId uint `gorm:"index" json:"userId"`

	// 纯前端数据，面板样式数据
	PanelJson string                 `json:"-"`
	Panel     map[string]interface{} `gorm:"-" json:"panel"`

	// 搜索引擎
	SearchEngineJson string                 `json:"-"`
	SearchEngine     map[string]interface{} `gorm:"-" json:"searchEngine"`
}
