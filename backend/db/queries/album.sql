-- name: CreateAlbum :one
INSERT INTO albums (
  title,
  artist_id,
  cover_image_url,
  release_date
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, title, artist_id, cover_image_url, release_date, created_at;

-- name: GetAlbumByID :one
SELECT id, title, artist_id, cover_image_url, release_date, created_at
FROM albums
WHERE id = $1
LIMIT 1;

-- name: ListAlbums :many
SELECT id, title, artist_id, cover_image_url, release_date, created_at
FROM albums
ORDER BY release_date DESC NULLS LAST, created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListAlbumsByArtistID :many
SELECT id, title, artist_id, cover_image_url, release_date, created_at
FROM albums
WHERE artist_id = $1
ORDER BY release_date DESC NULLS LAST, created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateAlbum :one
UPDATE albums
SET
  title = $2,
  artist_id = $3,
  cover_image_url = $4,
  release_date = $5
WHERE id = $1
RETURNING id, title, artist_id, cover_image_url, release_date, created_at;

-- name: DeleteAlbum :exec
DELETE FROM albums
WHERE id = $1;

-- name: ListAlbumSongs :many
SELECT id, title, album_id, filepath, duration, uploaded_by, created_at
FROM songs
WHERE album_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
