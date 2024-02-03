package handler

import (
	"gorm.io/gorm"
)

type HandlerContext struct {
	DB *gorm.DB
}

func NewHandlerContext(db *gorm.DB) *HandlerContext {
	return &HandlerContext{
		DB: db,
	}
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
