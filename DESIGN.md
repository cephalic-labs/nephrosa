---
name: Renal Relay
description: A care-team dashboard with the expressive contrast of a type studio.
colors:
  night: "#050505"
  paper: "#f5f2ec"
  paper-muted: "#d4d0c8"
  ink: "#141411"
  muted-ink: "#a9a69f"
  line: "#45443f"
  header-line: "#363631"
  lime: "#8cff00"
  lime-hover: "#78da00"
  lime-ink: "#34421f"
  lime-rule: "rgb(20 20 17 / 40%)"
  cyan: "#20c7d9"
  cyan-ink: "#19383b"
  cyan-muted: "#315356"
  cyan-rule: "rgb(20 20 17 / 42%)"
  critical: "#a42c24"
  success: "#18753c"
  offline-line: "#514d42"
  offline-ground: "#1b1a17"
  offline-ink: "#f6c273"
  episode-muted: "#55524c"
  episode-ink: "#282722"
  episode-rule: "#c5c0b7"
  header-ground: "rgb(5 5 5 / 94%)"
  light-ground: "#f1eee8"
  light-header: "rgb(241 238 232 / 94%)"
  light-ink: "#141411"
  light-muted: "#55524c"
  light-rule: "#c9c4ba"
  light-offline-ground: "#fff8e9"
  light-offline-ink: "#68440e"
  focus-dark: "#467500"
typography:
  display:
    fontFamily: "Noto Serif Relay, Georgia, Times New Roman, serif"
    fontSize: "clamp(58px, 7.2vw, 96px)"
    fontWeight: 450
    lineHeight: 0.99
    letterSpacing: "-0.058em"
  detail-display:
    fontFamily: "Noto Serif Relay, Georgia, Times New Roman, serif"
    fontSize: "clamp(46px, 6vw, 76px)"
    fontWeight: 450
    letterSpacing: "-0.05em"
  detail-unavailable:
    fontFamily: "Noto Serif Relay, Georgia, Times New Roman, serif"
    fontSize: "clamp(38px, 5vw, 58px)"
    fontWeight: 450
  detail-mobile:
    fontFamily: "Noto Serif Relay, Georgia, Times New Roman, serif"
    fontSize: "clamp(44px, 11vw, 66px)"
    fontWeight: 450
  detail-section:
    fontFamily: "Noto Serif Relay, Georgia, Times New Roman, serif"
    fontSize: "24px"
    fontWeight: 480
  headline:
    fontFamily: "Noto Serif Relay, Georgia, Times New Roman, serif"
    fontSize: "clamp(26px, 3vw, 37px)"
    fontWeight: 480
    lineHeight: 1.15
    letterSpacing: "-0.045em"
  body:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "15px"
    fontWeight: 400
    lineHeight: 1.55
  label:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "12px"
    fontWeight: 550
    lineHeight: 1.35
  brand:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "18px"
    fontWeight: 650
    letterSpacing: "-0.055em"
  navigation:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "13px"
    fontWeight: 550
  navigation-mobile:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "12px"
    fontWeight: 550
  hero-copy:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "16px"
    fontWeight: 400
    lineHeight: 1.65
  compact-copy:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "14px"
    fontWeight: 400
  workflow-heading:
    fontFamily: "Noto Serif Relay, Georgia, Times New Roman, serif"
    fontSize: "30px"
    fontWeight: 480
  workflow-heading-mobile:
    fontFamily: "Noto Serif Relay, Georgia, Times New Roman, serif"
    fontSize: "27px"
    fontWeight: 480
  workflow-copy:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "13px"
    fontWeight: 400
  workflow-index:
    fontFamily: "Geist Mono, monospace"
    fontSize: "10px"
    fontWeight: 400
  overview-heading:
    fontFamily: "Noto Serif Relay, Georgia, Times New Roman, serif"
    fontSize: "29px"
    fontWeight: 480
  overview-heading-mobile:
    fontFamily: "Noto Serif Relay, Georgia, Times New Roman, serif"
    fontSize: "25px"
    fontWeight: 480
  measure-mobile:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "43px"
    fontWeight: 590
  episode-count:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "19px"
    fontWeight: 400
  episode-copy:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "11px"
    fontWeight: 400
  episode-state:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "17px"
    fontWeight: 590
  workflow-label-small:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "10px"
    fontWeight: 550
  hero-tablet:
    fontFamily: "Noto Serif Relay, Georgia, Times New Roman, serif"
    fontSize: "clamp(54px, 7vw, 76px)"
    fontWeight: 450
  hero-mobile:
    fontFamily: "Noto Serif Relay, Georgia, Times New Roman, serif"
    fontSize: "clamp(54px, 13.5vw, 76px)"
    fontWeight: 450
  hero-narrow:
    fontFamily: "Noto Serif Relay, Georgia, Times New Roman, serif"
    fontSize: "52px"
    fontWeight: 450
  episodes-heading-mobile:
    fontFamily: "Noto Serif Relay, Georgia, Times New Roman, serif"
    fontSize: "28px"
    fontWeight: 480
  measure:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "clamp(36px, 4.4vw, 57px)"
    fontWeight: 590
    lineHeight: 1.02
    letterSpacing: "-0.065em"
