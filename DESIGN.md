---
name: Renal Relay
description: A clear care-team dashboard for kidney monitoring follow-through.
colors:
  ink: "#101210"
  muted-ink: "#48534c"
  paper: "#fbfcfa"
  rule: "#171a17"
  signal-mint: "#0af474"
  relay-green: "#057f3c"
  support-copy: "#344139"
  support-on-dark: "#d7e2da"
  paper-on-dark: "#f4f8f4"
  rule-on-dark: "#68736b"
  support-on-mint: "#153c27"
  row-rule: "#acb2ad"
  overdue: "#a13223"
  offline-red: "#8b3124"
  offline-surface: "#ffb4a8"
  metric-rule: "rgb(16 18 16 / 42%)"
typography:
  display:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "clamp(46px, 6.4vw, 88px)"
    fontWeight: 570
    lineHeight: 0.97
    letterSpacing: "-0.072em"
  headline:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "23px"
    fontWeight: 560
    lineHeight: 1.2
    letterSpacing: "-0.045em"
  body:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "15px"
    fontWeight: 400
    lineHeight: 1.55
  label:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "12px"
    fontWeight: 540
    lineHeight: 1.45
  measurement:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "10px"
    fontWeight: 400
    lineHeight: 1.35
rounded:
  none: "0px"
spacing:
  shell-gutter: "44px"
  cell-inset: "48px"
  mobile-gutter: "16px"
  rule: "1px"
components:
  dark-workflow-panel:
    backgroundColor: "{colors.ink}"
    textColor: "{colors.paper}"
    rounded: "{rounded.none}"
    padding: "clamp(25px, 3.6vw, 48px)"
  signal-overview-panel:
    backgroundColor: "{colors.signal-mint}"
    textColor: "{colors.ink}"
    rounded: "{rounded.none}"
    padding: "clamp(26px, 3.6vw, 48px)"
---

# Design System: Renal Relay

## Overview

**Creative North Star: “Split Studio”**

Renal Relay uses the disciplined contrast of a Swiss editorial grid to make care-team follow-through immediately visible. The dashboard pairs a quiet paper field and oversized, direct headline with an ink-dark workflow panel and a signal-mint overview. Hard rules define the four-panel composition; a small halftone mark nods to the selected Aspen Search reference without importing its marketing content.

The page leads with the purpose of the work, then shows the follow-through sequence, live team counts, and the current episode queue. Aorta’s straightforward task progression informs the workflow panel. The interface renders only values returned by the dashboard API; an unavailable service shows an explicit offline state and no invented records.

**Key Characteristics:**

- Split Studio four-panel grid with square corners and precise black rules.
- Signal mint used for the live count field and small workflow accents.
- Oversized, compact headline balanced by calm operational copy.
- Patient and episode content appears only when returned by the live service.

## Colors

The palette is built from near-black and warm paper, with electric mint as a deliberate signal and deep green for accessible links and state accents.

### Primary

- **Signal Mint** (`{colors.signal-mint}`): The team overview surface and small active marks; use black text on the mint field.
- **Relay Green** (`{colors.relay-green}`): Text links, the live connection indicator, and focus outline where a darker accent is needed.
- **Supporting Text** (`{colors.support-copy}`): Introductory paragraph ink on paper; `support-on-dark` and `paper-on-dark` are the text values for dark surfaces.
- **Supporting Details** (`{colors.support-on-mint}`): Secondary copy on the mint overview. `row-rule` separates episode rows, while `metric-rule` divides counts.
- **Overdue / Offline** (`{colors.overdue}`, `{colors.offline-red}`, `{colors.offline-surface}`): Reserved for overdue status and an unavailable connection, always paired with text.

### Neutral

- **Ink** (`{colors.ink}`): Main text and the dark workflow panel.
- **Muted Ink** (`{colors.muted-ink}`): Supporting copy on paper and mint.
- **Paper** (`{colors.paper}`): Main page field and text on the ink panel.
- **Rule** (`{colors.rule}`): One-pixel divisions between panels.
- **Offline Red** (`{colors.offline-red}`): Connection-unavailable indicator, paired with explicit status text.

