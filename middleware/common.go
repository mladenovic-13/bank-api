package middleware

import "gorm.io/gorm"

type MiddlewareContext struct {
	DB *gorm.DB
}

func NewMiddlewareContext(db *gorm.DB) *MiddlewareContext {
	return &MiddlewareContext{
		DB: db,
	}
}
