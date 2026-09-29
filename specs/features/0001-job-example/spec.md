# Feature: 0001 — create and read a Job

Status: validated

## Actor and outcome

An authenticated API client creates a Job and reads its status.

## Contract and acceptance examples

| Given | When | Then |
| --- | --- | --- |
| Valid API key, name and new idempotency key | POST `/v1/jobs` | 201, stable ID, one Job and one outbox event |
| Same key and same name | Repeat POST | Same Job; no new event |
| Same key and different name | Repeat POST | 409 conflict |
| Unknown ID | GET `/v1/jobs/{id}` | 404 |
| Missing/invalid key | Business request | 401 |

## Invariants and recovery

Name is 1–120 characters after trimming. The database transaction contains both Job and event. A process crash before commit leaves neither; a retry after commit returns the same Job. Event delivery is at least once, using a stable event ID.

## Non-goals

The example does not model authorization beyond a bootstrap API key, job processing, or a specific external broker. Replace those with product-specific specs before production.

## Dependency rationale

One PostgreSQL driver is necessary for durable transactional storage. No CQRS or event bus framework is needed for one example.
