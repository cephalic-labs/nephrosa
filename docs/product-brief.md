# Renal Relay product brief

## Summary

Renal Relay is a clinician-led follow-through system for kidney monitoring in people with diabetes. It tracks an ordered uACR and/or eGFR monitoring episode across the gap between visits: patient contact, test arrangement and completion, report receipt, clinician acknowledgement, and documented next action.

The differentiator is verified care continuity. A sent reminder or answered call does not count as a completed episode.

## Problem

An order alone does not establish that a patient reached a lab, that the report returned to the care team, or that a clinician reviewed it and documented a next step. Practical barriers can include time away from work, travel, cost, unclear instructions, and language access. The care team needs a shared view of the unfinished step, its owner, and the evidence needed to move forward.

## Users

| Role | Primary responsibilities |
| --- | --- |
| Clinician | Approve orders, due windows, outreach content, and escalation rules; review reports; document next actions; authorize closure. |
| Care coordinator | Work the follow-up queue, resolve access barriers, document contact and report receipt, and escalate exceptions. |
| Patient | Receive consented, understandable communication; share arrangement status or barriers; request human help or opt out. |
| Administrator | Manage prototype users, permissions, approved templates, and audit visibility. |

## Proposed experience

Each monitoring order becomes a persistent episode with an owner, due window, approved outreach content, and escalation rules. A care-team workspace surfaces open and overdue work. A consented multilingual voice flow can ask about test arrangements and practical barriers, then produce constrained structured follow-up information. Staff verify completion and report receipt. The clinician reviews the report and records the next action before standard closure.

## Prototype screens

- Dashboard with episode counts, overdue work, and shortcuts.
- Searchable episode worklist with status, owner, and due date.
- Episode creation and patient lookup.
- Episode detail with stage tracker, timeline, tasks, and audit history.
- Coordinator queue for barriers, callbacks, and exceptions.
- Clinician review queue for received reports and next-action documentation.
- Optional patient page for plain-language status, help, and callback requests.
- Settings for approved scripts, escalation rules, roles, and languages.

## Design direction

Use a calm, professional interface with deep teal `#0E343B`, teal `#087F78`, soft mint `#EAF7F5`, neutral `#F6F8FA`, and amber for attention states. Pair status colors with text and icons. Support readable tables, search, filters, keyboard access, responsive layouts, clear empty/error states, and large touch targets for the optional patient experience.

Candidate languages include Tamil, Hindi, and English; validate availability and quality before enabling them in a demo.

## Demonstration scenario

A fictional patient has uACR and eGFR orders. During a consented call, the patient reports that work prevents a lab visit. The care team sees a time-off-work barrier, coordinates an available lab window, records completion and report receipt, then the clinician acknowledges the report and records the next action. The episode timeline and audit history show each handoff.

## Scope boundary

The prototype is for workflow demonstration with fictional or appropriately de-identified records. It does not interpret laboratory results or provide medical advice. A real patient pilot requires clinical co-design, approved scripts and escalation rules, and institutional privacy and security review.
