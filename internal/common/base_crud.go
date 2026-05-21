package common

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/morkid/paginate"
)

type Context interface {
	Context() context.Context
}

type IBaseCrudRepo[TEntity any, TInput any, TPage any] interface {
	Create(c Context, input TInput) (uint, error)
	Read(c fiber.Ctx) (paginate.Page, []TPage)
	Update(c Context, id any, input TInput) error
	Delete(c Context, id any) error
	GetByID(c Context, id any) (TEntity, error)
}

type IBaseCrudService[TEntity any, TInput any, TPage any] interface {
	Create(c Context, input TInput) (uint, error)
	Read(c fiber.Ctx) (paginate.Page, []TPage)
	Update(c Context, id any, input TInput) error
	Delete(c Context, id any) error
	GetByID(c Context, id any) (TEntity, error)
}
