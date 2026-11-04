# Project structure

The repository keeps the initialized Next.js app at its root and adds the first local API/database boundaries. This describes current ownership, not a commitment to production topology.

```text
.
├── app/                         # Next.js dashboard and visual system
├── backend/
│   ├── cmd/api/                  # API, migration, seed, and healthcheck commands
│   └── internal/                 # Config, database, HTTP, migration, seed code
├── db/
│   ├── init/                     # Local PostgreSQL role initialization
│   ├── migrations/               # Embedded schema migrations
│   └── seed/                     # Fictional demonstration records
├── deploy/compose.yaml           # Local database and API services
├── docs/                         # Product, design, architecture, and guides
├── integrations/voice/           # Reserved for a future provider adapter
├── scripts/dev-secrets.sh         # Local-only credential generation
└── secrets/                       # Generated local credentials; ignored by Git
```

## Current ownership

- `app/`: server-rendered care-team dashboard; currently read-only.
- `backend/`: health checks and dashboard read endpoint, database configuration, and one-shot migration/seed commands. No authentication or workflow mutation rules exist yet.
- `db/`: initial workflow schema, local role setup, and fictional seed only.
- `deploy/`: local Compose service configuration; not a production deployment manifest.
- `integrations/voice/`: future boundary; no provider code or credentials.
- `docs/`: product truth, design system, prototype limits, proposed architecture, and contributor guides.

## Next foundation work

Implement and review authentication/authorization, validated workflow mutations, transactional audit events, clinical policy constraints, and a non-demo deployment boundary before operational use. Keep patient and report data out of the demo until institutional review and data-handling decisions are complete.
