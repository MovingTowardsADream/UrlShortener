package entity

import "time"

type Url struct {
	Url      string        `json:"url"`
	ShortUrl string        `json:"short_url"`
	Expiry   time.Duration `json:"expiry"`
}
