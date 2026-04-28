-- name: AddSongToUserLibrary :one
INSERT INTO user_library (
  user_id,
  song_id
) VALUES (
  $1, $2
)
RETURNING user_id, song_id, added_at;

-- name: ListUserLibrarySongs :many
SELECT
  ul.user_id,
  ul.song_id,
  ul.added_at,
  s.title,
  s.album_id,
  s.filepath,
  s.duration,
  s.uploaded_by,
  s.created_at
FROM user_library ul
JOIN songs s ON s.id = ul.song_id
WHERE ul.user_id = $1
ORDER BY ul.added_at DESC
LIMIT $2 OFFSET $3;

-- name: RemoveSongFromUserLibrary :exec
DELETE FROM user_library
WHERE user_id = $1 AND song_id = $2;
