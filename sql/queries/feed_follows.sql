-- name: CreateFeedFollow :one
WITH inserted AS (INSERT INTO feed_follows (id, created_at, updated_at, feed_id, user_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)

SELECT
    inserted.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted
INNER JOIN users ON users.id = inserted.user_id
INNER JOIN feeds ON feeds.id = inserted.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT * FROM feed_follows 
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;