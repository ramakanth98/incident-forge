package agents

import (
	"context"

	"github.com/ramakanth98/incident-forge/internal/models"
)

type Store interface {
	Incident() models.Incident
	Evidence() []models.Evidence
	AddFindings(...models.Finding)
}

type Agent interface {
	Name() string
	Run(ctx context.Context, s Store) error
}
