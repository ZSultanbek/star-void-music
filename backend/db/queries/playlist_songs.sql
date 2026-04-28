-- name: AddSongToPlaylist :one
INSERT INTO playlist_songs (
  playlist_id,
  song_id,
  position
) VALUES (
  $1, $2, $3
)
RETURNING playlist_id, song_id, position;

-- name: ListPlaylistSongs :many
SELECT
  ps.playlist_id,
  ps.song_id,
  ps.position,
  s.title,
  s.album_id,
  s.filepath,
  s.duration,
  s.uploaded_by,
  s.created_at
FROM playlist_songs ps
JOIN songs s ON s.id = ps.song_id
WHERE ps.playlist_id = $1
ORDER BY ps.position ASC
LIMIT $2 OFFSET $3;

-- name: UpdatePlaylistSongPosition :one
UPDATE playlist_songs
SET position = $3
WHERE playlist_id = $1 AND song_id = $2
RETURNING playlist_id, song_id, position;

-- name: RemoveSongFromPlaylist :exec
DELETE FROM playlist_songs
WHERE playlist_id = $1 AND song_id = $2;

-- name: ClearPlaylistSongs :exec
DELETE FROM playlist_songs
WHERE playlist_id = $1;
