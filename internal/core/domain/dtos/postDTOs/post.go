package postDTOs

import "time"

type PostDTO struct {
	ID          int        `json:"id"`
	UserID      string     `json:"user_id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Excerpt     *string    `json:"excerpt,omitempty"`
	Status      string     `json:"status"`
	ViewCount   int        `json:"view_count"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	Slug        *string    `json:"slug,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
