package repository

type newsRepository struct {
	SqlHandler
}

func NewNewsRepository(handler SqlHandler) *newsRepository {
	return &newsRepository{handler}
}
