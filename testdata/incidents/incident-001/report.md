# Incident Report

**ID:** incident-001

**Title:** Checkout API latency spike after deploy

**Window:** 2026-01-05 18:10:00Z → 2026-01-05 18:40:00Z

**Services:** checkout-api, payments-worker

## Findings

### Changes near incident start (0.75)

- **Agent:** change-correlation

Top change events close to incident start:
- checkout-api (18:08:30Z) Deployed v1.12.0 (added retry logic to downstream call)

**Evidence:**
- `chg-001` [change/checkout-api] Deployed v1.12.0 (added retry logic to downstream call)

