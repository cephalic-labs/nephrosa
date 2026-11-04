# Dashboard surface brief

## Purpose

Help a care coordinator see which ordered uACR/eGFR monitoring episodes remain open, who owns them, when they are due, and what follow-through is expected next.

## Mode and audience

Operate. The primary user is a clinic care-team member reviewing a daily queue at a workstation, with occasional narrow-screen use.

## First viewport

The first view establishes the workspace and current owner, then shows open/overdue/review measures and a table of episodes with patient reference, ordered tests, workflow status, due date, and owner. The episode table is the primary working surface. A distinct notice identifies fictional demonstration data when the page is not connected to the local API.

## Content and state

- Show only minimal patient display information and a prototype reference.
- Keep ordered test types and workflow status explicit.
- Show dates and accountable owner without relying on color or position alone.
- Show the next task and an open barrier when available.
- Provide empty states for no open episodes and no follow-through tasks.
- Keep disabled or unimplemented controls visibly unavailable; do not imply that a displayed action has changed persisted workflow state.

## Interaction boundary

The current prototype is read-only. The episode navigation and view links are in-page anchors, and episode creation is disabled. No authentication, patient search, episode mutation, or clinical decision support is implemented. The care team retains review, next-action, escalation, and closure authority.

## Visual direction

Follow the project-wide [DESIGN.md](../../DESIGN.md): calm light surfaces, neutral-first composition, teal for active navigation and primary actions, readable status text, and restrained warning color. The source proposal's teal/mint palette informs accents, while the Impeccable craft and legibility guidance governs the interface.
