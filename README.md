# Incident-forge

Incident-forge is a lightweight incident investigation engine that ingests
telemetry bundles (logs, metrics, and change events) and produces an
evidence-backed incident report.

The project is intentionally **CLI-first**, **deterministic**, and
**connector-agnostic**, with a focus on clarity, reproducibility, and safe
analysis over automation magic.

---

## Design Principles

- **Evidence-first**  
  Every finding or hypothesis must reference concrete evidence.

- **Deterministic by default**  
  The same input bundle should produce the same output.

- **Agent isolation**  
  Each agent has a narrow responsibility and bounded scope.

- **Human-in-the-loop**  
  The tool assists investigation; it does not auto-remediate.

---

## Current Capabilities

- Load incident bundles from disk (`incident.json`, `evidence.json`)
- Run multiple investigation agents concurrently
- Enforce analysis budgets (e.g. max evidence per agent)
- Maintain a full investigation journal with timing
- Generate a Markdown incident report with evidence references

---

## Usage

Run an investigation against a replayable incident bundle:

```bash
go run ./cmd/forge investigate ./testdata/incidents/incident-001
