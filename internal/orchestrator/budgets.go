package orchestrator

type Budgets struct {
	MaxEvidencePerAgent int
}

func (b Budgets) WithDefaults() Budgets {
	if b.MaxEvidencePerAgent <= 0 {
		b.MaxEvidencePerAgent = 500
	}
	return b
}
