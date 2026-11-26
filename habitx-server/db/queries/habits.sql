-- name: ListHabits :many
SELECT id, name, color, created_at
FROM habits
ORDER BY created_at;

-- name: GetHabit :one
SELECT *
FROM habits
WHERE id = $1 LIMIT 1;

-- name: CreateHabit :one
INSERT INTO habits (name, color)
VALUES ($1, $2)
RETURNING id, name, color, created_at;

-- name: UpdateHabit :one
UPDATE habits
SET
  name = $2,
  color = $3
WHERE id = $1
RETURNING *;

-- name: DeleteHabit :exec
DELETE FROM habits
WHERE id = $1;
