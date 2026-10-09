package models

import (
	"gorm.io/gorm"
)

type ItemIconGroup struct {
	BaseModel
	Icon        string `json:"icon"`
	Title       string `gorm:"type:varchar(50)" json:"title"`
	Description string `gorm:"type:varchar(1000)" json:"description"`
	Sort        int    `gorm:"type:int(11)" json:"sort"`
	UserId      uint   `json:"userId"`
	User        User   `json:"user"`
	// 分组卡片仅显示名称
	OnlyName string `gorm:"type:varchar(50);default:''" json:"-"`
	// 分组卡片公开访问时隐藏
	PublicVisitHide bool `gorm:"default:false" json:"-"`
	// 分组卡片样式（JSON）
	CardStyle string `gorm:"type:text" json:"-"`
}

func (m *ItemIconGroup) DeleteByUserId(db *gorm.DB, userId uint) (err error) {
	err = db.Delete(&ItemIconGroup{}, "user_id = ?", userId).Error
	return
}
