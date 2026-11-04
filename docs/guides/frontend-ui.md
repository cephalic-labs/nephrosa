# Frontend UI development

Use the Impeccable skill for Renal Relay interface work. The repository's current `app/page.tsx` is the unmodified Next.js starter screen; its tokens and styling are inherited setup, not an approved product identity. Do not treat those choices as the Renal Relay visual direction.

## Required workflow

1. Read [`../../PRODUCT.md`](../../PRODUCT.md) and the relevant requirements in [`../product/mvp-spec.md`](../product/mvp-spec.md).
2. Load Impeccable project context from the repository root with `impeccable context` and pass the target route/file when working on a surface.
3. For a new screen, establish its user, task, states, content, and responsive behavior before implementation. Keep visual direction open until the product owner selects it.
4. Before editing UI code, follow the applicable Impeccable reference and read its `craft-floor.md` guidance after the direction is settled.
5. Inspect the result at desktop and mobile widths in one bounded visual review; address the discovered defects in a single correction pass and confirm once.
6. Run the Impeccable detector once on changed UI targets after implementation. Record actual verification evidence in the handoff.

## Product-specific constraints

- Clinical actions, patient identity, due dates, barriers, report receipt, acknowledgement, and closure status must be represented clearly and accurately.
- Color cannot be the only status signal. Use text and icons as well.
- Preserve keyboard access, visible focus, readable tables, responsive layout, and accessible empty/error states.
- Do not invent clinical outcomes, patient quotes, analytics, or claims. Mark demonstration data as fictional.
- Patient-facing outreach is consented and plain-language; clinical questions route to care staff.

## Current design record

The current visual implementation is being documented separately in root `DESIGN.md` as an incumbent starter reference. It must not be mistaken for an approved Renal Relay design system. A later visual direction should be established through Impeccable before the application UI is replaced.
