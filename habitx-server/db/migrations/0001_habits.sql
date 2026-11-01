-- +goose Up
CREATE TABLE habits (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name        text NOT NULL,
  color       text NOT NULL,
  created_at  timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE habits;
