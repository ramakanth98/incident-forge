package agents

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/ramakanth98/incident-forge/internal/models"
)

type LogSignalAgent struct{}

func (a *LogSignalAgent) Name() string { return "log-signal" }

func (a *LogSignalAgent) Run(ctx context.Context, s Store) error {
	inc := s.Incident()
	ev := s.Evidence()

	keywords := []string{
		"timeout", "timed out",
		"error", "exception", "panic",
		"retry", "throttle", "429",
		"refused", "connection reset",
		"deadline exceeded",
	}

	type scored struct {
		e     models.Evidence
		score int
	}

	list := make([]scored, 0)

	for _, e := range ev {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if e.Type != models.EvidenceLog {
			continue
		}
		if len(inc.Services) > 0 && !contains(inc.Services, e.Service) {
			continue
		}

		text := strings.ToLower(e.Summary)
		if msg, ok := e.Raw["msg"].(string); ok && msg != "" {
			text += " " + strings.ToLower(msg)
		}

		score := 0
		for _, k := range keywords {
			if strings.Contains(text, k) {
				score += 2
			}
		}

		if lvl, ok := e.Raw["level"].(string); ok {
			l := strings.ToLower(lvl)
			if l == "error" || l == "fatal" {
				score += 3
			}
			if l == "warn" || l == "warning" {
				score += 1
			}
		}

		if score > 0 {
			list = append(list, scored{e: e, score: score})
		}
	}

	if len(list) == 0 {
		s.AddFindings(models.Finding{
			Agent:       a.Name(),
			Title:       "No strong log signals detected",
			Detail:      "No log evidence matched the initial keyword set (timeouts, retries, errors).",
			EvidenceIDs: nil,
			Confidence:  0.35,
		})
		return nil
	}

	sort.Slice(list, func(i, j int) bool { return list[i].score > list[j].score })

	top := min(5, len(list))
	eids := make([]string, 0, top)
	lines := make([]string, 0, top)

	for i := 0; i < top; i++ {
		eids = append(eids, list[i].e.ID)
		lines = append(lines, fmt.Sprintf("- %s (%s) score=%d: %s",
			list[i].e.Service,
			list[i].e.Timestamp.Format("15:04:05Z"),
			list[i].score,
			list[i].e.Summary,
		))
	}

	s.AddFindings(models.Finding{
		Agent:       a.Name(),
		Title:       "High-signal log entries",
		Detail:      "Top log signals:\n" + strings.Join(lines, "\n"),
		EvidenceIDs: eids,
		Confidence:  0.75,
	})

	return nil
}
