---
name: Renal Relay — inherited starter UI
description: Factual capture of the current Next.js starter; no product visual direction is approved.
colors:
  background: "#ffffff"
  foreground: "#171717"
  background-dark: "#0a0a0a"
  foreground-dark: "#ededed"
typography:
  app-sans:
    fontFamily: "Geist, Arial, Helvetica, sans-serif"
  app-mono:
    fontFamily: "Geist Mono, monospace"
rounded:
  full: "9999px"
  code: "0.25rem"
spacing:
  action-height: "3rem"
  action-inline: "1.25rem"
  page-inline: "4rem"
  page-block: "8rem"
components:
  starter-link-primary:
    backgroundColor: "{colors.foreground}"
    textColor: "{colors.background}"
    rounded: "{rounded.full}"
    height: "{spacing.action-height}"
    padding: "0 {spacing.action-inline}"
  starter-link-secondary:
    backgroundColor: "{colors.background}"
    textColor: "{colors.foreground}"
    rounded: "{rounded.full}"
    height: "{spacing.action-height}"
    padding: "0 {spacing.action-inline}"
---

# Design System: Renal Relay — inherited starter UI

## Overview

This document captures the current `create-next-app` page as it exists in the repository. It is a reference for the incumbent implementation only. The page still identifies itself as the Next.js starter and does not express Renal Relay product workflows or an approved Renal Relay visual identity.

The user has asked that the starter UI be documented while the product's visual direction remains open. No named creative north star, brand palette, component system, or target visual character is established here. Use the product requirements in [`PRODUCT.md`](PRODUCT.md) for product truth; do not treat this starter capture as a target for the future care-team application.

**Key Characteristics:**

- Single centered starter page with two outbound calls to action.
- Light/dark appearance follows the operating-system color preference.
- Tailwind utilities provide layout and framework-default colors; no shared component library exists.

## Colors

The only custom color tokens are the global light/dark foreground and background variables. The page also uses Tailwind default color utilities and a one-off hover color; those are not a Renal Relay semantic palette.

### Neutral

- **Light background** (`{colors.background}`): Global page background and light surface.
- **Light foreground** (`{colors.foreground}`): Default text and primary starter action.
- **Dark background** (`{colors.background-dark}`): Background when the OS requests dark appearance.
- **Dark foreground** (`{colors.foreground-dark}`): Text in dark appearance.

## Typography

**App sans:** Geist is loaded by the root layout and exposed to Tailwind. The page's main container requests the sans utility. The global `body` rule separately specifies Arial/Helvetica, so the repository currently has two overlapping font declarations.

**Mono:** Geist Mono is loaded for code-oriented use.

No reusable type scale has been defined. The starter page uses Tailwind text-size and weight utilities directly in the page markup.

## Layout

The page is a single centered column. Its main region is constrained to Tailwind's `max-w-3xl`, uses generous vertical and horizontal padding, and changes text alignment at the `sm` utility breakpoint. The outer layout fills the viewport height and centers its content. No shared application navigation, dashboard grid, or responsive table pattern exists yet.

## Elevation & Depth

The starter page uses flat surfaces. No box-shadow token or elevation system is defined. Separation comes from background changes and a thin border on the secondary link action.

## Shapes

The two main link actions use fully rounded pill shapes. Inline code uses a small rounded background. The rest of the page has no reusable surface or card shape rules.

## Components

### Starter link actions

These are page-local links styled as actions; they are not reusable Renal Relay controls.

- **Primary link:** Dark foreground surface, light text, pill shape, 48px height, and 20px horizontal padding. Hover changes the background. Dark appearance uses a light action surface.
- **Secondary link:** Light surface, subtle border, pill shape, same height and horizontal padding. Hover removes the border and adds a faint background; dark appearance has corresponding dark styles.
- Both links use a short color transition and become full-width on narrow layouts, with a fixed width at the `md` utility breakpoint.

## Do's and Don'ts

### Do:

- **Do** describe these values as inherited starter defaults when maintaining the current page.
- **Do** use [`PRODUCT.md`](PRODUCT.md) and the relevant product brief before designing a Renal Relay workflow screen.

### Don't:

- **Don't** present the starter colors or pill actions as approved Renal Relay branding.
- **Don't** infer a dashboard layout, clinical status palette, or patient-facing visual language from this single starter page.
- **Don't** treat Tailwind's default palette utilities as custom design tokens.
