# Renal Relay scaffold plan and progress ledger

This ledger tracks repository scaffolding and implementation evidence. A planned item is not complete because it appears in the directory tree or a design document. Update status only when the stated evidence exists.

## Status definitions

- **DONE** — artifact exists and its stated acceptance condition has been checked.
- **IN PROGRESS** — work has started, but acceptance evidence is incomplete.
- **PLANNED** — not implemented.
- **BLOCKED** — cannot proceed because a named decision or dependency is unresolved.

## Scaffold items

| ID | Work item | Status | Acceptance evidence / remaining work |
| --- | --- | --- | --- |
| SCAFFOLD-001 | Record product proposal, requirements, workflow and source notes. | DONE | Markdown source set exists under root and `docs/`; proposal status and source limitations are explicit. |
| SCAFFOLD-002 | Establish docs index, architecture baseline, ADR and progress ledger. | IN PROGRESS | The implementation is recorded; proposed architecture and ADR remain proposals pending owner acceptance. |
| SCAFFOLD-003 | Record Impeccable product context and design system. | DONE | `PRODUCT.md`, Renal Relay `DESIGN.md`, and Impeccable sidecar document the accepted prototype baseline and implemented dashboard system. |
| SCAFFOLD-004 | Create application and service foundations. | DONE | Next.js dashboard, Go/Gin API, PostgreSQL migration/seed files, Docker Compose and local secret generation are present. |
| SCAFFOLD-005 | Implement initial episode read model and persistence. | IN PROGRESS | Initial schema, read-only dashboard overview and episode detail, and fictional seed exist. Authentication, write transitions, and audit-producing workflow commands remain unimplemented. |
| SCAFFOLD-006 | Build care-team UI using Impeccable workflow. | DONE | Server-rendered dashboard, branded system, and narrow layout are implemented and inspected in desktop and mobile browsers. Full API/database container runtime remains unverified because this environment cannot access the Docker daemon. |
| SCAFFOLD-007 | Add voice/telephony integration. | BLOCKED | Provider, supported languages, consent wording, approved scripts, and escalation rules are undecided. |
| SCAFFOLD-008 | Prepare real-patient pilot. | BLOCKED | Requires clinical co-lead, institutional privacy/security review, retention policy, and approved clinical workflow. |

## Current implementation checkpoint

Phase 1 provides dashboard and episode-detail read paths over fictional episodes, initial PostgreSQL schema, a demo seed, API health/overview/detail endpoints, and local Compose orchestration with separate database roles and secret files. The workflow is not yet mutable through the API. Authentication, patient-facing flows, and voice integration are not implemented. No clinical pilot or outcome is claimed.

## Next work

1. Review and accept or revise the proposed ADR and architecture baseline.
2. Add authentication/authorization and validated, audited workflow mutations before any operational use.
3. Inspect the dashboard in a running browser and complete a bounded visual/accessibility review.
4. Keep voice, integrations, and real-patient pilot work gated on explicit clinical, privacy, and security decisions.
