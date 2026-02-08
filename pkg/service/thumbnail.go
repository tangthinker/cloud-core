package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/tangthinker/cloud-core/helper"
	"github.com/tangthinker/cloud-core/internal/model"
	"github.com/tangthinker/cloud-core/internal/model/schema"
	"github.com/tangthinker/cloud-core/pkg/storage"
)

type ThumbnailService struct {
	baseStorage    storage.Storage
	thumbnailCache *model.ThumbnailCacheModel
}

func NewThumbnailService(storageService storage.Storage) *ThumbnailService {
	return &ThumbnailService{
		baseStorage:    storageService,
		thumbnailCache: model.NewThumbnailCacheModel(),
	}
}

func (s *ThumbnailService) GetThumbnail(filepath string, imageWidth int, imageHeight int) (string, error) {
	fileStream, err := s.baseStorage.GetStream(filepath)
	if err != nil {
		return "", err
	}
	defer fileStream.Close()

	fileHash := sha256.New()
	if _, err := io.Copy(fileHash, fileStream); err != nil {
		return "", err
	}

	_, err = fileHash.Write([]byte(fmt.Sprintf("[%d,%d]", imageWidth, imageHeight)))
	if err != nil {
		return "", err
	}

	fileHashStr := fmt.Sprintf("%x", fileHash.Sum(nil))

	thumbnail, err := s.thumbnailCache.GetByHash(fileHashStr)
	if err != nil {
		return "", err
	}

	if thumbnail != nil {
		return thumbnail.Base64Thumbnail, nil
	}

	isVideo := strings.HasSuffix(strings.ToLower(filepath), ".mp4") || strings.HasSuffix(strings.ToLower(filepath), ".mov")

	var thumbnailBase64 string
	if isVideo {
		thumbnailBase64, err = s.generateVideoThumbnail(s.baseStorage.RealPath(filepath), imageWidth, imageHeight)
	} else {
		thumbnailBase64, err = s.generateThumbnail(filepath, imageWidth, imageHeight)
	}
	if err != nil {
		return "", err
	}

	thumbnail = &schema.ThumbnailCache{
		FileHash:        fileHashStr,
		Base64Thumbnail: thumbnailBase64,
	}
	if err := s.thumbnailCache.Create(thumbnail); err != nil {
		return "", err
	}

	return thumbnailBase64, nil
}

func (s *ThumbnailService) generateVideoThumbnail(filepath string, imageWidth int, imageHeight int) (string, error) {
	thumbnail, err := helper.GenerateVideoThumbnail(filepath, imageWidth, imageHeight)
	if err != nil {
		return "", err
	}

	var thuBuff bytes.Buffer
	if err := imaging.Encode(&thuBuff, thumbnail, imaging.JPEG); err != nil {
		return "", err
	}

	b := thuBuff.Bytes()

	var base64Data = make([]byte, base64.URLEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(base64Data, b)

	return string(base64Data), nil
}

func (s *ThumbnailService) generateThumbnail(filepath string, imageWidth int, imageHeight int) (string, error) {
	utils, err := s.baseStorage.Get(filepath)
	if err != nil {
		return "", err
	}

	src, err := helper.DecodeImage(utils, filepath)
	if err != nil {
		return "", err
	}

	// 生成缩略图
	thumbnail, err := helper.GenerateThumbnail(src, imageWidth, imageHeight)
	if err != nil {
		return "", err
	}

	var thuBuff bytes.Buffer
	if err := imaging.Encode(&thuBuff, thumbnail, imaging.JPEG); err != nil {
		return "", err
	}

	b := thuBuff.Bytes()

	// base64 编码
	var base64Data = make([]byte, base64.URLEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(base64Data, b)

	return string(base64Data), nil
}
