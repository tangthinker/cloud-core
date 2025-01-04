package middleware

import (
	"crypto/sha256"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/tangthinker/cloud-core/internal/encrypt"
)

func ContentEncrypt(ctx *fiber.Ctx) error {
	err := ctx.Next()
	if err != nil {
		return err
	}

	token := ctx.Locals("token").(string)

	if token == "" {
		ctx.Status(fiber.StatusForbidden)
		return ctx.SendString("Forbidden: Invalid Token")
	}

	uid := ctx.Locals("uid").(string)

	if uid == "" {
		ctx.Status(fiber.StatusForbidden)
		return ctx.SendString("Forbidden: Invalid UID")
	}

	key := generate32Key(token)

	response := ctx.Response().Body()

	encryptedBody, err := encrypt.Encrypt(response, []byte(key))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Encryption failed")
	}

	ctx.Response().SetBody([]byte(encryptedBody))

	println("ContentEncrypt:", encryptedBody)

	return ctx.SendStatus(fiber.StatusOK)
}

func generate32Key(source string) string {
	hash := sha256.Sum256([]byte(source)) // SHA-256 会生成一个 32 字节的哈希值
	fmt.Printf("SHA-256 Hash: %x\n", hash)
	return string(hash[:]) // 将哈希值作为密钥返回
}
