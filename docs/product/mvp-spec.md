# Renal Relay MVP product specification

| Field | Value |
| --- | --- |
| Status | User-confirmed prototype baseline; clinical approval and open decisions remain outstanding |
| Version | 0.1 |
| Date | 4 November 2026 |
| Primary user | Doctor and care team |
| Initial use case | Follow-through for clinician-ordered kidney monitoring in people with diabetes |
| Initial tests | uACR and eGFR |
| Product truth | See [`../../PRODUCT.md`](../../PRODUCT.md) |

## 1. Purpose and product definition

Renal Relay is a clinician-led follow-through workspace that keeps a kidney-monitoring episode visible until the ordered test is verified, report receipt is recorded, a clinician acknowledges the report, and a next action is documented. It addresses the operational gap between ordering a test and recording what the care team did with its result.

The product measures verified care continuity. Message delivery and call completion are intermediate events, not episode success.

## 2. Target outcomes

- Care-team members can see every open episode, its owner, current stage, due window, barrier, and next task.
- Patients can receive consented, plain-language outreach in an agreed language and request human help or opt out.
- Staff can separately verify test completion and report receipt.
- Clinicians acknowledge the report and record a next action before standard closure.
- Every material action has an attributable audit event.

## 3. Prototype scope

### Included in core workflow

- Manual patient and order entry; CSV import is a proposed P1 extension.
- uACR/eGFR monitoring episodes with due window, owner, contact consent, and approved outreach configuration.
- Clinician/coordinator dashboard, episode worklist, episode detail, coordinator tasks, clinician review, and audit timeline.
- Contact and barrier logging, test-completion verification, report receipt, clinician acknowledgement, next action, and guarded closure.
- Fictional or appropriately de-identified demonstration records.

### Deferred extensions

- Multilingual telephony and constrained structured extraction.
- Optional patient mobile page.
- CSV import and operational analytics.
- EHR, laboratory, or other institutional integrations.

### Explicit non-goals

- Diagnosis, medical risk calculation, interpretation of uACR/eGFR results, medication selection, or treatment recommendations.
- Autonomous clinical decisions or episode closure based solely on a call or model output.
- Real-patient deployment before clinical co-design and institutional privacy/security review.

## 4. Users and permissions

| Role | Core actions |
| --- | --- |
| Clinician | Approve orders and due windows, approve outreach/escalation configuration, review reports, document next action, authorize closure. |
| Care coordinator | Work assigned queue, log contact and access barriers, verify test completion, record report receipt, resolve or escalate tasks. |
| Patient | Receive consented contact, share arrangement status or barriers, request callback/help, opt out. |
| Administrator | Manage demo users, access roles, approved templates, and audit visibility. |

Authorization is enforced by the service that owns the workflow state, not by the interface or voice agent.

## 5. Information architecture

| Surface | Primary user goal |
| --- | --- |
| Dashboard | Identify open, overdue, awaiting-report, and awaiting-review work. |
| Episode worklist | Search and filter episodes by stage, owner, and due date. |
| New episode | Select a patient, tests, due window, owner, consent, and approved follow-up setup. |
| Episode detail | Understand current stage, timeline, evidence, barriers, tasks, and audit history. |
| Coordinator queue | Resolve barriers, callbacks, missing reports, and escalations. |
| Clinician review queue | Acknowledge received reports and document next actions. |
| Optional patient page | View plain-language status and ask for help or callback. |
| Settings | Manage approved scripts, escalation rules, roles, and supported languages. |

## 6. Core journey

1. Clinician creates an episode for the required test(s), sets a due window and accountable owner.
2. Care team records channel consent and preferred language before outreach.
3. Patient is contacted using approved wording; an arrangement status or practical barrier is captured.
4. Coordinator resolves barriers and documents the action.
5. Staff separately record verified test completion and report receipt, including provenance.
6. Clinician reviews the report, acknowledges it, and documents the next action.
7. Service validates required evidence before standard closure.

See [`../workflow.md`](../workflow.md) for states and exception handling.

## 7. Acceptance criteria

| ID | Requirement | Acceptance condition |
| --- | --- | --- |
| MVP-01 | Episode creation | Authorized clinician saves patient, test(s), due window, and owner. |
| MVP-02 | Workflow integrity | Only allowed transitions are accepted; closure is blocked while required evidence/documentation is missing. |
| MVP-03 | Team queue | Open/overdue episodes, assigned owner, barrier, and next action are visible and assignable. |
| MVP-04 | Consent and contact | Consent, channel, language, contact result, response, callback, and opt-out are recordable. |
| MVP-05 | Separate evidence | Test completion and report receipt are separate records; verification source is captured. |
| MVP-06 | Clinician sign-off | Standard closure requires clinician acknowledgement and a documented next action. |
| MVP-07 | Audit | Material actions are attributable to a user and visible in an append-only event history. |
| MVP-08 | Safe agent boundary | Agent output cannot interpret results or directly make clinical decisions or close episodes. |

## 8. Success measures

Primary proposed process measure: share of ordered monitoring episodes with test completed, report acknowledged, and next action recorded within the clinician-set window.

Proposed secondary measures: reach/opt-out, barrier type and resolution time, staff time, and escalation rate. The measure definitions, baseline, collection method, and evaluation owner remain to be agreed with the clinical team.

## 9. Product decisions still open

- Clinical co-lead and approval process for scripts, due-window defaults, and escalation rules.
- Demo languages and validation approach.
- Voice provider, authentication approach, and safe patient verification.
- Report representation (metadata, upload, or later integration).
- Data retention, secure storage, audit access, and real-patient privacy/security review.
- Demo date, team capacity, and final implementation scope.

No implementation or clinical outcome is implied by this proposed specification.
