package router

import (
	"github.com/mladenovic-13/bank-api/handler"
	"github.com/mladenovic-13/bank-api/middleware"
	"gorm.io/gorm"
)

type RouterContext struct {
	*gorm.DB
}

func NewRouterContext(db *gorm.DB) *RouterContext {
	return &RouterContext{
		DB: db,
	}
}

type RouteContext struct {
	*handler.HandlerContext
	*middleware.MiddlewareContext
}

func NewRouteContext(c *RouterContext) *RouteContext {
	return &RouteContext{
		HandlerContext:    handler.NewHandlerContext(c.DB),
		MiddlewareContext: middleware.NewMiddlewareContext(c.DB),
	}
}