**The Live-Data Rule.** Counts and episode rows come from the dashboard response. When it is unavailable, show an unavailable state instead of sample values.

## Typography

**Display Font:** Geist (with Arial, Helvetica, sans-serif fallback)

**Body Font:** Geist (with Arial, Helvetica, sans-serif fallback)

**Label/Mono Font:** Geist Mono for workflow indices; system monospace fallback.

**Character:** Geist gives the page a compact, contemporary voice without competing with the grid. Large display type is tightly tracked; operational text stays open and readable.

### Hierarchy

- **Display** (570, `clamp(46px, 6.4vw, 88px)`; 90px on extra-wide screens; 0.97): Main dashboard statement; on mobile the size scales with viewport and caps at 76px.
- **Headline** (560, 21–23px, 1.2): Workflow, overview, and episode section titles.
- **Body** (400, 13–17px, 1.35–1.65): Explanatory and episode copy; keep paragraphs to roughly 50ch.
- **Label** (400–650, 9–18px, 1.35–1.5): Connection, navigation, workflow index, patient reference, and due-date details.
- **Measurement** (590, 34–56px, 1.08): API-backed summary counts, with tabular numerals.

## Layout

The content is capped at 1440px. A horizontal navigation bar sits above a two-column, two-row grid with zero gaps and shared one-pixel rules. The top row gives equal space to the headline and dark workflow sequence; the lower row pairs mint counts with the episode queue. Desktop shell gutters are 44px and panel insets are fluid, around 25–48px.

At 900px the gutter and panel insets tighten. At 620px and below the panels stack in reading order, the header navigation moves to its own row, and episode details reflow into a two-column list. Keep the first heading visible, keep the four workflow steps distinct, and avoid horizontal page scrolling.

## Elevation & Depth

This system is flat and uses no persistent shadows. Black rules, large color fields, and spacing establish separation; focus uses a visible 3px green outline with a 4px offset. Do not add shadows as a substitute for the grid.

## Shapes

Surfaces and panel corners are square. Use one-pixel borders and crisp geometric marks. The workflow’s dot pattern is a small authored SVG texture tied to the supplied Aspen reference; do not expand it into a generic background grid.

## Components

### Header and Navigation

- **Shape:** Open horizontal bar with a one-pixel bottom rule.
- **Navigation:** Text links with a green underline for the current section; no pills or enclosing cards.
- **Connection:** A small status dot and explicit “Live connection” or “Connection unavailable” text.
- **Responsive:** Brand and status remain on the first row; navigation occupies the second row on narrow screens.

### Workflow Panel

- **Shape:** Square ink field, padded 25–48px.
- **Content:** Four steps in sequence: order placed, test completed, report received, next action recorded.
- **Color:** Paper text, mint indices and directional marks.
- **Behavior:** The diagram is explanatory content, not a claim about an individual episode’s current state.

### Team Overview

- **Shape:** Square signal-mint field divided into a two-by-two count grid.
- **Content:** Open episodes, overdue, awaiting report, and clinician review values from the API; unavailable values render as an em dash.
- **Contrast:** Ink labels and numbers; the due-within-seven-days count is supporting copy.

### Episode Queue

- **Shape:** Open paper field with one-pixel row rules.
- **Content:** Only API-returned patient reference, tests, status, due date, and owner are shown.
- **States:** A connected empty response gets a neutral empty message. A failed request gets a clear error and a working retry control.
- **Status:** Use explicit status text; color is supplementary.

## Do's and Don'ts

### Do:

- Keep the dashboard heading and live connection state easy to find.
- Preserve the split grid, high contrast, and square rules.
- Pair status colors with readable labels and visible keyboard focus.
- Keep clinical review and next-action authority with care teams.
- Show only API-backed patient and episode information.

### Don't:

- Add fictional patients, fabricated counts, or demo banners to the dashboard.
- Add controls that do not perform a real action or links that lead nowhere.
- Use color alone to convey status, or imply automated clinical decisions.
- Add rounded card stacks, decorative gradients, or unsupported clinical interpretations.
