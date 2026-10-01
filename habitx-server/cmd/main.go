package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Habit struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

var habits = []Habit{
	{Id: "habit1", Name: "Habit 1", Color: "Red"},
	{Id: "habit2", Name: "Habit 2", Color: "Blue"},
	{Id: "habit3", Name: "Habit 3", Color: "Yellow"},
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{Status: "ok"}); err != nil {
		http.Error(w, "encoding fail", http.StatusInternalServerError)
		return
	}
}

func listHabits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(habits); err != nil {
		http.Error(w, "encoding fail", http.StatusInternalServerError)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz)
	mux.HandleFunc("GET /habits", listHabits)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
