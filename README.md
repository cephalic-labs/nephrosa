# Renal Relay

Renal Relay is a clinician-led follow-through workspace for kidney monitoring in people with diabetes. It keeps clinician-ordered uACR and eGFR monitoring episodes visible until testing, report receipt, clinician acknowledgement, and a documented next action are accounted for.

## Project status

This repository contains the Phase 1 Renal Relay prototype foundation: a Next.js care-team dashboard, Go API, PostgreSQL schema and fictional demo seed, plus a local Compose environment. It is a demonstration scaffold, not a clinical system. Authentication, authorization, workflow mutations, integrations, and production deployment remain unimplemented.

## Product boundary

- Care teams own clinical decisions, approved outreach content, escalation rules, report review, and episode closure.
- Patient outreach requires recorded consent and respects language preference and opt-out.
- Voice and language-model components may capture structured follow-up information. They must not diagnose, interpret laboratory results, calculate clinical risk, or recommend treatment.
- Use fictional or appropriately de-identified records in the prototype. Do not use real patient data without institutional privacy and security review.

## Documentation

- [Documentation index and source precedence](docs/README.md)
- [Product record](PRODUCT.md)
- [MVP product specification](docs/product/mvp-spec.md)
- [Proposed HLD/LLD architecture](docs/architecture/renal-relay-hld-lld.md)
- [Architecture decision record](docs/adr/0001-monorepo-scaffold-and-runtime-boundaries.md)
- [Scaffold plan and progress ledger](docs/scaffold-plan.md)
- [Local development](docs/guides/local-development.md)
- [Security and local secrets](docs/guides/security-secrets.md)
- [Dashboard surface brief](docs/design/dashboard.md)

The documents distinguish proposal, accepted decisions, scaffold, and implemented behavior. See the [source notes](docs/source-notes.md) for the supplied concept materials and unresolved decisions.

## Current application

The root Next.js app renders a care-team overview from the Go dashboard API when it is available and clearly labeled fictional sample records otherwise. The API exposes liveness/readiness and a read-only overview endpoint. PostgreSQL migrations establish workflow tables; a separate command seeds two fictional demonstration episodes. Docker Compose wires the local services and secret files. No authentication, workflow mutation endpoints, health-system connection, voice integration, or production hardening is included. See the local development guide for startup steps and limits.

## Prototype outcome

The central measure is the share of ordered monitoring episodes for which the test is completed, the report is acknowledged by a clinician, and a next action is recorded within the clinician-set window. The proposed pilot is 60–90 days in one diabetes clinic, with baseline agreed with the clinical team and weekly exception review.
