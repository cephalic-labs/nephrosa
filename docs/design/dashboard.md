# Dashboard surfaces

## Public landing page (`/`)

Introduce Renal Relay before sign-in. Lead with the user-supplied care-team photograph, explain kidney-monitoring follow-through in plain language, and provide a direct route to care-team sign-in. Keep clinical accountability explicit. Do not show patient data, fabricated summary measures, or workflow controls publicly.

## Role workspaces (`/dashboard/clinician`, `/dashboard/coordinator`)

Require an authenticated care-team session. `/dashboard` routes each signed-in user to their role-specific workspace. Adapt the welcome copy and focus summary for clinicians and care coordinators, then present live API workload measures and a clearly labeled episode queue. The queue identifies the patient by minimum display name and prototype reference, ordered tests, current workflow step, due date, and care owner. The patient link opens episode detail. The clinician can create an episode; both roles can inspect workflow and episode history.

Unavailable API data uses an explicit retryable state and never substitutes example rows or counts. Closed history is selected through an explicit history control.

## Navigation and workflow

The shared care-team navigation links to `/dashboard`, `/workflow`, and the episode queue. The workflow page explains the ordered handoffs from clinician order through verified testing, report receipt, clinician acknowledgement, and closure. Every visible navigation item must lead to a functioning route or section.

## Visual and content rules

- Follow [DESIGN.md](../../DESIGN.md): warm, light-first editorial styling, responsive layouts, and an available dark theme.
- Keep measures and queue information legible at narrow widths; use clear column labels on wide screens and stacked, labeled episode rows on mobile.
- Distinguish due status with text and a marker, not color alone.
- Keep clinical decisions with the clinician. Use fictional data only.
