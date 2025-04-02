// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING id, created_at, updated_at, title, url, description, published_at, feed_id
`

type CreatePostParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	Description sql.NullString
	PublishedAt time.Time
	FeedID      uuid.UUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) error {
	_, err := q.db.ExecContext(ctx, createPost,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.PublishedAt,
		arg.FeedID,
	)
	return err
}

const getPostByUrl = `-- name: GetPostByUrl :one
SELECT id, created_at, updated_at, title, url, description, published_at, feed_id FROM posts WHERE url = $1
`

func (q *Queries) GetPostByUrl(ctx context.Context, url string) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostByUrl, url)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.FeedID,
	)
	return i, err
}

const getPostForUser = `-- name: GetPostForUser :many
SELECT posts.id, posts.created_at, posts.updated_at, title, posts.url, description, published_at, posts.feed_id, feed_follows.id, feed_follows.created_at, feed_follows.updated_at, feed_follows.feed_id, feed_follows.user_id, feeds.id, feeds.created_at, feeds.updated_at, name, feeds.url, feeds.user_id, last_fetched_at FROM posts
INNER JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
INNER JOIN feeds ON posts.feed_id = feeds.id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2
`

type GetPostForUserParams struct {
	UserID uuid.UUID
	Limit  int32
}

type GetPostForUserRow struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Title         string
	Url           string
	Description   sql.NullString
	PublishedAt   time.Time
	FeedID        uuid.UUID
	ID_2          uuid.UUID
	CreatedAt_2   time.Time
	UpdatedAt_2   time.Time
	FeedID_2      uuid.UUID
	UserID        uuid.UUID
	ID_3          uuid.UUID
	CreatedAt_3   time.Time
	UpdatedAt_3   time.Time
	Name          string
	Url_2         string
	UserID_2      uuid.UUID
	LastFetchedAt sql.NullTime
}

func (q *Queries) GetPostForUser(ctx context.Context, arg GetPostForUserParams) ([]GetPostForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostForUser, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostForUserRow
	for rows.Next() {
		var i GetPostForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
			&i.ID_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
			&i.FeedID_2,
			&i.UserID,
			&i.ID_3,
			&i.CreatedAt_3,
			&i.UpdatedAt_3,
			&i.Name,
			&i.Url_2,
			&i.UserID_2,
			&i.LastFetchedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
