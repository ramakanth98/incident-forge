# Incident Report

**ID:** incident-001

**Title:** Checkout API latency spike after deploy

**Window:** 2026-01-05 18:10:00Z → 2026-01-05 18:40:00Z

**Services:** checkout-api, payments-worker

## Investigation Journal

- 14:54:47Z [agent] agent started (change-correlation)
- 14:54:47Z [agent] agent started (log-signal)
- 14:54:47Z [agent] agent started (metrics-anomaly)
- 14:54:47Z [agent] agent finished (change-correlation, 0ms)
- 14:54:47Z [agent] agent finished (log-signal, 0ms)
- 14:54:47Z [agent] agent finished (metrics-anomaly, 0ms)

## Findings

### Metric anomalies in incident window (0.80)

- **Agent:** metrics-anomaly

Top anomalies:
- checkout-api http_server_latency_p95_ms: +991% (220 → 2400)

**Evidence:**
- `met-001` [metric/checkout-api] p95 latency jumped from 220ms to 2.4s

### Changes near incident start (0.75)

- **Agent:** change-correlation

Top change events close to incident start:
- checkout-api (18:08:30Z) Deployed v1.12.0 (added retry logic to downstream call)

**Evidence:**
- `chg-001` [change/checkout-api] Deployed v1.12.0 (added retry logic to downstream call)

### High-signal log entries (0.75)

- **Agent:** log-signal

Top log signals:
- checkout-api (18:13:10Z) score=5: Timeout calling payments-service; retrying

**Evidence:**
- `log-001` [log/checkout-api] Timeout calling payments-service; retrying

