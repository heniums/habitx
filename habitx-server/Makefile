# Server
run:
	HABITX_DB_DSN="postgres://habitx:habitx@localhost:5433/habitx?sslmode=disable" go run ./cmd/habitx-server

test:
	go test ./...

tidy:
	go mod tidy

.PHONY: run test tidy
