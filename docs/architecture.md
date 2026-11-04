# Proposed architecture

This is a concise architecture summary for navigation. The detailed architecture baseline is [`architecture/renal-relay-hld-lld.md`](architecture/renal-relay-hld-lld.md). Both describe proposed services and do not claim that those services are implemented.

## Components

| Component | Proposed responsibility |
| --- | --- |
| Web workspace | Next.js with TypeScript for clinician and coordinator queues, episode detail, and the optional patient page. |
| API service | Go with Gin for REST endpoints, authorization, input validation, episode transitions, and audit events. |
| Database | PostgreSQL for patients, consents, episodes, orders, tasks, contacts, reports, reviews, and events. |
| Voice adapter | Telephony provider plus a speech service candidate such as Sarvam for supported Indic-language speech-to-text and text-to-speech. |
| Constrained language model | Use only clinician-approved scripts and structured extraction for follow-up status/barriers. |
| Local services | Docker Compose is proposed for local prototype services. |
| Initial data input | Manual entry and CSV import; no EHR dependency in the first demo. |

## Trust boundary

The voice/model component is an untrusted input producer. It can suggest a structured follow-up outcome. The Go service remains responsible for validating consent, schema, identity, permissions, ownership, and allowed transitions before any state change. The care team verifies test completion and report receipt; the clinician owns report acknowledgement, next action, and closure.

## Integration sequence

1. Demonstrate the core workflow with manual/CSV entry and fictional data.
2. Add a voice adapter only after a provider, supported languages, consent behavior, and approved scripts are selected.
3. Consider additional integrations only after the prototype boundaries and institutional review are clear.

## Security and privacy considerations

The prototype should minimize collected health details and avoid unnecessary call recordings. A real pilot requires decisions and review for authentication, authorization, secure storage, patient verification, access logging, retention, and institutional privacy/security requirements. No claim of regulatory compliance is implied by this proposal.

## Repository starting point

The repository is currently a root-level Next.js app. The proposed monorepo layout is captured in [Project structure](project-structure.md). This documentation phase does not add service code, migrations, containers, dependencies, or integrations.
