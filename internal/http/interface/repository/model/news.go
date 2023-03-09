package model

type News struct {
	NewsID int    `db:"news_id"`
	Title  string `db:"title"`
	Body   string `db:"body"`
}
