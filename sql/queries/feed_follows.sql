-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *, (
    SELECT users.name FROM users
    WHERE id = user_id
    ) AS user_name, (
    SELECT feeds.name FROM feeds
    WHERE id = feed_id
    ) AS feed_name;

-- name: GetFeedFollowsForUser :many
SELECT *, (
    SELECT users.name FROM users
    WHERE id = feed_follows.user_id
    ) AS user_name, (
    SELECT feeds.name FROM feeds
    WHERE id = feed_follows.feed_id
    ) AS feed_name FROM feed_follows
WHERE feed_follows.user_id = $1;

-- name: RemoveFollowFeed :exec
DELETE FROM feed_follows
WHERE feed_id = (
    SELECT feeds.id FROM feeds
    WHERE url = $1
) AND feed_follows.user_id = $2;