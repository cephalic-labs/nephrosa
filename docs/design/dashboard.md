# Dashboard surface brief

## Purpose

Help a care coordinator see which ordered uACR/eGFR monitoring episodes remain open, who owns them, when they are due, and what follow-through is expected next.

## Mode and audience

Operate. The primary user is a clinic care-team member reviewing a daily queue at a workstation, with occasional narrow-screen use.

## First viewport

The first view establishes the workspace, then shows API-backed open/overdue/review measures and a queue of open episodes with patient reference, ordered tests, workflow status, due date, and owner. Each patient name links to a read-only episode detail route. When the API is unavailable, measures use an em dash and both dashboard and detail views show a clear retryable unavailable state; they never substitute sample records.

## Content and state

- Show only minimal patient display information and a prototype reference.
- Keep ordered test types and workflow status explicit.
- Show dates and accountable owner without relying on color or position alone.
- Episode detail shows recorded tasks, barriers, contact history, completion verification, report receipt metadata, clinician acknowledgement, and audit event labels.
- Provide empty states for no open episodes and no follow-through tasks.
- Keep disabled or unimplemented controls visibly unavailable; do not imply that a displayed action has changed persisted workflow state.

## Interaction boundary

The current prototype is read-only. The dashboard links to episode details; no episode mutations or clinical decision support are implemented. The read API has no authentication or authorization and is for fictional local data only. The care team retains review, next-action, escalation, and closure authority.

## Visual direction

Follow the project-wide [DESIGN.md](../../DESIGN.md): warm-paper light mode by default, an available dark theme, editorial serif display, lime and cyan work surfaces, readable status text, and restrained warning color.