rounded:
  control: "4px"
  status: "5px"
  panel: "12px"
spacing:
  shell-gutter: "48px"
  header-height: "80px"
  shell-max: "1440px"
  workflow-inset: "32px"
  episode-inset: "clamp(24px, 3.8vw, 48px)"
components:
  workflow-panel:
    backgroundColor: "{colors.cyan}"
    textColor: "{colors.ink}"
    rounded: "{rounded.panel}"
    padding: "32px"
  overview-panel:
    backgroundColor: "{colors.lime}"
    textColor: "{colors.ink}"
    rounded: "{rounded.panel}"
    padding: "32px clamp(24px, 4vw, 54px) 29px"
  episode-panel:
    backgroundColor: "{colors.paper}"
    textColor: "{colors.ink}"
    rounded: "{rounded.panel}"
    padding: "clamp(24px, 3.8vw, 48px)"
  button-retry:
    backgroundColor: "{colors.lime}"
    textColor: "{colors.ink}"
    rounded: "{rounded.control}"
    height: "43px"
    padding: "0 17px"
---

# Design System: Renal Relay

## Overview

**Creative North Star: “Electric Type Studio”**

This system adapts Glyphs.app’s expressive type-studio character to a practical care-team workspace: oversized serif display, crisp sans-serif interface text, and confident blocks of electric lime and cyan. A compact header switch starts in light mode and lets staff move to the original night-black canvas. The light theme uses warm paper ground, dark ink, and the same tactile accent work surfaces; both themes preserve the task-first hierarchy.

The dashboard introduces the follow-through path, shows current team counts, and then presents the live episode queue. It does not reproduce Glyphs’ software copy, type specimens, product images, or testimonials. The bundled Noto Serif variable font supplies the display voice; its OFL license is included in `public/fonts/OFL.txt`.

**Key Characteristics:**

- Light warm-paper shell by default, with a switch to the night-black theme and oversized serif display.
- Lime for live workload and connected-state signals; cyan for the workflow sequence.
- A light paper episode surface provides a distinct, readable data zone.
- Only API-returned counts and episode rows are displayed.

## Colors

Near-black, warm paper, electric lime, and bright cyan give the dashboard the contrast and color confidence of the Glyphs reference while keeping status and episode content easy to read.

### Primary

- **Electric Lime** (`{colors.lime}`): The team overview surface, connected state, keyboard focus, and retry action. Keep its text near-black.
- **Studio Cyan** (`{colors.cyan}`): The workflow panel, with near-black text and a dark-cyan secondary label.

### Secondary

- **Urgency Red** (`{colors.critical}`): The overdue status marker, accompanied by the word “Overdue.”
- **Success Green** (`{colors.success}`): A non-overdue episode status marker; status text remains present.

### Neutral

- **Night** (`{colors.night}`): The page ground and sticky header.
- **Warm Paper** (`{colors.paper}`) and **Soft Paper** (`{colors.paper-muted}`): Queue surface and supporting text on night.
- **Queue Ink** (`{colors.episode-ink}`) and **Queue Muted Ink** (`{colors.episode-muted}`): Patient and status text on the warm-paper surface.
- **Ink** (`{colors.ink}`): Text on the light, lime, and cyan surfaces.
- **Muted Ink** (`{colors.muted-ink}`): Supporting text on the night surface.
- **Rules** (`{colors.line}`, `{colors.header-line}`, `{colors.episode-rule}`): Section and row divisions.
- **Connection Unavailable** (`{colors.offline-ground}`, `{colors.offline-ink}`): Explicit degraded-service state, distinct from the lime live state.

**The Live-Data Rule.** Every summary measure and episode row comes from the dashboard API. An unavailable service never falls back to sample records or fabricated counts.

## Typography

**Display Font:** Noto Serif Relay (bundled variable font; Georgia and Times New Roman fallbacks)

**Body Font:** Geist (with Arial and Helvetica fallbacks)

**Label/Mono Font:** Geist for interface labels; Geist Mono for workflow indices.

**Character:** A high-contrast serif headline carries the editorial presence of a type studio. Geist keeps navigation, labels, and patient information direct and familiar.

### Hierarchy

