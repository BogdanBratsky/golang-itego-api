package models

import "time"

type Article struct {
	ArticleId    int       `json:"id"`
	ArticleTitle string    `json:"title"`
	ArticleText  string    `json:"text"`
	ArticleViews int       `json:"views"`
	AuthorId     int       `json:"author"`
	CategoryId   int       `json:"category"`
	CreatedAt    time.Time `json:"createdAt"`
}
