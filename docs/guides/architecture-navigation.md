# Architecture navigation

This guide identifies the current repository areas and proposed ownership boundaries.

## Current tree

```text
nephrosa/
├── app/                    # Existing Next.js App Router starter
├── public/                 # Existing static starter assets
├── docs/                   # Product, architecture, decisions, and guides
├── backend/                # Empty planned API boundaries
├── db/                     # Empty planned migration and seed boundaries
├── integrations/voice/      # Empty planned voice adapter boundary
└── deploy/                 # Empty planned local deployment boundary
```

## Proposed authority boundaries

| Area | Proposed owner |
| --- | --- |
| `app/` | Clinician/coordinator and optional patient interfaces. |
| `backend/` | Authentication/authorization, consent checks, lifecycle validation, and audit-event writes. |
| `db/` | Versioned PostgreSQL schema and fictional seed data. |
| `integrations/voice/` | Provider transport and constrained speech/LLM integration only. |
| `deploy/` | Local service orchestration after the runtime choices are accepted. |
| `docs/` | Product, architecture, workflow, and implementation-status source of truth. |

See [the architecture baseline](../architecture/renal-relay-hld-lld.md) and [ADR-0001](../adr/0001-monorepo-scaffold-and-runtime-boundaries.md). The directories do not imply implemented services.
