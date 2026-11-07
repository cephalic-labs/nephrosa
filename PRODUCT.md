# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Users

The primary users are doctors and care-team members following up on clinician-ordered kidney monitoring for people with diabetes. Clinicians set the order and due window, review the report, record a next action, and authorize closure. Care coordinators work the follow-up queue, address practical barriers, and record evidence and handoffs. Patients receive consented outreach, share whether a test is arranged or completed, request human help, or opt out. Prototype administrators manage demo access and approved templates.

## Product Purpose

Renal Relay is a clinician-led follow-through workspace for uACR/eGFR monitoring episodes. It keeps each episode visible through patient contact, test arrangement and verified completion, report receipt, clinician acknowledgement, and a documented next action. The prototype baseline is the supplied Renal Relay proposal, confirmed by the user on 4 November 2026.

Success means that the care team can verify the required sequence within the clinician-set window. A sent reminder or completed call alone is not a successful episode.

## Positioning

Renal Relay tracks verified closed-loop follow-through instead of counting messages sent. It combines patient access support with an accountable staff workflow and leaves medical judgement with clinicians.

## Operating Context

The product supports follow-up between clinic visits. The first demonstration uses fictional or appropriately de-identified patient and order data, with manual entry as the initial input. The care team works from shared queues and episode timelines. Voice follow-up and a patient mobile page are proposed later extensions, not prerequisites for the initial workflow.

## Capabilities and Constraints

- The prototype scope includes uACR and eGFR episodes, ownership and due windows, contact and barrier logging, separate test-completion and report-receipt records, clinician acknowledgement, next-action documentation, and guarded closure.
- Outreach requires recorded consent and an agreed language preference. Patients may ask for human help or opt out.
- A voice or language-model component may capture structured follow-up outcomes. It must not diagnose, interpret results, calculate clinical risk, or recommend treatment.
- Staff verify test completion and report receipt. Clinicians own report acknowledgement, next action, and standard closure.
- Real-patient use requires clinician co-design and institutional privacy/security review.
- The current local demonstration uses a Next.js web workspace, a Go API, PostgreSQL, authenticated clinician/coordinator sessions, audited episode actions, and an optional Sarvam voice adapter. These are prototype implementation choices, not accepted production architecture.

## Evidence on Hand

- The supplied Renal Relay concept and prototype plan and the eight-slide Health-a-thon presentation define the proposed use case, workflow, safeguards, technical options, and pilot measures. Source filenames and reconciliation notes are recorded in [`docs/source-notes.md`](docs/source-notes.md).
- The repository contains a public product landing page, authenticated role-aware care-team dashboard, workflow guide, episode creation and detail pages, backend workflow actions, and fictional seed records. Clinical integrations and clinical decision support are not implemented.
- No real-patient data, validated clinical scripts, clinical integration, or pilot outcome is included in this repository.

## Product Principles

1. Keep ownership and unfinished work visible between visits.
2. Count verified care continuity, not contact attempts.
3. Keep clinical decisions and episode closure with qualified staff.
4. Capture practical barriers and route them to a person who can resolve them.
5. Minimize sensitive data and use fictional/de-identified records for the prototype.

## Accessibility & Inclusion

The care-team workspace should support readable tables, search and filters, keyboard access, clear status text and icons, and responsive layouts. The optional patient experience should use plain language, large touch targets, strong contrast, and minimal steps. English, Tamil, and Hindi are candidate languages in the proposal; the supported languages and their quality must be confirmed before demo use. No formal conformance level has been approved.
