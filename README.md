# Renal Relay

Renal Relay is a clinician-led follow-through workspace for kidney monitoring in people with diabetes. It keeps clinician-ordered uACR and eGFR monitoring episodes visible until testing, report receipt, clinician acknowledgement, and a documented next action are accounted for.

## Project status

This repository contains the initial Next.js application and documentation/scaffold for the user-confirmed prototype baseline. The Renal Relay workflows and planned services are not implemented. The architecture and clinical integration choices remain proposals or open decisions.

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

The documents distinguish proposal, accepted decisions, scaffold, and implemented behavior. See the [source notes](docs/source-notes.md) for the supplied concept materials and unresolved decisions.

## Current application

The repository contains the initialized Next.js app at its root and empty directories marking the planned backend, database, voice integration, and deployment boundaries. No service code, migrations, containers, or integrations have been implemented. See `package.json` for the available commands.

## Prototype outcome

The central measure is the share of ordered monitoring episodes for which the test is completed, the report is acknowledged by a clinician, and a next action is recorded within the clinician-set window. The proposed pilot is 60–90 days in one diabetes clinic, with baseline agreed with the clinical team and weekly exception review.
