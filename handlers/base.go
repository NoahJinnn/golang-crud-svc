package handlers

import "gorm.io/gorm"

type HttpHandler struct {
	DB *gorm.DB
}

func NewHttpHandler(DB *gorm.DB) *HttpHandler {
	return &HttpHandler{
		DB: DB,
	}
}
