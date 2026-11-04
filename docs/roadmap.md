# Roadmap and open decisions

## Proposed phases

### Phase 1 — Foundation

Establish monorepo conventions, local PostgreSQL setup and migrations, Go API baseline, web design system, and fictional seed data.

### Phase 2 — Core workflow

Implement patient and episode management, allowed state transitions, coordinator queue, test and report verification, clinician acknowledgement, next-action capture, and audit history.

### Phase 3 — Patient and voice integration

After approval, add clinician-reviewed scripts, consent checks, multilingual telephony adapter, constrained structured extraction, callbacks, and escalation handling.

### Phase 4 — Quality and demo

Address role and workflow edge cases, accessibility, operational measures, deployment, and a rehearsed end-to-end demonstration.

The phases are proposed sequencing, not a committed schedule. The demo timeline and available developers have not been supplied.

## Pilot proposal

The concept proposes a 60–90-day pilot in one diabetes clinic. Agree on the baseline with the clinical team and review exceptions weekly. The primary process outcome is the share of ordered kidney-monitoring episodes with testing completed, report acknowledged, and a next action recorded within the clinician-set window.

Secondary measures proposed are patient reach/opt-out, barrier categories and resolution time, staff time, and escalation rate. These measures require operational definitions and an approved evaluation plan.

## Open decisions

- Confirm the clinical co-lead and clinician review of scripts and escalation criteria.
- Select the initial demo languages and validate speech quality and availability.
- Select a telephony provider; assess consent flow, costs, and supported calling regions.
- Choose authentication, roles, patient verification, and account lifecycle.
- Decide whether report handling is manual acknowledgement, metadata entry, or document upload.
- Define data retention, storage, audit access, and institutional privacy/security review before any real-patient use.
- Confirm the demo date, team capacity, and scope before estimating delivery.
- Define an exception closure policy and the accountable human role.
