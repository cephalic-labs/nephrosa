# Renal Relay documentation

This directory is the working source of truth for Renal Relay product scope, system boundaries, implementation decisions, and contributor guidance. The prototype baseline is user-confirmed; clinical policy and technical architecture decisions remain proposed or open where marked.

The current demonstration implementation and its verification evidence are described in [`guides/samvaad-demo.md`](guides/samvaad-demo.md). The [progress ledger](scaffold-plan.md) tracks current phase status; the roadmap remains a proposed sequence.

## Source precedence

When project documents disagree, use this order:

1. **User-confirmed product truth and scope** in [`../PRODUCT.md`](../PRODUCT.md) and [`product/mvp-spec.md`](product/mvp-spec.md).
2. **Accepted ADRs** in [`adr/`](adr/) for accepted engineering decisions. ADR-0001 is currently proposed, so it is not yet binding.
3. **Architecture baseline** in [`architecture/renal-relay-hld-lld.md`](architecture/renal-relay-hld-lld.md) for proposed implementation boundaries and contracts.
4. **Detailed workflow, data/API, requirements, and structure** in the linked supporting documents below.
5. **Scaffold ledger** in [`scaffold-plan.md`](scaffold-plan.md) for current implementation status; the ledger records progress but does not redefine product behavior.

`PRODUCT.md` records durable product truth. `DESIGN.md` records the interface design system. They guide future work but cannot override a safety boundary or an accepted architecture decision.

## Product and design

- [`../PRODUCT.md`](../PRODUCT.md) — users, purpose, positioning, constraints, evidence, and unresolved product facts.
- [`../DESIGN.md`](../DESIGN.md) — current Renal Relay care-team dashboard design system, documented from the Phase 1 implementation.
- [`product/mvp-spec.md`](product/mvp-spec.md) — user-confirmed prototype scope, roles, journeys, functional requirements, safeguards, and acceptance.
- [`product-brief.md`](product-brief.md) — concise concept summary.
- [`requirements.md`](requirements.md) — prioritized requirements and safety constraints.
- [`workflow.md`](workflow.md) — episode lifecycle, closure guard, and audit expectations.
- [`roadmap.md`](roadmap.md) — proposed delivery phases and unresolved decisions.

## Architecture and decisions

- [`architecture/renal-relay-hld-lld.md`](architecture/renal-relay-hld-lld.md) — system context, runtime boundaries, data/API overview, security, and non-functional requirements.
- [`data-and-api.md`](data-and-api.md) — proposed entities and REST routes.
- [`adr/0001-monorepo-scaffold-and-runtime-boundaries.md`](adr/0001-monorepo-scaffold-and-runtime-boundaries.md) — repository layout and authority boundaries.
- [`project-structure.md`](project-structure.md) — current scaffold and planned ownership by directory.
- [`architecture.md`](architecture.md) — concise architecture summary; the HLD/LLD above is the detailed baseline.
- [`scaffold-plan.md`](scaffold-plan.md) — phased implementation ledger and completion evidence.

## Contributor guides

- [`guides/architecture-navigation.md`](guides/architecture-navigation.md) — repository paths and responsibility boundaries.
- [`guides/samvaad-demo.md`](guides/samvaad-demo.md) — agent configuration, real calling, webhook security, live-test results, and the presenter script.
- [`guides/local-development.md`](guides/local-development.md) — local web/API/database workflow and its current limits.
- [`guides/security-secrets.md`](guides/security-secrets.md) — local credential handling and production secret requirements.
- [`design/dashboard.md`](design/dashboard.md) — dashboard audience, hierarchy, states, and prototype interaction boundary.
- [`../backend/README.md`](../backend/README.md) — API commands, authenticated workflow endpoints, and configuration.
- [`../db/README.md`](../db/README.md) — migration, seed, and database-role conventions.
- [`../deploy/README.md`](../deploy/README.md) — local Compose startup, reset, and security boundary.
- [`guides/frontend-ui.md`](guides/frontend-ui.md) — required Impeccable workflow for future UI work.

## Source material

- [`source-notes.md`](source-notes.md) — supplied concept and presentation, reconciliation notes, and unresolved questions.
- The original attachments are maintained outside this repository; see the source notes for their names and use.
