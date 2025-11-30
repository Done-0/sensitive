// Package main demonstrates content moderation API for production web services
// Creator: Done-0
// Created: 2025-01-15
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Done-0/sensitive"
)

var detector *sensitive.Detector

type CheckRequest struct {
	Content string `json:"content"`
}

type CheckResponse struct {
	Safe     bool     `json:"safe"`
	Reason   string   `json:"reason,omitempty"`
	Matches  []string `json:"matches,omitempty"`
	Filtered string   `json:"filtered,omitempty"`
}

func init() {
	detector = sensitive.New(
		sensitive.WithFilterStrategy(sensitive.StrategyMask),
	)

	words := map[string]sensitive.Level{
		"illegal":  sensitive.LevelHigh,
		"violence": sensitive.LevelHigh,
		"abuse":    sensitive.LevelMedium,
		"spam":     sensitive.LevelLow,
	}
	detector.AddWords(words)
	detector.Build()

	log.Println("Detector initialized with", len(words), "words")
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result := detector.Detect(req.Content)

	resp := CheckResponse{
		Safe:     !result.HasSensitive,
		Filtered: result.FilteredText,
	}

	if result.HasSensitive {
		resp.Reason = "Content contains sensitive words"
		resp.Matches = make([]string, len(result.Matches))
		for i, m := range result.Matches {
			resp.Matches[i] = m.Word
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/api/check", checkHandler)

	addr := ":8080"
	log.Printf("Content moderation API server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
