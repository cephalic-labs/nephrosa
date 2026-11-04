# Proposed data model and API

This document records the supplied prototype design proposal and distinguishes the currently implemented read-only routes. The broader schema and write contracts remain proposals; field types, identifiers, validation rules, and privacy policy still need design review.

## Implemented read routes

| Method and route | Behavior |
| --- | --- |
| `GET /api/v1/dashboard/overview` | Returns the current summary counts and up to ten open episode summaries. |
| `GET /api/v1/episodes/:id` | Returns one open episode's detail from PostgreSQL, including test orders/completion verification, contact history, barriers, tasks, report receipt metadata, clinician acknowledgement/next action, and event labels with actor and time. |

The detail read uses a read-only repeatable-read transaction. It omits report contents and audit-event metadata. Missing or closed episodes return 404; database failures return a generic 503. Neither route authenticates callers. They are limited to the local fictional-data prototype and are not approved for real-patient data or untrusted networks. All proposed mutation routes below remain unimplemented.

## Proposed entities

| Entity | Proposed fields and purpose |
| --- | --- |
| users | `id`, `role`, `display_name`, `active`, `created_at` — prototype users and roles. |
| patients | `id`, `external_reference`, `name`, `preferred_language`, contact fields — minimum patient context needed for follow-through. |
| consents | `patient_id`, `channel`, `status`, `recorded_at`, `source` — channel consent and provenance. |
| episodes | `id`, `patient_id`, `owner_id`, `status`, `due_at`, `created_by`, timestamps — persistent monitoring workflow. |
| test_orders | `id`, `episode_id`, `test_type` (uACR/eGFR), `ordered_at`, `due_at` — requested tests. |
| contact_attempts | `id`, `episode_id`, `channel`, `outcome`, `language`, `attempted_at` — outreach history. |
| barriers | `id`, `episode_id`, `category`, `description`, `status`, `resolved_at` — access obstacles. |
| tasks | `id`, `episode_id`, `assigned_to`, `type`, `priority`, `status`, `due_at` — accountable staff work. |
| test_completions | `id`, `test_order_id`, `completed_at`, `verification_source` — verification separate from an order. |
| reports | `id`, `episode_id`, `received_at`, `reference`, `recorded_by` — report receipt metadata. |
| clinical_reviews | `id`, `episode_id`, `clinician_id`, `acknowledged_at`, `next_action` — clinician acknowledgement and plan. |
| audit_events | `id`, `actor_id`, `episode_id`, `event_type`, `timestamp`, `metadata` — attributable event trail. |

Avoid storing unnecessary health details or call recordings. Data minimization, retention, access rules, and whether reports are uploaded or represented only by metadata need to be decided before real use.

## Proposed REST routes

| Method and route | Purpose |
| --- | --- |
| `POST /api/v1/patients` | Register a patient. |
| `GET /api/v1/patients` | Search/list patients. |
| `POST /api/v1/episodes` | Create an episode and its test orders. |
| `GET /api/v1/episodes` | List/filter episodes. |
| `POST /api/v1/episodes/:id/contacts` | Log a consented contact outcome. |
| `POST /api/v1/episodes/:id/barriers` | Record a practical barrier. |
| `PATCH /api/v1/tasks/:id` | Assign or resolve a coordinator task. |
| `POST /api/v1/test-orders/:id/completion` | Record verified test completion. |
| `POST /api/v1/episodes/:id/report` | Record report receipt. |
| `POST /api/v1/episodes/:id/review` | Record clinician acknowledgement and next action. |
| `POST /api/v1/episodes/:id/close` | Validate closure prerequisites. |
| `GET /api/v1/dashboard/summary` | Read operational counts. |
| `GET /api/v1/episodes/:id/audit` | Read the audit timeline. |
| `POST /api/v1/voice/events` | Receive a validated voice workflow event during integration phase. |

## Write-path rules

All writes should validate input, caller role, episode ownership/assignment where relevant, consent where outreach is involved, and allowed state transitions. An event from a voice integration is a proposed observation, not authority to close an episode or take a clinical action. The application service enforces closure prerequisites.
