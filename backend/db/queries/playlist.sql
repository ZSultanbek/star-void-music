-- name: CreatePlaylist :one
INSERT INTO playlists (
  user_id,
  name
) VALUES (
  $1, $2
)
RETURNING id, user_id, name, created_at;

-- name: GetPlaylistByID :one
SELECT id, user_id, name, created_at
FROM playlists
WHERE id = $1
LIMIT 1;

-- name: ListPlaylistsByUserID :many
SELECT id, user_id, name, created_at
FROM playlists
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePlaylistName :one
UPDATE playlists
SET name = $2
WHERE id = $1
RETURNING id, user_id, name, created_at;

-- name: DeletePlaylist :exec
DELETE FROM playlists
WHERE id = $1;
