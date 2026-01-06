package agents

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/ramakanth98/incident-forge/internal/models"
)

type ChangeCorrelationAgent struct{}

func (a *ChangeCorrelationAgent) Name() string { return "change-correlation" }

func (a *ChangeCorrelationAgent) Run(ctx context.Context, s Store) error {
	inc := s.Incident()
	ev := s.EvidenceLimited(500)

	// Pick change events near incident start, rank by closeness.
	type scored struct {
		e models.Evidence
		d int64
	}
	scoredList := make([]scored, 0)

	for _, e := range ev {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if e.Type != models.EvidenceChange {
			continue
		}
		// Only consider services in the incident scope (if provided).
		if len(inc.Services) > 0 && !contains(inc.Services, e.Service) {
			continue
		}

		delta := e.Timestamp.Sub(inc.StartTime)
		if delta < 0 {
			delta = -delta
		}
		scoredList = append(scoredList, scored{e: e, d: int64(delta.Seconds())})
	}

	sort.Slice(scoredList, func(i, j int) bool { return scoredList[i].d < scoredList[j].d })

	top := min(3, len(scoredList))
	if top == 0 {
		s.AddFindings(models.Finding{
			Agent:       a.Name(),
			Title:       "No relevant change events found",
			Detail:      "No deploy/config change events were present in the incident window.",
			EvidenceIDs: nil,
			Confidence:  0.3,
		})
		return nil
	}

	eids := make([]string, 0, top)
	lines := make([]string, 0, top)
	for i := 0; i < top; i++ {
		eids = append(eids, scoredList[i].e.ID)
		lines = append(lines, fmt.Sprintf("- %s (%s) %s", scoredList[i].e.Service, scoredList[i].e.Timestamp.Format("15:04:05Z"), scoredList[i].e.Summary))
	}

	s.AddFindings(models.Finding{
		Agent:       a.Name(),
		Title:       "Changes near incident start",
		Detail:      "Top change events close to incident start:\n" + strings.Join(lines, "\n"),
		EvidenceIDs: eids,
		Confidence:  0.75,
	})

	return nil
}

func contains(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
