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

	agentList := []agents.Agent{
		&agents.ChangeCorrelationAgent{},
		&agents.MetricsAnomalyAgent{},
		&agents.LogSignalAgent{},
	}

	type agentResult struct {
		name string
		err  error
	}

	results := make(chan agentResult, len(agentList))

	for _, ag := range agentList {
		ag := ag // capture loop variable
		go func() {
			err := ag.Run(ctx, st)
			results <- agentResult{name: ag.Name(), err: err}
		}()
	}

	for i := 0; i < len(agentList); i++ {
		r := <-results
		if r.err != nil {
			return fmt.Errorf("%s failed: %w", r.name, r.err)
		}
	}

	outPath, err := report.WriteMarkdown(bundlePath, st.Incident(), st.Evidence(), st.Findings())
	if err != nil {
		return err
	}

	fmt.Println("wrote:", outPath)
	return nil
}
