# Architecture navigation

This guide distinguishes the Phase 1 implementation from proposed authority boundaries. The deployed technology choices remain a prototype baseline pending architecture review.

## Current tree

```text
nephrosa/
├── app/                    # Next.js read-only care-team dashboard
├── backend/                # Go API, configuration, migration and seed commands
├── db/                     # Initial schema and fictional seed
├── deploy/                 # Local PostgreSQL/API Compose setup
├── docs/                   # Product, design, architecture, decisions, and guides
├── integrations/voice/      # Reserved future provider adapter boundary
└── scripts/                # Local-only development setup helpers
```

## Current authority and limits

| Area | Current responsibility |
| --- | --- |
| `app/` | Display the overview and episode queue; no workflow mutation or clinical decision support. |
| `backend/` | Health endpoints, dashboard read model, DB access, migrations, and fictional seeding. No authentication or writes. |
| `db/` | Initial PostgreSQL schema and fictional seed records. |
| `integrations/voice/` | No implementation; provider and scripts remain undecided. |
| `deploy/` | Local service orchestration only. |
| `docs/` | Product truth, proposed architecture and contracts, and implementation status. |

See [the architecture baseline](../architecture/renal-relay-hld-lld.md), [ADR-0001](../adr/0001-monorepo-scaffold-and-runtime-boundaries.md), and [the scaffold ledger](../scaffold-plan.md). None of the current services are approved for real patient data or production use.
