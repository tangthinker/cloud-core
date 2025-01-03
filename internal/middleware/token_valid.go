package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/tangthinker/user-center/pkg"
)

func TokenValid(ctx *fiber.Ctx) error {
	headers := ctx.GetReqHeaders()
	if len(headers) == 0 {
		ctx.Status(fiber.StatusForbidden)
		return ctx.SendString("Forbidden: Invalid Token")
	}
	authorization := headers["Authorization"]
	if len(authorization) == 0 {
		ctx.Status(fiber.StatusForbidden)
		return ctx.SendString("Forbidden: Invalid Token")
	}
	uid, err := pkg.TokenValid(authorization[0])
	if err != nil {
		ctx.Status(fiber.StatusForbidden)
		return ctx.SendString("Forbidden: Invalid Token")
	}

	fmt.Println("authorization successful", uid)

	ctx.Set("uid", uid)

	return ctx.Next()

}
