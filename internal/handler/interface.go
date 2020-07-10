package handler

import "net/http"

type HandlerInterface interface {
	GetIndex() http.HandlerFunc
	GetAll() http.HandlerFunc
}
