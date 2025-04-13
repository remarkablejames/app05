package repo_impl

import (
	"app05/internal/core/domain/dtos/postDTOs"
	"context"
	"database/sql"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

// GetAllPosts retrieves all posts from the database
func (r *PostRepository) GetAllPosts(ctx context.Context) ([]*postDTOs.PostDTO, error) {
	posts := []*postDTOs.PostDTO{}

	query := `SELECT id, user_id, title, content, excerpt, status, slug, 
                   view_count, published_at, created_at, updated_at 
              FROM posts`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post postDTOs.PostDTO
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.Excerpt,
			&post.Status,
			&post.Slug,
			&post.ViewCount,
			&post.PublishedAt,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
