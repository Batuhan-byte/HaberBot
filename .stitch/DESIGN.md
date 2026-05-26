---
name: Terminal Monolith
colors:
  surface: '#131313'
  surface-dim: '#131313'
  surface-bright: '#3a3939'
  surface-container-lowest: '#0e0e0e'
  surface-container-low: '#1c1b1b'
  surface-container: '#201f1f'
  surface-container-high: '#2a2a2a'
  surface-container-highest: '#353534'
  on-surface: '#e5e2e1'
  on-surface-variant: '#c4c7c8'
  inverse-surface: '#e5e2e1'
  inverse-on-surface: '#313030'
  outline: '#8e9192'
  outline-variant: '#444748'
  surface-tint: '#c6c6c7'
  primary: '#ffffff'
  on-primary: '#2f3131'
  primary-container: '#e2e2e2'
  on-primary-container: '#636565'
  inverse-primary: '#5d5f5f'
  secondary: '#c7c6c6'
  on-secondary: '#303031'
  secondary-container: '#464747'
  on-secondary-container: '#b5b5b5'
  tertiary: '#ffffff'
  on-tertiary: '#2f3131'
  tertiary-container: '#e2e2e2'
  on-tertiary-container: '#636565'
  error: '#ffb4ab'
  on-error: '#690005'
  error-container: '#93000a'
  on-error-container: '#ffdad6'
  primary-fixed: '#e2e2e2'
  primary-fixed-dim: '#c6c6c7'
  on-primary-fixed: '#1a1c1c'
  on-primary-fixed-variant: '#454747'
  secondary-fixed: '#e3e2e2'
  secondary-fixed-dim: '#c7c6c6'
  on-secondary-fixed: '#1b1c1c'
  on-secondary-fixed-variant: '#464747'
  tertiary-fixed: '#e2e2e2'
  tertiary-fixed-dim: '#c6c6c7'
  on-tertiary-fixed: '#1a1c1c'
  on-tertiary-fixed-variant: '#454747'
  background: '#131313'
  on-background: '#e5e2e1'
  surface-variant: '#353534'
typography:
  display-lg:
    fontFamily: Geist
    fontSize: 48px
    fontWeight: '600'
    lineHeight: '1.1'
    letterSpacing: -0.04em
  headline-lg:
    fontFamily: Geist
    fontSize: 32px
    fontWeight: '600'
    lineHeight: '1.2'
    letterSpacing: -0.02em
  headline-md:
    fontFamily: Geist
    fontSize: 24px
    fontWeight: '500'
    lineHeight: '1.3'
    letterSpacing: -0.01em
  body-lg:
    fontFamily: Geist
    fontSize: 16px
    fontWeight: '400'
    lineHeight: '1.6'
    letterSpacing: '0'
  body-md:
    fontFamily: Geist
    fontSize: 14px
    fontWeight: '400'
    lineHeight: '1.5'
    letterSpacing: '0'
  label-md:
    fontFamily: JetBrains Mono
    fontSize: 13px
    fontWeight: '500'
    lineHeight: '1'
    letterSpacing: 0.02em
  label-sm:
    fontFamily: JetBrains Mono
    fontSize: 11px
    fontWeight: '400'
    lineHeight: '1'
    letterSpacing: 0.05em
  headline-lg-mobile:
    fontFamily: Geist
    fontSize: 28px
    fontWeight: '600'
    lineHeight: '1.2'
rounded:
  sm: 0.125rem
  DEFAULT: 0.25rem
  md: 0.375rem
  lg: 0.5rem
  xl: 0.75rem
  full: 9999px
spacing:
  unit: 4px
  gutter: 16px
  margin-mobile: 16px
  margin-desktop: 32px
  max-width: 1280px
---

## Brand & Style

The design system is engineered for high-performance developer tools, prioritizing technical precision and cognitive clarity over decorative elements. It targets a sophisticated audience that values efficiency, speed, and information density.

The aesthetic is **Ultra-Minimalist / High-Contrast**, drawing heavily from the "Linear" and "Vercel" school of design. It utilizes a strict monochromatic palette to eliminate visual noise, relying on razor-sharp borders and intentional whitespace to define structure. The emotional response should be one of professional reliability, high-tech sophistication, and absolute focus.

