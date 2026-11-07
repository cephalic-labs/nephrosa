# Renal Relay

Renal Relay is a clinician-led follow-through workspace for kidney monitoring in people with diabetes. It keeps clinician-ordered uACR and eGFR monitoring episodes visible until testing, report receipt, clinician acknowledgement, and a documented next action are accounted for.

## Project status

This repository contains a care-team demonstration implementation: authenticated clinician/coordinator sessions, patient/episode creation, audited workflow actions, a Sarvam Samvaad instant-calling adapter, and guarded clinician closeout. The running services, a real Samvaad call and ngrok callback, and fictional coordinator/clinician closeout were verified on 7 November 2026. Use fictional clinical records and a consenting demonstration phone.

Phase 1 foundations and Phase 2 core workflow are implemented. Phase 3 voice integration has a demonstrated English calling/callback path, with outcome and language verification still pending. Phase 4 demo and quality work is in progress. See the [progress ledger](docs/scaffold-plan.md) for current evidence and remaining gates.

Start with the [Samvaad demo guide](docs/guides/samvaad-demo.md) for required configuration, agent instructions, callback security, test commands, and the presenter script.

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
- [Episode detail surface brief](docs/design/episode-detail.md)

The documents distinguish proposal, accepted decisions, scaffold, and implemented behavior. See the [source notes](docs/source-notes.md) for the supplied concept materials and unresolved decisions.

## Current application

The root Next.js app renders the dashboard and episode details from the Go API. Authenticated actions create episodes, record consent, initiate calls, resolve coordinator work, verify tests, record report receipt, acknowledge reports, and close episodes. Backend checks enforce permissions and closure requirements. Closed episodes remain accessible through dashboard history. Samvaad handles voice/telephony; the application stores validated patient-reported outcomes and short notes, not transcripts or recordings. Calls remain disabled until the documented configuration and callback-security checkpoint is complete.

## Prototype outcome

The central measure is the share of ordered monitoring episodes for which the test is completed, the report is acknowledged by a clinician, and a next action is recorded within the clinician-set window. The proposed pilot is 60–90 days in one diabetes clinic, with baseline agreed with the clinical team and weekly exception review.
