---
name: Renal Relay
description: A calm, legible care-team workspace for kidney monitoring follow-through.
colors:
  ink: "#18353a"
  ink-muted: "#65787b"
  teal: "#087f78"
  teal-dark: "#086660"
  teal-pale: "#eaf5f3"
  canvas: "#f5f7f7"
  surface: "#ffffff"
  line: "#e4eaea"
  line-strong: "#d6e0df"
  warning: "#a65b16"
  warning-bg: "#fff4e7"
  focus-ring: "#39a99d"
typography:
  interface:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "14px"
    fontWeight: 400
    lineHeight: 1.5
  page-title:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "clamp(24px, 2.3vw, 30px)"
    fontWeight: 600
    lineHeight: 1.2
    letterSpacing: "-1px"
  section-title:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "17px"
    fontWeight: 600
    lineHeight: 1.35
    letterSpacing: "-0.35px"
  metadata:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "11px"
    fontWeight: 400
    lineHeight: 1.45
  compact:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "10px"
    fontWeight: 600
    lineHeight: 1.4
  small:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "12px"
    fontWeight: 400
    lineHeight: 1.5
  wordmark:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "20px"
    fontWeight: 700
    lineHeight: 1.2
  metric:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "30px"
    fontWeight: 600
    lineHeight: 1.1
  mobile-title:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "24px"
    fontWeight: 600
    lineHeight: 1.2
  control-symbol:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
    fontSize: "18px"
    fontWeight: 400
    lineHeight: 1.5
rounded:
  xs: "4px"
  sm: "8px"
  md: "12px"
  lg: "16px"
  pill: "999px"
  circle: "50%"
spacing:
  sidebar-width: "248px"
  header-height: "70px"
  content-max: "1240px"
  control-height: "40px"
components:
  nav-item:
    backgroundColor: "{colors.teal-pale}"
    textColor: "{colors.teal-dark}"
    rounded: "{rounded.sm}"
    height: "42px"
    padding: "0 10px"
  button-primary:
    backgroundColor: "{colors.teal}"
    textColor: "{colors.surface}"
    rounded: "{rounded.sm}"
    height: "{spacing.control-height}"
    padding: "0 14px"
  data-surface:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.ink}"
    rounded: "{rounded.md}"
    padding: "18px"
  status-chip:
    backgroundColor: "{colors.teal-pale}"
    textColor: "{colors.teal-dark}"
    rounded: "{rounded.pill}"
    padding: "5px 8px"
---

## Overview

Renal Relay is a care-team operations workspace. Its visual system gives episode status, due dates, ownership, and the next accountable step a clear reading order. The interface is intentionally quiet so exceptions and unresolved follow-through remain easy to find. PRODUCT.md defines product truth; this file records the system visible in the first implemented dashboard.

## Colors

- **Ink** (`{colors.ink}`) is the primary text color; muted ink is reserved for supporting information and remains readable on white.
- **Teal** (`{colors.teal}`) marks primary actions and active navigation. Deep teal is used for text on pale teal surfaces.
- **Pale teal** (`{colors.teal-pale}`) identifies selected workspace areas and positive workflow states.
- **Canvas** (`{colors.canvas}`) separates the work area from white data surfaces. Borders use the two low contrast line colors to define table and surface boundaries.
- **Warning** (`{colors.warning}`) and warning background are reserved for overdue or time-sensitive work. Status meaning is also written as text; color never carries meaning alone.

The page is light by design for a clinic workstation. Do not introduce dark mode without checking the information density, contrast, and ambient-use context of the full workspace.

## Typography

Geist is bundled locally through the `geist` package and loaded with `next/font/local`, avoiding a runtime or build-time Google Fonts request. Its SIL Open Font License is included with the package. The page-title scale is fluid from 24–30px; section titles use 17px; data and helper text step down to 10–14px. Weights and restrained negative tracking establish hierarchy. Dates and counts use tabular numerals where available.

## Layout

The desktop workspace uses a fixed 248px navigation rail, a 70px utility header, and a content column capped at 1240px. Summary measures sit above a wide episode table; the follow-through list follows directly below. At tablet widths the rail narrows and metrics form two columns. At 760px and below, navigation becomes a horizontal row and the episode table keeps its readable minimum width inside a horizontal scroll region. Keep headings, status, due date, and owner visible in the first pass through the content.

## Elevation & Depth

The system is mostly flat. Canvas and white surfaces, thin neutral borders, and spacing define groups. Shadows are reserved for transient overlays; do not add a shadow to every card. The focus ring uses `{colors.focus-ring}` with a 3px offset outline.

## Shapes

Use 8px corners for controls and navigation, 12px for data surfaces, and 16px only for larger surfaces. Status tags and small count badges may be pill-shaped. Patient identity and status also have text labels; do not rely on avatar color or color-only state.

## Components

- **Navigation item:** 42px target height with a drawn, consistent 1.5px line icon. Active state uses pale teal and deep teal text.
- **Primary action:** 40px high, teal fill, white text, 8px corners. Disabled actions remain visibly disabled and explain their availability through adjacent context or a title.
- **Summary measure:** a neutral, bordered surface with label, prominent count, and brief operational explanation. Avoid decorative charts that displace episode work.
- **Episode table:** patient reference and display name, ordered test labels, explicit workflow status, due date, and owner. Preserve a horizontal scroll affordance on narrow screens instead of compressing clinical workflow fields into ambiguous abbreviations.
- **Status chip:** text plus a small indicator, with semantic status conveyed in words. Apply warning colors only to genuinely time-sensitive states.
- **Follow-through row:** one next action tied to a patient/episode and accountable owner. Keep the action readable and avoid presenting a decorative checkbox as a working control.

## Do's and Don'ts

### Do:

- Keep synthetic records visibly identified as sample data.
- Make ownership, due date, state, and next action scannable before adding secondary detail.
- Use semantic headings, table headers, keyboard-visible focus, and text labels alongside status color.
- Keep clinical interpretation and next-action authority with clinicians.

### Don't:

- Present the prototype as connected to a health system or as containing real patient records.
- Use color alone for overdue, review, consent, or completion state.
- Add decorative gradients, glass effects, or repeated icon cards that compete with episode content.
- suggest diagnoses, interpret laboratory values, or imply automated clinical decisions.
