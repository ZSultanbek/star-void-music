-- name: CreateArtist :one
INSERT INTO artists (
  name,
  slug
) VALUES (
  $1, $2
)
RETURNING id, name, slug, created_at;

-- name: GetArtistByID :one
SELECT id, name, slug, created_at
FROM artists
WHERE id = $1
LIMIT 1;

-- name: ListArtists :many
SELECT id, name, slug, created_at
FROM artists
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateArtist :one
UPDATE artists
SET
  name = $2,
  slug = $3
WHERE id = $1
RETURNING id, name, slug, created_at;

-- name: DeleteArtist :exec
DELETE FROM artists
WHERE id = $1;

-- name: ListArtistAlbums :many
SELECT a.id, a.title, a.artist_id, a.cover_image_url, a.release_date, a.created_at
FROM albums a
WHERE a.artist_id = $1
ORDER BY a.release_date DESC NULLS LAST, a.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListArtistSongs :many
SELECT s.id, s.title, s.album_id, s.filepath, s.duration, s.uploaded_by, s.created_at
FROM songs s
JOIN albums a ON a.id = s.album_id
WHERE a.artist_id = $1
ORDER BY s.created_at DESC
LIMIT $2 OFFSET $3;