- **Display** (450, `clamp(58px, 7.2vw, 96px)`, 0.99): Dashboard headline on wide layouts; scales down at tablet and phone sizes.
- **Headline** (480, 26–37px, 1.15): Episode and panel headings in the display serif.
- **Body** (400, 14–16px, 1.55–1.65): Explanatory and operational copy; limit hero text to about 48ch.
- **Label** (400–700, 10–13px, 1.35): Navigation, status, workflow labels, and episode metadata.
- **Measure** (590, 36–57px, 1.02): API-backed workload counts with tabular numerals.

**The Contrast Rule.** Keep the large serif voice in headings. Use Geist for all labels, controls, and episode data.

## Layout

The page is a single centered workspace capped at 1440px, with 48px desktop gutters and a sticky, compact header. The header includes an accessible light/dark switch. The first content region pairs the oversized headline and introduction with a cyan workflow panel. A lime band contains five live summary measures; a warm-paper panel holds the episode queue and its loading, unavailable, and empty states. Light mode uses warm paper for the shell and near-white for the queue; dark mode restores the night shell.

At 1020px the metrics move to three columns and the episode rows wrap. At 700px the hero and workflow stack, the metric band becomes two columns, and episode rows reflow into compact pairs. Keep the page free of horizontal overflow and preserve the nav and service status in the sticky header.

## Elevation & Depth

Persistent shadows are not used. Large color panels and one-pixel rules define the composition. The sticky header uses a translucent theme-matched ground and a restrained 14px backdrop blur so it remains legible over the long page; focus uses a 3px high-contrast outline with a 4px offset.

## Shapes

Feature and episode panels use gently rounded corners (12px). The connection indicator is a compact, slightly squared control (5px); the retry button is more square (4px). Rules remain 1px. Avoid round status pills and nested card stacks.

## Components

### Navigation

- **Style:** A sticky theme-matched header with a fine divider, high-contrast wordmark, and plain sans-serif links.
- **Active state:** Lime underline under the current section; links shift to the theme's primary ink on hover.
- **Connection state:** Lime field when live; dark inset field with amber text when unavailable.
- **Responsive:** Navigation moves to its own row on narrow screens while the service state stays visible.
- **Theme switch:** A text-labeled, keyboard-operable control switches between light and dark themes; light is the initial theme, and the chosen theme is saved in browser storage.

### Workflow Panel

- **Shape:** Rounded cyan field (12px) with near-black text.
- **Content:** Four numbered steps in a two-column sequence: order, test completion, report receipt, and recorded next action.
- **Behavior:** Describes the product’s workflow; it does not imply that an individual episode has completed these stages.

### Team Overview

- **Shape:** Lime panel (12px) with five measures separated by fine dark rules.
- **Content:** Open, overdue, awaiting-report, clinician-review, and due-within-seven-days values from the API. If unavailable, use an em dash and “Awaiting live data.”
- **States:** Connected and unavailable states use distinct backgrounds and explicit labels.

### Episode Queue

- **Shape:** Warm-paper panel (12px), with one-pixel row rules and dark text.
- **Content:** Patient display name and reference, ordered tests, status, due date, and owner from the API only.
- **States:** Show a real empty state for a successful zero-row response; show a clear service error and a working retry action after a failed request.
- **Status:** Use explicit status text; red and green dots are supplementary.

### Episode Detail

- **Purpose:** Follow one open episode's recorded path across tests, coordinator follow-through, contacts, reports, clinician review, and audit history.
- **Layout:** Editorial serif patient heading over two responsive columns of ruled sections; stack to one column on phones.
- **Content:** Show only API-returned data. Distinguish recorded evidence from missing records with explicit copy. Show report reference metadata only; never display or interpret laboratory values.
- **Unavailable / missing:** A missing episode returns the not-found route. If the service is unavailable, show an explicit retry action. Do not replace service failure with empty records.
- **Navigation:** Episode names in the dashboard queue open this detail view; the shared header returns to dashboard sections.

### Retry Button

- **Shape:** Squared lime control (4px), 43px high.
- **Hover / Focus:** A small lift and darker lime on hover; a clear lime outline on keyboard focus.
- **Action:** Submits a GET navigation to reload the dashboard response.

## Do's and Don'ts

### Do:

- Keep the display serif for section and hero headings; keep product controls and patient data in Geist.
- Use lime and cyan as broad, deliberate work surfaces, not as gradients.
- Keep each workflow step readable on desktop and mobile.
- Pair every colored state with text and preserve keyboard focus visibility.
- Render only live API data; distinguish service failure from a successful empty response.

### Don't:

- Add fictional names, sample counts, fabricated episodes, or demonstration banners.
- Use marketing testimonials, mock product imagery, or controls unrelated to care-team work.
- Hide service failure behind zero counts or an empty-queue message.
- Imply diagnosis, interpretation of lab results, or automated clinical decisions.
