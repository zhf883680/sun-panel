package models

type File struct {
	BaseModel
	Src      string `json:"src"`
	UserId   uint   `json:"userId"`
	FileName string `json:"fileName" gorm:"varchar(255)"` // 文件名
	Method   int    `gorm:"int(5)" json:"method"`         // 上传方式
	Ext      string `gorm:"varchar(255)" json:"ext"`      // 扩展名
	// 文件类型（预留）
	Type int `gorm:"int(5);default:0" json:"type"`
	// 是否公开到公共图库
	IsPublicGallery bool `gorm:"default:false" json:"isPublicGallery"`
	// 图库分类 icon / wallpaper / unknown
	Category string `gorm:"type:varchar(50);index" json:"category"`
	// 文件大小（字节）
	Size int64 `gorm:"int(20);default:0" json:"size"`
}

// 添加一个文件记录
func (m *File) AddFile(userId uint, fileName, ext, src string, size int64, category string) (File, error) {
	file := File{
		UserId:   userId,
		FileName: fileName,
		Src:      src,
		Ext:      ext,
		Size:     size,
		Category: category,
		Type:     CategoryToType(category),
	}
	err := Db.Create(&file).Error

	return file, err
}

// 图库分类与 v1.8.1 的 file.type 映射：0未知 1壁纸 2图标
func CategoryToType(category string) int {
	switch category {
	case "wallpaper":
		return 1
	case "icon":
		return 2
	default:
		return 0
	}
}
