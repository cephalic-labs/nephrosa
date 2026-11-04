# Episode detail surface brief

## Purpose and mode

Operate. Help a care-team member inspect one open monitoring episode and verify what has been recorded so far.

## Content and hierarchy

- Identify the patient by the same minimum display information and prototype reference already shown in the dashboard.
- Show clinician-ordered tests with due dates and separately recorded completion verification.
- Show coordinator tasks, barriers, and contact history when present.
- Show report receipt metadata, clinician acknowledgement, next action, and audit event labels when present.
- Never display report contents, laboratory values, generated clinical interpretation, or invented records.

## States and boundaries

- A missing or closed episode uses the not-found state.
- An unavailable API uses an explicit retryable error state, not empty sections.
- Each absent record type has a plain-language empty message.
- This surface is read-only and uses the shared theme, header, and responsive rules in [DESIGN.md](../../DESIGN.md).
- The current read API has no authentication or authorization and is restricted to fictional local prototype data.
