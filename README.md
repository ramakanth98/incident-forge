# incident-forge

incident-forge is a lightweight incident investigation engine that ingests
telemetry bundles (logs, metrics, change events) and produces an evidence-backed
incident report.

The focus is on:
- deterministic analysis
- evidence-first findings
- reproducible investigations
- clear agent boundaries

This project is intentionally CLI-first and connector-agnostic.

## Current Capabilities
- Load incident bundles from disk
- Correlate change events near incident start
- Generate a markdown incident report with evidence references

## Usage
```bash
go run ./cmd/forge investigate ./testdata/incidents/incident-001
