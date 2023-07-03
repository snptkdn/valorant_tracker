package dto

type WebhookRequestDto struct {
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds"`
}

type Embed struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Thumbnail   Thumbnail `json:"thumbnail"`
	Fields      []Field   `json:"fields"`
	Color       int       `json:"color"`
}

type Thumbnail struct {
	URL string `json:"url"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}
