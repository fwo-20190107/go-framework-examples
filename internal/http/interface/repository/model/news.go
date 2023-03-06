package model

type News struct {
	NewsID int    `json:"news_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
