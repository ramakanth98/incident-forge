package models

import "time"

type Incident struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Services  []string  `json:"services"`
}

type EvidenceType string

const (
	EvidenceLog    EvidenceType = "log"
	EvidenceMetric EvidenceType = "metric"
	EvidenceChange EvidenceType = "change"
)

type Evidence struct {
	ID        string         `json:"id"`
	Type      EvidenceType   `json:"type"`
	Service   string         `json:"service"`
	Timestamp time.Time      `json:"timestamp"`
	Summary   string         `json:"summary"`
	Raw       map[string]any `json:"raw"`
}

type Finding struct {
	Agent       string   `json:"agent"`
	Title       string   `json:"title"`
	Detail      string   `json:"detail"`
	EvidenceIDs []string `json:"evidence_ids"`
	Confidence  float64  `json:"confidence"` // 0..1
}
