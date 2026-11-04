# Proposed architecture

This is a concise architecture summary for navigation. The detailed architecture baseline is [`architecture/renal-relay-hld-lld.md`](architecture/renal-relay-hld-lld.md). The initial local read path exists; this proposal is not accepted as production architecture.

## Components

| Component | Proposed responsibility |
| --- | --- |
| Web workspace | Implemented Next.js read-only dashboard for the care-team overview. Episode detail and optional patient page remain future work. |
| API service | Implemented Go/Gin health and dashboard overview read endpoints. Authentication, authorization, writes, and audit event generation remain future work. |
| Database | Implemented PostgreSQL initial schema and fictional seed for patients, consents, episodes, orders, tasks, contacts, reports, reviews, and events. |
| Voice adapter | Telephony provider plus a speech service candidate such as Sarvam for supported Indic-language speech-to-text and text-to-speech. |
| Constrained language model | Use only clinician-approved scripts and structured extraction for follow-up status/barriers. |
| Local services | Docker Compose currently runs PostgreSQL and API for local prototype use. |
| Initial data input | Manual entry and CSV import; no EHR dependency in the first demo. |

## Trust boundary

The proposed voice/model component is an untrusted input producer. It can suggest a structured follow-up outcome. Before writes are introduced, the Go service must validate consent, schema, identity, permissions, ownership, and allowed transitions. The care team verifies test completion and report receipt; the clinician owns report acknowledgement, next action, and closure. The current API has no write endpoints.

## Integration sequence

1. Demonstrate the core workflow with manual/CSV entry and fictional data.
2. Add a voice adapter only after a provider, supported languages, consent behavior, and approved scripts are selected.
3. Consider additional integrations only after the prototype boundaries and institutional review are clear.

## Security and privacy considerations

The prototype should minimize collected health details and avoid unnecessary call recordings. A real pilot requires decisions and review for authentication, authorization, secure storage, patient verification, access logging, retention, and institutional privacy/security requirements. No claim of regulatory compliance is implied by this proposal.

## Repository starting point

The repository uses a root-level Next.js app with the initial Go API, PostgreSQL schema/seed, and local Compose setup. See [Project structure](project-structure.md) for current paths and boundaries. The architecture described here remains a proposal outside those scoped Phase 1 artifacts.
