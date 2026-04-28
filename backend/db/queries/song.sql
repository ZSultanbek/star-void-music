-- name: CreateSong :one
INSERT INTO songs (
  title,
  album_id,
  filepath,
  duration,
  uploaded_by
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, title, album_id, filepath, duration, uploaded_by, created_at;

-- name: GetSongByID :one
SELECT id, title, album_id, filepath, duration, uploaded_by, created_at
FROM songs
WHERE id = $1
LIMIT 1;

-- name: ListSongs :many
SELECT id, title, album_id, filepath, duration, uploaded_by, created_at
FROM songs
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListSongsByAlbumID :many
SELECT id, title, album_id, filepath, duration, uploaded_by, created_at
FROM songs
WHERE album_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateSong :one
UPDATE songs
SET
  title = $2,
  album_id = $3,
  filepath = $4,
  duration = $5
WHERE id = $1
RETURNING id, title, album_id, filepath, duration, uploaded_by, created_at;

-- name: DeleteSong :exec
DELETE FROM songs
WHERE id = $1;
