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
| SCAFFOLD-005 | Implement episode persistence and audited workflow. | DONE | Authenticated clinician/coordinator sessions, episode creation, workflow actions, audit history, and guarded closure are implemented. The fictional-data workflow and read-back are recorded in the [demo guide](guides/samvaad-demo.md). |
| SCAFFOLD-006 | Build care-team UI using Impeccable workflow. | IN PROGRESS | Public landing, role dashboards, episode creation/detail, and care pathway are implemented. Theme contrast and responsive/accessibility verification remain open; recent CSS changes have not had a complete rendered regression pass. |
| SCAFFOLD-007 | Add voice/telephony integration. | IN PROGRESS | Samvaad calling and callbacks were demonstrated with consenting recipients. Live structured-barrier creation, multilingual accuracy, outcome review against source audio, and approved clinical scripts remain open. |
| SCAFFOLD-008 | Prepare real-patient pilot. | BLOCKED | Requires clinical co-lead, institutional privacy/security review, retention policy, and approved clinical workflow. |

## Current implementation checkpoint

Phase 1 foundations and Phase 2 core workflow are implemented for the fictional-data demonstration. Phase 3 voice integration has a verified English calling/callback path, with the gaps listed above. Phase 4 quality and demo work is in progress: the end-to-end demonstration is recorded, while interface verification and deployment preparation remain open. No clinical pilot or patient outcome is claimed.

## Verification checkpoint — 7 November 2026

| Check | Result | Scope |
| --- | --- | --- |
| Local frontend and API readiness | PASS | Landing route and API readiness each returned HTTP 200. |
| Frontend lint and TypeScript | PASS | ESLint and `tsc --noEmit` completed. |
| Production build | PASS | `next build --webpack` compiled, checked types, and generated routes after selecting Next's documented TypeScript compiler API. |
| Go suite | PASS | `go test ./...` completed with its required module cache and local socket permissions. Database integration tests require a disposable `TEST_DATABASE_URL` and were not run in this pass. |
| Rendered theme and responsive review | NOT RUN | The current theme changes have not had a complete desktop/mobile browser pass. |
| Remaining voice outcomes and multilingual quality | NOT RUN | Requires configured provider tests and an approved evaluation method. |

## Next work

1. Complete light/dark, desktop/mobile, keyboard, and contrast checks across the care-team journey; fix observed defects.
2. Verify the remaining voice outcomes and extraction quality with an approved evaluation method.
3. Review and accept or revise the proposed ADR, clinical scripts, escalation rules, consent policy, and privacy/security controls.
4. Prepare an operational deployment and pilot plan only after the named clinical and institutional decisions are made.
