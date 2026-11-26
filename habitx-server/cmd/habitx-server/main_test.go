package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"uuid"

	"github.com/heniums/habitx-server/db"
	"github.com/heniums/habitx-server/internal/store"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestRoutes(t *testing.T) {
	queries, cleanup := store.TestDB(t)
	defer cleanup()
	s := &server{queries: queries}

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{name: "healthz responds 200", method: "GET", path: "/healthz", body: "", wantStatus: 200},
		{name: "habits responds 200", method: "GET", path: "/habits", body: "", wantStatus: 200},
		{name: "wrong method rejected", method: "POST", path: "/healthz", body: "", wantStatus: 405},
		{name: "post habit without body", method: "POST", path: "/habits", body: "", wantStatus: 400},
		{name: "post habit invalid json", method: "POST", path: "/habits", body: "not json", wantStatus: 400},
		{name: "post habit empty name", method: "POST", path: "/habits", body: `{"name":"","color":"red"}`, wantStatus: 400},
		{name: "get habit bad uuid", method: "GET", path: "/habits/not-a-uuid", body: "", wantStatus: 400},
		{name: "get habit missing id", method: "GET", path: "/habits/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", body: "", wantStatus: 404},
		{name: "put habit bad uuid", method: "PUT", path: "/habits/not-a-uuid", body: `{"name":"x","color":"y"}`, wantStatus: 400},
		{name: "put habit missing id", method: "PUT", path: "/habits/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", body: `{"name":"x","color":"y"}`, wantStatus: 404},
		{name: "put habit empty name", method: "PUT", path: "/habits/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", body: `{"name":"","color":"y"}`, wantStatus: 400},
		{name: "delete habit bad uuid", method: "DELETE", path: "/habits/not-a-uuid", body: "", wantStatus: 400},
		{name: "delete habit missing id is idempotent", method: "DELETE", path: "/habits/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", body: "", wantStatus: 204},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))

			s.routes().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestListHabits(t *testing.T) {
	queries, cleanup := store.TestDB(t)
	defer cleanup()
	s := &server{queries: queries}

	seedParams := []struct {
		Name  string
		Color string
	}{
		{Name: "eat clean", Color: "red"},
		{Name: "drink enough water", Color: "blue"},
		{Name: "exercise", Color: "yellow"},
	}

	for _, seedParam := range seedParams {
		if _, err := queries.CreateHabit(t.Context(), seedParam); err != nil {
			t.Fatalf("seed %v", err)
		}
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/habits", nil)

	s.routes().ServeHTTP(rec, req)

	var got []habitJSON

	if err := json.NewDecoder(rec.Result().Body).Decode(&got); err != nil {
		t.Fatalf("decode %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("got %d habits, want 3 habits", len(got))
	}
	if _, err := uuid.Parse(got[0].ID); err != nil {
		t.Errorf("ID is not a UUID: %v", err)
	}
	if _, err := time.Parse(time.RFC3339, got[0].CreatedAt); err != nil {
		t.Errorf("created_at RFC3339: %v", err)
	}
	if got[0].Name != "eat clean" || got[1].Name != "drink enough water" || got[2].Name != "exercise" {
		t.Errorf("order mismatch: %s, %s, %s", got[0].Name, got[1].Name, got[2].Name)
	}
}

// seedHabit inserts one habit through the real query and fails the test on error.
func seedHabit(t *testing.T, queries *db.Queries, name, color string) db.Habit {
	t.Helper()

	habit, err := queries.CreateHabit(t.Context(), db.CreateHabitParams{Name: name, Color: color})
	if err != nil {
		t.Fatalf("seed habit %q: %v", name, err)
	}
	return habit
}

func TestCreateHabit(t *testing.T) {
	queries, cleanup := store.TestDB(t)
	defer cleanup()
	s := &server{queries: queries}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/habits", strings.NewReader(`{"name":"stretch","color":"blue"}`))

	s.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("got status %d, want 201", rec.Code)
	}

	var got habitJSON
	if err := json.NewDecoder(rec.Result().Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Name != "stretch" || got.Color != "blue" {
		t.Errorf("roundtrip mismatch: got name=%q color=%q", got.Name, got.Color)
	}
	if _, err := uuid.Parse(got.ID); err != nil {
		t.Errorf("ID is not a UUID: %v", err)
	}

	// the API said 201; prove the row actually persisted
	var id pgtype.UUID
	if err := id.Scan(got.ID); err != nil {
		t.Fatalf("response id %q is not a valid UUID: %v", got.ID, err)
	}
	if _, err := queries.GetHabit(t.Context(), id); err != nil {
		t.Errorf("created habit not readable via GetHabit: %v", err)
	}
}

func TestUpdateHabit(t *testing.T) {
	queries, cleanup := store.TestDB(t)
	defer cleanup()
	s := &server{queries: queries}

	seed := seedHabit(t, queries, "old name", "red")

	rec := httptest.NewRecorder()
	path := "/habits/" + seed.ID.String()
	req := httptest.NewRequest("PUT", path, strings.NewReader(`{"name":"new name","color":"green"}`))

	s.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}

	var got habitJSON
	if err := json.NewDecoder(rec.Result().Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Name != "new name" || got.Color != "green" {
		t.Errorf("response shows old values: got name=%q color=%q", got.Name, got.Color)
	}

	// assert the database changed, not just the response body
	stored, err := queries.GetHabit(t.Context(), seed.ID)
	if err != nil {
		t.Fatalf("get after update: %v", err)
	}
	if stored.Name != "new name" || stored.Color != "green" {
		t.Errorf("DB did not change: got name=%q color=%q", stored.Name, stored.Color)
	}
}

func TestDeleteHabit(t *testing.T) {
	queries, cleanup := store.TestDB(t)
	defer cleanup()
	s := &server{queries: queries}

	seed := seedHabit(t, queries, "doomed", "gray")
	path := "/habits/" + seed.ID.String()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", path, nil)
	s.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("got status %d, want 204", rec.Code)
	}

	// deletion is proven by absence
	rec = httptest.NewRecorder()
	req = httptest.NewRequest("GET", path, nil)
	s.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("GET after DELETE got status %d, want 404", rec.Code)
	}

	// idempotence: deleting the same id again still succeeds
	rec = httptest.NewRecorder()
	req = httptest.NewRequest("DELETE", path, nil)
	s.routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("second DELETE got status %d, want 204", rec.Code)
	}
}
