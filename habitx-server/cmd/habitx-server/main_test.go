package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"
	"uuid"

	"github.com/heniums/habitx-server/internal/store"
)

func TestRoutes(t *testing.T) {
	queries, cleanup := store.TestDB(t)
	defer cleanup()
	s := &server{queries: queries}

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{name: "healthz responds 200", method: "GET", path: "/healthz", wantStatus: 200},
		{name: "habits responds 200", method: "GET", path: "/habits", wantStatus: 200},
		{name: "wrong method rejected", method: "POST", path: "/healthz", wantStatus: 405},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.path, nil)

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
