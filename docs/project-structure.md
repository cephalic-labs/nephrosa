# Project structure plan

## Current repository

The initialized repository is a Next.js application at its root (`app/`, `public/`, and root configuration files). Empty tracked directories now mark the planned backend, database, voice integration, and deployment boundaries. They contain no implementation code. The directory tree below describes the scaffold and later implementation destinations.

## Proposed monorepo layout

```text
.
├── app/                         # Existing Next.js application
├── public/                      # Existing static assets
├── docs/                        # Product and technical documentation
├── backend/
│   ├── cmd/api/                  # Go API entry point
│   └── internal/                 # API, domain workflow, and persistence packages
├── db/
│   ├── migrations/               # PostgreSQL schema changes
│   └── seed/                     # Fictional demo data
├── integrations/
│   └── voice/                    # Provider adapter boundary (later phase)
└── deploy/                       # Local service orchestration (later phase)
```

The exact Go package boundaries, migration tool, and deployment files are undecided. The existing root Next.js app can remain the web workspace; avoid moving it until a concrete repository need justifies that change.

## Planned ownership boundaries

- `app/`: clinician/coordinator UI and optional patient-facing UI.
- `backend/`: authorization, workflow rules, validated state transitions, and audit creation.
- `db/`: versioned schema and fictional seed records.
- `integrations/voice/`: provider-specific transport and speech integration, isolated from clinical workflow decisions.
- `deploy/`: local development dependencies once the chosen stack is ready.
- `docs/`: source of truth for product requirements, workflow, architecture, data/API proposal, and open decisions.

## Scaffolding sequence

Create implementation directories together with their first meaningful artifacts in later implementation work. Avoid empty folders, placeholder source files, and speculative dependencies. The first code phase should establish the local service contract and fictional data flow before voice integration.
