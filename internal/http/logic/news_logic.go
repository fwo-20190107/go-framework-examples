package logic

import "examples/internal/http/logic/repository"

type NewsLogic interface {
}

type newsLogic struct {
	newsRepository repository.NewsRepository
}

func NewNewsLogic(newsRepository repository.NewsRepository) *newsLogic {
	return &newsLogic{
		newsRepository: newsRepository,
	}
}
