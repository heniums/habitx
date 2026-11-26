package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/heniums/habitx-server/db"
	"github.com/heniums/habitx-server/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

func (s *server) getHabit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var id pgtype.UUID
	if err := id.Scan(r.PathValue("id")); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	habit, err := s.queries.GetHabit(r.Context(), id)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "fail to get habit", http.StatusInternalServerError)
		return
	}

	habitJSON := habitJSON{
		ID:        habit.ID.String(),
		Name:      habit.Name,
		Color:     habit.Color,
		CreatedAt: habit.CreatedAt.Time.Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(habitJSON); err != nil {
		http.Error(w, "encoding fail", http.StatusInternalServerError)
		return
	}
}

func (s *server) updateHabit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var id pgtype.UUID
	if err := id.Scan(r.PathValue("id")); err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Decode request body
	var request struct {
		Name string `json:"name"`
		Color string `json:"color"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Body validation
	if request.Name == "" || request.Color == "" {
		http.Error(w, "name and color required", http.StatusBadRequest)
		return
	}

	updateParams := db.UpdateHabitParams {
		ID: id,
		Name: request.Name,
		Color: request.Color,
	}

	habit, err := s.queries.UpdateHabit(r.Context(), updateParams)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "fail to update habit", http.StatusInternalServerError)
		return
	}

	habitJSON := habitJSON{
		ID:        habit.ID.String(),
		Name:      habit.Name,
		Color:     habit.Color,
		CreatedAt: habit.CreatedAt.Time.Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(habitJSON); err != nil {
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
	mux.HandleFunc("GET /habits/{id}", s.getHabit)
	mux.HandleFunc("PUT /habits/{id}", s.updateHabit)

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
