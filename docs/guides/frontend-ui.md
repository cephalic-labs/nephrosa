# Frontend UI guidance

Use the Impeccable skill for interface work. Read the relevant skill references before shaping a new surface, then ground implementation in [PRODUCT.md](../../PRODUCT.md), [DESIGN.md](../../DESIGN.md), and the route-specific brief under `docs/design/`.

## Current system

The dashboard at `app/page.tsx` is the first Renal Relay surface. `app/globals.css` holds its tokens and components. The design is neutral-first with teal interaction accents, clear operational hierarchy, visible synthetic-data labeling, semantic tables, and responsive behavior. Do not reintroduce the Next.js starter page.

## UI implementation expectations

- Keep server data access on the server and never place secrets in `NEXT_PUBLIC_*` variables.
- Follow the installed Next.js guides under `node_modules/next/dist/docs/` as required by `AGENTS.md`.
- Preserve clinical authority boundaries: the UI displays workflow state but does not interpret laboratory results or recommend treatment.
- Use semantic landmarks and table headings, readable text status, keyboard focus, and narrow-screen behavior.
- Document only durable design decisions that are present in the rendered implementation.
- Keep prototype and fictional-data limitations visible where data could be mistaken for live records.
