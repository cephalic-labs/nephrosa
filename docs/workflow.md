# Episode workflow and safeguards

## Standard path

1. **Order:** A clinician records the patient, uACR and/or eGFR orders, due window, responsible owner, approved outreach content, and escalation rule.
2. **Consent and contact:** The team confirms communication consent and language preference. Approved outreach asks whether the test is arranged or completed and offers human help.
3. **Barrier resolution:** A reported obstacle, such as cost, travel, work schedule, or unclear instructions, becomes a coordinator task with an owner and next action.
4. **Completion verification:** Staff record when testing was completed and how completion was verified.
5. **Report receipt:** Staff separately record that the report is available to the care team. A completed test does not imply a received report.
6. **Clinical acknowledgement:** A clinician reviews the report outside the agent, acknowledges it in the system, and records the next action.
7. **Closure:** Standard closure is allowed only when required evidence and documentation are present. An exception route must be defined and human-approved.

## Suggested states

`ordered → outreach_pending → awaiting_booking / booked → awaiting_test_completion → awaiting_report → awaiting_clinician_review → closed`

An open barrier, overdue condition, opt-out, or escalation can be represented as a task or flag alongside the episode's main state. The implementation should avoid forcing these concurrent concerns into mutually exclusive states.

## Closure guard

Standard closure requires all of the following:

- Required test completion has been verified and its source recorded.
- Report receipt has been recorded separately.
- A clinician has acknowledged the report.
- The clinician has documented a next action.

A successful call, patient statement, reminder delivery, or elapsed time cannot substitute for missing evidence. Any exception closure path needs an explicit human decision and audit record.

## Agent boundary and escalation

The voice/LLM layer may propose structured fields for scheduling status or practical barriers. The application service must validate consent, input shape, user permissions, and permitted state changes before recording them. The agent must not interpret uACR/eGFR, diagnose, calculate medical risk, or recommend medication or treatment. Medical questions or urgent concerns go to a clinical human under an approved escalation process.

## Audit expectations

Record who performed each material action and when, including order/configuration changes, consent and contact outcomes, task assignment/resolution, verification source, report receipt, clinician acknowledgement, next action, escalation, and closure. The detailed event schema and retention controls remain design decisions.
