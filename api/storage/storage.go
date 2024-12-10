package storage

import (
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/tangthinker/cloud-core/pkg/storage"
	"mime"
)

type BaseResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Api struct {
	baseStorage storage.Storage
}

func NewApi(rootPath string) *Api {
	return &Api{
		baseStorage: storage.NewCommonStorage(rootPath),
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
	path := ctx.Query("path")

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
