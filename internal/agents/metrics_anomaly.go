package agents

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/ramakanth98/incident-forge/internal/models"
)

type MetricsAnomalyAgent struct{}

func (a *MetricsAnomalyAgent) Name() string { return "metrics-anomaly" }

func (a *MetricsAnomalyAgent) Run(ctx context.Context, s Store) error {
	inc := s.Incident()
	ev := s.EvidenceLimited(500)

	type metricHit struct {
		e   models.Evidence
		pct float64
	}

	hits := make([]metricHit, 0)

	for _, e := range ev {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if e.Type != models.EvidenceMetric {
			continue
		}
		if len(inc.Services) > 0 && !contains(inc.Services, e.Service) {
			continue
		}

		before, okB := asFloat(e.Raw["before"])
		after, okA := asFloat(e.Raw["after"])

		// If we can't compute a percentage, still keep it as a "weak hit".
		if !okB || !okA || before <= 0 {
			hits = append(hits, metricHit{e: e, pct: 0})
			continue
		}

		pct := ((after - before) / before) * 100.0
		if pct >= 50.0 {
			hits = append(hits, metricHit{e: e, pct: pct})
		}
	}

	if len(hits) == 0 {
		s.AddFindings(models.Finding{
			Agent:       a.Name(),
			Title:       "No significant metric anomalies detected",
			Detail:      "No metric evidence exceeded the anomaly threshold (+50%) in the current evidence set.",
			EvidenceIDs: nil,
			Confidence:  0.35,
		})
		return nil
	}

	sort.Slice(hits, func(i, j int) bool { return hits[i].pct > hits[j].pct })

	top := min(5, len(hits))
	eids := make([]string, 0, top)
	lines := make([]string, 0, top)

	for i := 0; i < top; i++ {
		h := hits[i]
		eids = append(eids, h.e.ID)

		if h.pct > 0 {
			lines = append(lines, fmt.Sprintf("- %s %s: +%.0f%% (%v → %v)",
				h.e.Service,
				metricName(h.e),
				h.pct,
				h.e.Raw["before"],
				h.e.Raw["after"],
			))
		} else {
			lines = append(lines, fmt.Sprintf("- %s %s: %s",
				h.e.Service,
				metricName(h.e),
				h.e.Summary,
			))
		}
	}

	s.AddFindings(models.Finding{
		Agent:       a.Name(),
		Title:       "Metric anomalies in incident window",
		Detail:      "Top anomalies:\n" + strings.Join(lines, "\n"),
		EvidenceIDs: eids,
		Confidence:  0.8,
	})

	return nil
}

func metricName(e models.Evidence) string {
	if v, ok := e.Raw["metric"].(string); ok && v != "" {
		return v
	}
	return "metric"
}

func asFloat(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}
