-- name: ListHabits :many
SELECT id, name, color, created_at
FROM habits
ORDER BY created_at;

-- name: CreateHabit :one
INSERT INTO habits (name, color)
VALUES ($1, $2)
RETURNING id, name, color, created_at;
