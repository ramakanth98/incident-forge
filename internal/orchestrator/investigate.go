package orchestrator

import (
	"context"
	"fmt"
	"time"

	"github.com/ramakanth98/incident-forge/internal/agents"
	"github.com/ramakanth98/incident-forge/internal/connectors"
	"github.com/ramakanth98/incident-forge/internal/report"
	"github.com/ramakanth98/incident-forge/internal/store"
)

func RunInvestigate(bundlePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	st := store.NewMemStore()

	loader := connectors.NewBundleLoader()
	inc, err := loader.LoadIncident(bundlePath)
	if err != nil {
		return err
	}
	ev, err := loader.LoadEvidence(bundlePath)
	if err != nil {
		return err
	}

	st.PutIncident(inc)
	st.AddEvidence(ev...)

	// Agents run sequentially first (simple). We'll make this concurrent next.
	agentList := []agents.Agent{
		&agents.ChangeCorrelationAgent{},
	}

	for _, ag := range agentList {
		if err := ag.Run(ctx, st); err != nil {
			return fmt.Errorf("%s failed: %w", ag.Name(), err)
		}
	}

	outPath, err := report.WriteMarkdown(bundlePath, st.Incident(), st.Evidence(), st.Findings())
	if err != nil {
		return err
	}

	fmt.Println("wrote:", outPath)
	return nil
}
