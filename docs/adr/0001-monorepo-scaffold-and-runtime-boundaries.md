# ADR-0001: Establish the Renal Relay scaffold and runtime boundaries

## Status

Proposed — pending product/engineering acceptance.

## Date

4 November 2026

## Context

The initialized repository is a root-level Next.js application. The supplied prototype plan proposes a TypeScript web workspace, Go/Gin REST API, PostgreSQL, and a later voice/LLM integration. It also places workflow validation and clinical authority with the application service and care team.

The repository needs a visible scaffold without implying that those services, their contracts, or clinical safeguards are already implemented.

## Decision proposal

1. Keep the existing Next.js app at the repository root during the first prototype phase.
2. Mark future `backend/`, `db/`, `integrations/voice/`, and `deploy/` boundaries without adding service code until each has an agreed contract and implementation task.
3. If the Go/Gin service is accepted, it owns authorization, consent validation, episode transitions, closure checks, and audit-event creation.
4. PostgreSQL is the proposed system of record for episode workflow data.
5. Voice and language-model components may propose structured follow-up observations but do not own workflow state or clinical decisions.
6. Begin with manual or CSV demo data. Do not require EHR/LIS integration for the prototype.
7. Use fictional or appropriately de-identified data until institutional privacy/security review authorizes anything else.

## Alternatives considered

### Build the voice agent before the workflow service

Rejected for the proposed prototype sequence because the core value is verified follow-through and the agent must remain behind consent, validation, and human review boundaries.

### Treat contact completion as episode completion

Rejected because the source concept requires separate evidence for test completion, report receipt, clinician acknowledgement, and next action.

### Add clinical system integration at the outset

Deferred because no institution, interface, identity model, or integration contract has been selected.

## Consequences

- Empty scaffold directories are intentional and do not indicate a running service.
- Product and architecture docs distinguish proposal from accepted decisions and implementation evidence.
- Later service work must update this ADR if it changes runtime authority or repository boundaries.
- The data model must minimize sensitive information and support attributable workflow events.

## Source documents

- [`../architecture/renal-relay-hld-lld.md`](../architecture/renal-relay-hld-lld.md)
- [`../product/mvp-spec.md`](../product/mvp-spec.md)
- [`../../docs/source-notes.md`](../../docs/source-notes.md)
