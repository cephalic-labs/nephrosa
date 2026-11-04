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
| SCAFFOLD-002 | Establish docs index, architecture baseline, ADR and progress ledger. | IN PROGRESS | Documents drafted; architecture and ADR require project-owner acceptance. |
| SCAFFOLD-003 | Record Impeccable product context and document incumbent UI. | DONE | Product baseline confirmed; `PRODUCT.md`, starter-only `DESIGN.md`, and Impeccable sidecar exist. Visual direction remains explicitly open. |
| SCAFFOLD-004 | Create empty application/service boundaries without source code. | DONE | `backend/`, `db/`, `integrations/voice/`, and `deploy/` contain tracked placeholders only. |
| SCAFFOLD-005 | Implement core patient/episode workflow and persistence. | PLANNED | Requires accepted architecture, domain contract, and a separate implementation task. |
| SCAFFOLD-006 | Build care-team UI using Impeccable workflow. | PLANNED | Requires product/UI direction to be recorded and a dedicated UI implementation task. |
| SCAFFOLD-007 | Add voice/telephony integration. | BLOCKED | Provider, supported languages, consent wording, approved scripts, and escalation rules are undecided. |
| SCAFFOLD-008 | Prepare real-patient pilot. | BLOCKED | Requires clinical co-lead, institutional privacy/security review, retention policy, and approved clinical workflow. |

## Current implementation checkpoint

The repository has the initialized Next.js scaffold plus documentation and empty planned-service boundaries. The Renal Relay workflow, backend, database, authentication, patient-facing flows, and voice integration are not implemented. No clinical pilot or outcome is claimed.

## Next work

1. Accept or revise the proposed ADR and architecture baseline.
2. Select the first implementation slice and add its contracts before service code.
3. Establish a Renal Relay visual direction with Impeccable before replacing the starter page.
4. Use Impeccable for UI implementation and bounded visual review.
