package orchestrator

type Budgets struct {
	MaxEvidencePerAgent int
}

func DefaultBudgets() Budgets {
	return Budgets{
		MaxEvidencePerAgent: 500, // safe default
	}
}
