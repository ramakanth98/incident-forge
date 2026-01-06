package report

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ramakanth98/incident-forge/internal/models"
)

func WriteMarkdown(outDir string, inc models.Incident, evidence []models.Evidence, findings []models.Finding, journal []models.JournalEvent) (string, error) {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir: %w", err)
	}

	// Evidence index by ID for linking.
	evByID := make(map[string]models.Evidence, len(evidence))
	for _, e := range evidence {
		evByID[e.ID] = e
	}

	sort.Slice(findings, func(i, j int) bool {
		return findings[i].Confidence > findings[j].Confidence
	})

	var b strings.Builder
	b.WriteString("# Incident Report\n\n")
	b.WriteString(fmt.Sprintf("**ID:** %s\n\n", inc.ID))
	b.WriteString(fmt.Sprintf("**Title:** %s\n\n", inc.Title))
	b.WriteString(fmt.Sprintf("**Window:** %s → %s\n\n", inc.StartTime.Format("2006-01-02 15:04:05Z"), inc.EndTime.Format("2006-01-02 15:04:05Z")))
	if len(inc.Services) > 0 {
		b.WriteString(fmt.Sprintf("**Services:** %s\n\n", strings.Join(inc.Services, ", ")))
	}

	b.WriteString("## Investigation Journal\n\n")
	// Sort by time just in case.
	sort.Slice(journal, func(i, j int) bool { return journal[i].Timestamp.Before(journal[j].Timestamp) })

	for _, je := range journal {
		if je.Agent != "" {
			if je.Agent != "" && je.Message == "agent finished" {
				b.WriteString(fmt.Sprintf("- %s [%s] %s (%s, %dms)\n",
					je.Timestamp.Format("15:04:05Z"),
					je.Type,
					je.Message,
					je.Agent,
					je.DurationMs,
				))
				continue
			} else {
				b.WriteString(fmt.Sprintf("- %s [%s] %s (%s)\n",
					je.Timestamp.Format("15:04:05Z"),
					je.Type,
					je.Message,
					je.Agent,
				))
			}
		} else {
			b.WriteString(fmt.Sprintf("- %s [%s] %s\n",
				je.Timestamp.Format("15:04:05Z"),
				je.Type,
				je.Message,
			))
		}
	}
	b.WriteString("\n")
	b.WriteString("## Findings\n\n")
	for _, f := range findings {
		b.WriteString(fmt.Sprintf("### %s (%.2f)\n\n", f.Title, f.Confidence))
		b.WriteString(fmt.Sprintf("- **Agent:** %s\n", f.Agent))
		b.WriteString("\n")
		b.WriteString(f.Detail)
		b.WriteString("\n\n")
		if len(f.EvidenceIDs) > 0 {
			b.WriteString("**Evidence:**\n")
			for _, id := range f.EvidenceIDs {
				e, ok := evByID[id]
				if !ok {
					continue
				}
				b.WriteString(fmt.Sprintf("- `%s` [%s/%s] %s\n", e.ID, e.Type, e.Service, e.Summary))
			}
			b.WriteString("\n")
		}
	}

	p := filepath.Join(outDir, "report.md")
	if err := os.WriteFile(p, []byte(b.String()), 0o644); err != nil {
		return "", fmt.Errorf("write report: %w", err)
	}
	return p, nil
}
