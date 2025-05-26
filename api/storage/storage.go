package storage

import (
	"encoding/base64"
	"fmt"
	"mime"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/tangthinker/cloud-core/internal/service"
	"github.com/tangthinker/cloud-core/pkg/storage"
	"github.com/tangthinker/cloud-core/pkg/video_trans"
)

type BaseResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Api struct {
	baseStorage      storage.Storage
	transService     video_trans.Service
	thumbnailService *service.ThumbnailService
}

func NewApi(rootPath string) *Api {
	return &Api{
		baseStorage:      storage.NewCommonStorage(rootPath),
		transService:     video_trans.NewService(rootPath),
		thumbnailService: service.NewThumbnailService(rootPath),
	}
}

func (a *Api) LS(ctx *fiber.Ctx) error {
	var req LSReq
	if ctx.BodyParser(&req) != nil {
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}

	ls, err := a.baseStorage.LS(req.Path)
	if err != nil {
		return ctx.JSON(BaseResp{
			Code: 1,
			Msg:  "ls failed: " + err.Error(),
		})
	}

	return ctx.JSON(BaseResp{
		Code: 0,
		Msg:  "success",
		Data: ls,
	})
}

func (a *Api) Get(ctx *fiber.Ctx) error {
	path := ctx.Query("filepath")

	utils, err := a.baseStorage.Get(path)
	if err != nil {
		return ctx.JSON(BaseResp{
			Code: 1,
			Msg:  "get failed: " + err.Error(),
		})
	}

	// base64 编码
	var base64Data = make([]byte, base64.URLEncoding.EncodedLen(len(utils)))
	base64.StdEncoding.Encode(base64Data, utils)

	return ctx.JSON(BaseResp{
		Code: 0,
		Msg:  "success",
		Data: string(base64Data),
	})
}

func (a *Api) GetThumbnail(ctx *fiber.Ctx) error {
	path := ctx.Query("filepath")
	width := ctx.Query("width", "200")
	height := ctx.Query("height", "200")

	widthNum, err := strconv.Atoi(width)
	if err != nil {
		widthNum = 200
	}

	heightNum, err := strconv.Atoi(height)
	if err != nil {
		heightNum = 200
	}

	thumbnail, err := a.thumbnailService.GetThumbnail(path, widthNum, heightNum)
	if err != nil {
		return ctx.JSON(BaseResp{
			Code: 1,
			Msg:  "get thumbnail failed: " + err.Error(),
		})
	}

	return ctx.JSON(BaseResp{
		Code: 0,
		Msg:  "success",
		Data: thumbnail,
	})
}

func (a *Api) Stat(ctx *fiber.Ctx) error {
	var req StatReq
	if ctx.BodyParser(&req) != nil {
		ctx.Status(fiber.StatusBadRequest)
		return nil
	}

	stat, err := a.baseStorage.Stat(req.Path)
	if err != nil {
		return ctx.JSON(BaseResp{
			Code: 1,
			Msg:  "stat failed: " + err.Error(),
		})
	}

	return ctx.JSON(BaseResp{
		Code: 0,
		Msg:  "success",
		Data: stat,
	})
}

func (a *Api) Download(ctx *fiber.Ctx) error {
	filepath := ctx.Query("filepath")

	stat, err := a.baseStorage.Stat(filepath)
	if err != nil {
		return ctx.JSON(BaseResp{
			Code: 1,
			Msg:  "stat failed: " + err.Error(),
		})
	}

	// 设置响应头
	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", stat.Filename))
	mimeType := mime.TypeByExtension(".txt") // 根据扩展名获取 MIME 类型
	ctx.Set("Content-Type", mimeType)
	ctx.Set("Content-Length", fmt.Sprintf("%d", stat.Size))

	writer := ctx.Response().BodyWriter()

	err = a.baseStorage.Download(filepath, writer)
	if err != nil {
		return ctx.JSON(BaseResp{
			Code: 1,
			Msg:  "download failed: " + err.Error(),
		})
	}

	return nil

}

func (a *Api) Upload(ctx *fiber.Ctx) error {
	filename := ctx.FormValue("filename")
	path := ctx.FormValue("path")

	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(BaseResp{
			Code: 1,
			Msg:  "upload failed: " + err.Error(),
		})
	}

	fileReader, err := file.Open()
	if err != nil {
		return ctx.JSON(BaseResp{
			Code: 1,
			Msg:  "upload failed: " + err.Error(),
		})
	}

	err = a.baseStorage.Upload(path, filename, fileReader)
	if err != nil {
		return ctx.JSON(BaseResp{
			Code: 1,
			Msg:  "upload failed: " + err.Error(),
		})
	}

	return ctx.JSON(BaseResp{
		Code: 0,
		Msg:  "success",
	})
}

func (a *Api) Trans2M3U8(ctx *fiber.Ctx) error {
	filepath := ctx.Query("filepath")

	state, err := a.transService.TransState(filepath)
	if err != nil {
		return ctx.JSON(BaseResp{
			Code: 1,
			Msg:  "access trans state failed: " + err.Error(),
		})
	}

	return ctx.JSON(BaseResp{
		Code: 0,
		Msg:  "success",
		Data: state,
	})

}
