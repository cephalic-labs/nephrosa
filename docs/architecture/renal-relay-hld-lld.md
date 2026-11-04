# Renal Relay high-level and low-level architecture

| Field | Value |
| --- | --- |
| Status | Proposed baseline for review; not yet an accepted implementation architecture |
| Version | 0.1 |
| Date | 4 November 2026 |
| Product spec | [`../product/mvp-spec.md`](../product/mvp-spec.md) |
| Decision record | [`../adr/0001-monorepo-scaffold-and-runtime-boundaries.md`](../adr/0001-monorepo-scaffold-and-runtime-boundaries.md) |

## Document map

| Part | Coverage |
| --- | --- |
| A. System context | Users, runtime responsibilities, and trust boundary. |
| B. High-level design | Web, API, database, voice, and local service topology. |
| C. Low-level design | Episode state, API and data contract proposals. |
| D. Engineering controls | Privacy, security, reliability, accessibility, and observability. |
| E. Delivery | Scaffold sequence and evidence-based acceptance. |

## A. Architecture drivers

| Driver | Required response |
| --- | --- |
| Clinical authority | Clinicians retain report review, next-action, escalation-policy, and closure authority. |
| Workflow integrity | The service validates identity, permissions, consent, evidence, and allowed state transitions. |
| Auditability | Material state changes record actor, time, event type, and relevant provenance. |
| Human follow-through | Barriers and missing evidence become accountable staff work. |
| Privacy | Minimize patient data, use fictional/de-identified demo data, and define retention before real use. |
| Replaceable integrations | Voice/LLM and telephony providers sit behind an adapter and cannot own episode state. |
| Prototype simplicity | Start with manual/CSV data entry; do not depend on EHR/LIS integration. |

## B. High-level system context

```text
Clinician / Coordinator ── HTTPS ──> Next.js care-team workspace
                                            │
                                            ▼
                                   Go workflow API
                                  /               \
                                 ▼                 ▼
                         PostgreSQL        Voice/LLM adapter
                                                   │
                                                   ▼
                                         Telephony / speech provider
```

The diagram shows proposed boundaries. The initial prototype may run without the voice/LLM adapter. The optional patient experience is a later web surface. All state-changing events from integrations pass through the workflow API.

## C. Component responsibilities

| Component | Proposed technology | Authority and responsibility |
| --- | --- | --- |
| Web workspace | Existing root Next.js + TypeScript app | Render care-team tasks and submit user actions. It is not the policy authority. |
| Workflow API | Go + Gin | Authenticate/authorize, validate input and consent, enforce transitions and closure rules, write audit events. |
| System of record | PostgreSQL | Persist users, patients, consent, episodes, orders, tasks, completion evidence, reports, reviews, and audit events. |
| Voice adapter | Telephony + speech provider candidate (Sarvam is one proposal) | Deliver approved outreach and return constrained follow-up observations. Disabled until selected and reviewed. |
| Language model | Hosted model candidate | Extract only bounded scheduling outcomes/barriers from consented interactions using approved scripts. |
| Local orchestration | Docker Compose proposed | Run prototype dependencies locally after service contracts are accepted. |

## D. Trust and data flow

1. Clinician/coordinator submits an action through the web workspace.
2. API authenticates the actor and checks role, assignment, consent, payload shape, and current episode state.
3. API persists the state change and audit event transactionally.
4. A future voice event is treated as untrusted proposed data; API validates it before staff review or allowed workflow update.
5. Staff separately verify test completion and report receipt.
6. Clinician review and next action are required before normal closure.

The model does not receive authority to close episodes, read/interpret lab results, diagnose, calculate risk, or recommend treatment.

## E. Domain state model

Suggested main progression:

```text
ordered → outreach_pending → awaiting_booking / booked
        → awaiting_test_completion → awaiting_report
        → awaiting_clinician_review → closed
```

Barrier, overdue, opted-out, and escalated conditions are orthogonal flags/tasks unless later workflow analysis demonstrates they must be lifecycle states. Full definitions and closure invariants are in [`../workflow.md`](../workflow.md).

## F. API and persistence baseline

The proposed REST surface and entities are documented in [`../data-and-api.md`](../data-and-api.md). These are design proposals, not existing endpoints or a deployed schema. Before implementation, agree on identifiers, validation, audit retention, report handling, and access policy.

Initial writes include patient/episode creation, contact/barrier capture, task assignment, completion verification, report receipt, clinician review, closure, and (later) validated voice events. Mutations must enforce authorization, ownership/assignment rules, consent where applicable, and valid transition constraints.

## G. Quality and safety requirements

- **Security/privacy:** least privilege, secure storage, limited data collection, defined retention, and an institutional review before patient use.
- **Integrity:** record completion evidence separately from report receipt and clinician acknowledgement.
- **Availability:** the core manual workflow should function without voice-provider availability.
- **Accessibility:** care-team workspace should support keyboard use, visible focus, semantic structure, clear status text/icons, and responsive behavior. The target conformance level requires confirmation before UI implementation.
- **Observability:** record service errors and workflow events without logging unnecessary patient details or conversation content.
- **Recovery:** define backups and restoration before any non-disposable environment stores real data.

## H. Deployment and integration sequence

1. Keep the current Next.js app at repository root; add backend and database artifacts when their interfaces are ready.
2. Demonstrate a manual workflow using fictional records.
3. Add CSV import only when the demo data contract is defined.
4. Add voice and telephony only after provider, language quality, consent, scripts, cost, and escalation behavior are reviewed.
5. Consider clinical-system integrations only after institutional requirements and ownership are established.

## I. Assumptions and unresolved architecture choices

- Go/Gin and PostgreSQL are proposed by the supplied plan, not installed dependencies or accepted ADR decisions yet.
- Authentication, role provisioning, patient identity verification, and authorization granularity remain open.
- Report data may be metadata only, document upload, or future integration; no choice is made here.
- Voice provider and model hosting region/retention are undecided.
- Operational deployment and backup requirements are not specified.

