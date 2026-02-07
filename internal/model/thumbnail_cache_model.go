package model

import (
	"errors"
	"fmt"

	"github.com/tangthinker/cloud-core/internal/model/schema"

	"github.com/tangthinker/cloud-core/internal/db"
	"gorm.io/gorm"
)

type ThumbnailCacheModel struct {
	DB *gorm.DB
}

func NewThumbnailCacheModel() *ThumbnailCacheModel {
	d := db.GetDB()
	if err := d.AutoMigrate(&schema.ThumbnailCache{}); err != nil {
		panicStr := fmt.Sprintf("auto migrate thumbnail cache table failed: %v", err)
		panic(panicStr)
	}

	return &ThumbnailCacheModel{DB: d}
}

func (m *ThumbnailCacheModel) GetByHash(fileHash string) (*schema.ThumbnailCache, error) {
	var cache schema.ThumbnailCache
	if err := m.DB.Where("file_hash = ?", fileHash).First(&cache).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cache, nil
}

func (m *ThumbnailCacheModel) Create(req *schema.ThumbnailCache) error {
	return m.DB.Create(&schema.ThumbnailCache{
		FileHash:        req.FileHash,
		Base64Thumbnail: req.Base64Thumbnail,
	}).Error
}

func (m *ThumbnailCacheModel) Update(fileHash string, base64Thumbnail string) error {
	return m.DB.Model(&schema.ThumbnailCache{}).Where("file_hash = ?", fileHash).Update("base64_thumbnail", base64Thumbnail).Error
}

func (m *ThumbnailCacheModel) Delete(fileHash string) error {
	return m.DB.Where("file_hash = ?", fileHash).Delete(&schema.ThumbnailCache{}).Error
}
