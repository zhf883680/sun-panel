package models

import (
	"sun-panel/models/datatype"

	"gorm.io/gorm"
)

type ItemIcon struct {
	BaseModel
	IconJson        string                    `gorm:"type:varchar(1000)" json:"-"`
	Icon            datatype.ItemIconIconInfo `gorm:"-" json:"icon"`
	Title           string                    `gorm:"type:varchar(50)" json:"title"`
	Url             string                    `gorm:"type:varchar(1000)" json:"url"`
	LanUrl          string                    `gorm:"type:varchar(1000)" json:"lanUrl"`
	Description     string                    `gorm:"type:varchar(1000)" json:"description"`
	OpenMethod      int                       `gorm:"type:tinyint(1)" json:"openMethod"`
	Sort            int                       `gorm:"type:int(11)" json:"sort"`
	ItemIconGroupId int                       `json:"itemIconGroupId"`
	UserId          uint                      `json:"userId"`
	User            User                      `json:"user"`
	// 仅显示名称（不显示图标）
	OnlyName string `gorm:"type:varchar(50);default:''" json:"onlyName"`
	// 卡片类型 1:详情 2:小图标
	CardType int `gorm:"type:tinyint(1);default:1" json:"cardType"`
	// 卡片背景色
	BackgroundColor string `gorm:"type:varchar(50);default:''" json:"backgroundColor"`
	// 扩展参数（预留）
	ExpandParam string `gorm:"type:text" json:"-"`
}

func (m *ItemIcon) DeleteByItemIconGroupIds(db *gorm.DB, userId uint, itemIconGroupIds []uint) (err error) {
	err = db.Delete(&ItemIcon{}, "item_icon_group_id in ? AND user_id=?", itemIconGroupIds, userId).Error
	return
}

func (m *ItemIcon) DeleteByUserId(db *gorm.DB, userId uint) (err error) {
	return db.Delete(&ItemIcon{}, "user_id=?", userId).Error
}
