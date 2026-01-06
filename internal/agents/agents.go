package agents

import (
	"context"

	"github.com/ramakanth98/incident-forge/internal/models"
)

type Store interface {
	Incident() models.Incident
	Evidence() []models.Evidence
	AddFindings(...models.Finding)
	AddJournal(...models.JournalEvent)
	EvidenceLimited(n int) []models.Evidence
}

type Agent interface {
	Name() string
	Run(ctx context.Context, s Store) error
}
