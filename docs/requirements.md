# Requirements and acceptance criteria

This document translates the supplied Renal Relay concept into prototype requirements. Priority labels reflect the source proposal and are not delivery commitments.

## P0 — core workflow

| Capability | Acceptance criteria |
| --- | --- |
| Monitoring episode creation | An authorized clinician can create an episode with patient, uACR/eGFR order(s), due window, owner, and approved follow-up configuration. |
| Episode lifecycle | Allowed state changes are enforced; a standard episode cannot close without required test completion evidence, report receipt, clinician acknowledgement, and a documented next action. |
| Coordinator work queue | Open work, overdue episodes, owners, barriers, and next actions are visible and assignable. |
| Contact and barrier records | Consent status, contact attempts, patient responses, language, and structured practical barriers can be recorded. |
| Test and report tracking | Test completion and report receipt are represented as separate facts and can each be verified by staff. |
| Clinician sign-off | A clinician can acknowledge a received report and record the next action before standard closure. |
| Roles and audit | Actions are attributable to an authorized user and recorded in an append-only event history. |

## P1 — extensions

| Capability | Acceptance criteria |
| --- | --- |
| Voice follow-up | An approved, consent-aware call flow captures structured follow-up outcomes and hands medical questions to staff. |
| CSV import | Fictional demo patients and orders can be loaded without an EHR integration. |
| Optional patient page | A patient can see a plain-language status/help experience with an agreed safe authentication flow. |
| Operational analytics | Episode completion, delay, reach/opt-out, barrier resolution, staff time, and escalation measures can be summarized. |

## Safety and privacy requirements

- Only consented channels may be used for outreach; opt-out and callback requests must be represented.
- Outreach wording and escalation rules are clinician-approved before real-patient use.
- The agent may capture arrangement status and barriers, but must not interpret results, diagnose, calculate risk, or recommend treatment.
- Clinical questions and urgent concerns follow a human escalation process approved by the clinical team.
- The demo uses fictional or appropriately de-identified records; minimize sensitive data and avoid unnecessary call recordings.
- Access control, secure storage, retention, and audit requirements need institutional review before a real pilot.

## Exclusions for the initial prototype

- EHR or laboratory system integration; initial inputs are manual entry or CSV.
- Automated laboratory result interpretation or treatment guidance.
- Autonomous closure of an episode without clinician acknowledgement and next-action documentation.
- A claim that the proposed pilot or any clinical outcome has already been achieved.

## Definition of a successful episode

The episode has the required test completion evidence, separately recorded report receipt, clinician acknowledgement, and a documented next action. A completed contact attempt alone is insufficient.
