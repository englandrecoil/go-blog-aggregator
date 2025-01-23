-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows(id, created_at, updated_at , user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5) 
    RETURNING *
)
SELECT 
    inserted_feed_follow.*,
    users.name AS user_name,
    feeds.name AS feed_name
FROM inserted_feed_follow
INNER JOIN users on users.id = inserted_feed_follow.user_id
INNER JOIN feeds on feeds.id = inserted_feed_follow.feed_id;

-- name: GetFeedByURL :one
SELECT * FROM feeds
    WHERE url = $1;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*,
    users.name AS user_name,
    feeds.name AS feed_name 
FROM feed_follows
INNER JOIN users on users.id = feed_follows.user_id
INNER JOIN feeds on feeds.id = feed_follows.feed_id
WHERE users.name = $1;