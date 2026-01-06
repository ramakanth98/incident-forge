package connectors

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ramakanth98/incident-forge/internal/models"
)

type BundleLoader struct{}

func NewBundleLoader() *BundleLoader { return &BundleLoader{} }

func (b *BundleLoader) LoadIncident(bundlePath string) (models.Incident, error) {
	p := filepath.Join(bundlePath, "incident.json")
	raw, err := os.ReadFile(p)
	if err != nil {
		return models.Incident{}, fmt.Errorf("read incident.json: %w", err)
	}

	// Custom parse to ensure time formats are correct.
	var tmp struct {
		ID       string   `json:"id"`
		Title    string   `json:"title"`
		Start    string   `json:"start_time"`
		End      string   `json:"end_time"`
		Services []string `json:"services"`
	}
	if err := json.Unmarshal(raw, &tmp); err != nil {
		return models.Incident{}, fmt.Errorf("parse incident.json: %w", err)
	}

	start, err := time.Parse(time.RFC3339, tmp.Start)
	if err != nil {
		return models.Incident{}, fmt.Errorf("parse start_time: %w", err)
	}
	end, err := time.Parse(time.RFC3339, tmp.End)
	if err != nil {
		return models.Incident{}, fmt.Errorf("parse end_time: %w", err)
	}

	return models.Incident{
		ID:        tmp.ID,
		Title:     tmp.Title,
		StartTime: start,
		EndTime:   end,
		Services:  tmp.Services,
	}, nil
}

func (b *BundleLoader) LoadEvidence(bundlePath string) ([]models.Evidence, error) {
	p := filepath.Join(bundlePath, "evidence.json")
	raw, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("read evidence.json: %w", err)
	}

	var tmp []struct {
		ID        string              `json:"id"`
		Type      models.EvidenceType `json:"type"`
		Service   string              `json:"service"`
		Timestamp string              `json:"timestamp"`
		Summary   string              `json:"summary"`
		Raw       map[string]any      `json:"raw"`
	}

	if err := json.Unmarshal(raw, &tmp); err != nil {
		return nil, fmt.Errorf("parse evidence.json: %w", err)
	}

	out := make([]models.Evidence, 0, len(tmp))
	for _, x := range tmp {
		ts, err := time.Parse(time.RFC3339, x.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("parse evidence timestamp for %s: %w", x.ID, err)
		}
		out = append(out, models.Evidence{
			ID:        x.ID,
			Type:      x.Type,
			Service:   x.Service,
			Timestamp: ts,
			Summary:   x.Summary,
			Raw:       x.Raw,
		})
	}
	return out, nil
}
