package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/heniums/habitx-server/db"
	"github.com/heniums/habitx-server/internal/store"
)

type habitJSON struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	CreatedAt string `json:"created_at"`
}

type server struct {
	queries *db.Queries
}

func (s *server) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{Status: "ok"}); err != nil {
		http.Error(w, "encoding fail", http.StatusInternalServerError)
		return
	}
}

func (s *server) listHabits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	habits, err := s.queries.ListHabits(r.Context())
	if err != nil {
		http.Error(w, "fail to list habits", http.StatusInternalServerError)
		return
	}

	if habits == nil {
		habits = []db.Habit{}
	}

	out := make([]habitJSON, 0, len(habits))
	for _, h := range habits {
		out = append(out, habitJSON{
			ID:        h.ID.String(),
			Name:      h.Name,
			Color:     h.Color,
			CreatedAt: h.CreatedAt.Time.Format(time.RFC3339),
		})
	}

	if err := json.NewEncoder(w).Encode(out); err != nil {
		http.Error(w, "encoding fail", http.StatusInternalServerError)
		return
	}
}

func (s *server) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", s.healthz)
	mux.HandleFunc("GET /habits", s.listHabits)

	return mux
}

func main() {
	queries, pool, err := store.New(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	s := &server{queries: queries}

	log.Fatal(http.ListenAndServe(":8080", s.routes()))
}
