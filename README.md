# Renal Relay

Renal Relay is a clinician-led follow-through workspace for kidney monitoring in people with diabetes. It keeps clinician-ordered uACR and eGFR monitoring episodes visible until testing, report receipt, clinician acknowledgement, and a documented next action are accounted for.

## Project status

This repository is an initialized Next.js application and the starting point for the proposed prototype. The product and technical documents describe intended behavior and planned components; they do not claim those workflows or integrations are implemented. The present task establishes documentation and a project structure plan only.

## Product boundary

- Care teams own clinical decisions, approved outreach content, escalation rules, report review, and episode closure.
- Patient outreach requires recorded consent and respects language preference and opt-out.
- Voice and language-model components may capture structured follow-up information. They must not diagnose, interpret laboratory results, calculate clinical risk, or recommend treatment.
- Use fictional or appropriately de-identified records in the prototype. Do not use real patient data without institutional privacy and security review.

## Documentation

- [Product brief](docs/product-brief.md)
- [Requirements and acceptance criteria](docs/requirements.md)
- [Workflow and safeguards](docs/workflow.md)
- [Architecture and boundaries](docs/architecture.md)
- [Data model and API proposal](docs/data-and-api.md)
- [Project structure](docs/project-structure.md)
- [Roadmap and open decisions](docs/roadmap.md)
- [Source notes](docs/source-notes.md)

## Current application

The repository contains the initialized Next.js app at its root and empty directories marking the planned backend, database, voice integration, and deployment boundaries. No service code, migrations, containers, or integrations have been implemented. See `package.json` for the available commands.

## Prototype outcome

The central measure is the share of ordered monitoring episodes for which the test is completed, the report is acknowledged by a clinician, and a next action is recorded within the clinician-set window. The proposed pilot is 60–90 days in one diabetes clinic, with baseline agreed with the clinical team and weekly exception review.
