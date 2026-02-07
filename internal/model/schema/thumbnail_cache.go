package schema

import "gorm.io/gorm"

type ThumbnailCache struct {
	gorm.Model
	FileHash        string `gorm:"type:varchar(255);not null;uniqueIndex" json:"file_hash"`
	Base64Thumbnail string `gorm:"type:text;null" json:"base64_thumbnail"`
}

func (ThumbnailCache) TableName() string {
	return "tb_thumbnail_cache"
}
