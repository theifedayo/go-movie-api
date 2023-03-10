package models

type Movie struct {
	ID            uint     `json:"id" gorm:"primary_key"`
	Title         string   `json:"title"`
	EpisodeID     uint     `json:"episode_id"`
	OpeningCrawl  string   `json:"opening_crawl"`
	URL           string   `json:"url"`
	ReleaseDate   string   `json:"release_date"`
	CharactersURL []string `json:"characters"`
}
