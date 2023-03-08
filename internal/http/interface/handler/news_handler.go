package handler

import "examples/internal/http/logic"

type newsHandler struct {
	newsLogic logic.NewsLogic
}

func NewNewsHandler(newsLogic logic.NewsLogic) *newsHandler {
	return &newsHandler{
		newsLogic: newsLogic,
	}
}
