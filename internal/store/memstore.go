package store

import (
	"sync"

	"github.com/ramakanth98/incident-forge/internal/models"
)

type MemStore struct {
	mu       sync.RWMutex
	incident models.Incident
	evidence []models.Evidence
	findings []models.Finding
	journal  []models.JournalEvent
}

func NewMemStore() *MemStore {
	return &MemStore{
		journal:  make([]models.JournalEvent, 0, 128),
		evidence: make([]models.Evidence, 0, 256),
		findings: make([]models.Finding, 0, 64),
	}
}

func (s *MemStore) AddJournal(e ...models.JournalEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.journal = append(s.journal, e...)
}

func (s *MemStore) Journal() []models.JournalEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]models.JournalEvent, len(s.journal))
	copy(out, s.journal)
	return out
}

func (s *MemStore) PutIncident(inc models.Incident) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.incident = inc
}

func (s *MemStore) Incident() models.Incident {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.incident
}

func (s *MemStore) AddEvidence(ev ...models.Evidence) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evidence = append(s.evidence, ev...)
}

func (s *MemStore) Evidence() []models.Evidence {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]models.Evidence, len(s.evidence))
	copy(out, s.evidence)
	return out
}

func (s *MemStore) AddFindings(f ...models.Finding) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.findings = append(s.findings, f...)
}

func (s *MemStore) Findings() []models.Finding {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]models.Finding, len(s.findings))
	copy(out, s.findings)
	return out
}

func (s *MemStore) EvidenceLimited(n int) []models.Evidence {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if n <= 0 || n >= len(s.evidence) {
		out := make([]models.Evidence, len(s.evidence))
		copy(out, s.evidence)
		return out
	}
	out := make([]models.Evidence, n)
	copy(out, s.evidence[:n])
	return out
}
