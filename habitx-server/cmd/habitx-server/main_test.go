package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
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
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			newMux().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestListHabits(t *testing.T) {
	req := httptest.NewRequest("GET", "/habits", nil)
	rec := httptest.NewRecorder()
	newMux().ServeHTTP(rec, req)

	var got []Habit
	if err := json.NewDecoder(rec.Result().Body).Decode(&got); err != nil {
		t.Fatalf("decode %v", err)
	}
	if len(got) != 3 || got[0].ID != "habit1" {
		t.Errorf("got %d habits, first habitID %q; want 3 habits; first ID habit1", len(got), got[0].ID)
	}
}
