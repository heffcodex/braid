package braid

import (
	"github.com/gofiber/fiber/v2"
)

type DTO struct {
	ctx *fiber.Ctx
}

func (o *DTO) Ctx() *fiber.Ctx       { return o.ctx }
func (o *DTO) SetCtx(ctx *fiber.Ctx) { o.ctx = ctx }