**Key Principles:**
- **Extreme Reduction:** If a visual element doesn't serve a functional purpose, it is removed.
- **Micro-interactions:** Motion is fast (150ms) and functional, used only to confirm actions or state changes.
- **Developer-Centric:** Layouts favor logical grouping and data readability.

## Colors

This design system operates on a **Pure Dark Mode** architecture. The palette is strictly limited to grayscale to ensure that content—specifically code and data—remains the focal point.

- **Background:** Pure black (#000000) is used for the primary canvas to maximize OLED contrast and reduce eye strain.
- **Surface:** Tiered grays (#0A0A0A for primary surfaces, #111111 for cards/modals) create subtle depth without breaking the minimalist aesthetic.
- **Accent:** White (#FFFFFF) is reserved for high-priority text and primary actions. 
- **Borders:** A consistent #222222 border is the primary tool for element separation, replacing shadows.
- **Status:** Functional colors (Success: #00FF00, Error: #FF0000) should be used sparingly as 1px indicators or small icons only.

## Typography

The typography system relies on **Geist** for its geometric, technical feel and high legibility. For technical metadata and code snippets, **JetBrains Mono** is utilized to provide a "developer-first" context.

- **Contrast:** High contrast between primary headings (White) and secondary body text (Gray #888888).
- **Scale:** Large display type is used sparingly for hero sections. Functional UI uses a compact 14px base for high information density.
- **Mono Integration:** Labels, tags, and IDs must always use the monospaced font to differentiate data from prose.

## Layout & Spacing

The layout follows a **Strict Grid System** based on 4px increments. It utilizes a fixed-fluid hybrid model where content resides in a centered container on large screens.

- **Grid:** 12-column desktop grid with 16px gutters.
- **Margins:** 32px on desktop, 16px on mobile.
- **Padding:** Consistent internal padding of 24px for cards and containers to maintain a spacious, premium feel despite the dark theme.
- **Density:** Information density is high, but balanced by generous "macro-whitespace" between major sections.

## Elevation & Depth

This design system rejects traditional shadows in favor of **Tonal Layering** and **Subtle Outlines**.

- **Level 0 (Base):** #000000 (Pure Black).
- **Level 1 (Cards/Sidebar):** #0A0A0A with a 1px solid #222222 border.
- **Level 2 (Modals/Popovers):** #111111 with a 1px solid #333333 border.
- **Interactions:** Hover states are signaled by increasing the border brightness (to #444444) or a subtle background shift (to #161616). No blurs or glassmorphism are permitted; surfaces remain opaque and grounded.

## Shapes

The shape language is **Strict and Geometric**. Small border radii are used to provide a modern "software" feel while maintaining a technical edge.

- **Standard Elements:** 4px (0.25rem) radius for buttons, inputs, and small chips.
- **Containers:** 8px (0.5rem) radius for cards and larger modules.
- **Strictness:** Rounding should never exceed 8px. Pill shapes are only permitted for status badges and toggles.

## Components

### Buttons
- **Primary:** Solid White background, Black text. No border.
- **Secondary:** Transparent background, #222222 border, White text. Hover: Border becomes #444444.
- **Ghost:** Transparent background, Gray text. Hover: White text, subtle #111111 background.

### Input Fields
- Background: #000000. Border: #222222.
- Focus: Border becomes #FFFFFF.
- Type: Monospaced font for technical inputs (APIs, IDs).

### Cards
- Background: #0A0A0A. Border: 1px #222222.
- Header: Separated by a 1px horizontal line (#222222).

### Data Tables
- Header: #888888 text, uppercase, label-sm font style.
- Rows: Separated by 1px horizontal borders (#111111).
- Hover: Row background shifts to #0A0A0A.

### Toggle Switches
- Track: #222222 when off, #FFFFFF when on.
- Thumb: Circle, #000000 (Black) when on to create a high-contrast cutout effect.

### Chips / Badges
- Small, 4px rounded, monospaced text. Background #111111, border #222222.
