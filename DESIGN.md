---
name: Hospitality Excellence
colors:
  surface: '#0b1326'
  surface-dim: '#0b1326'
  surface-bright: '#31394d'
  surface-container-lowest: '#060e20'
  surface-container-low: '#131b2e'
  surface-container: '#171f33'
  surface-container-high: '#222a3d'
  surface-container-highest: '#2d3449'
  on-surface: '#dae2fd'
  on-surface-variant: '#c7c4d7'
  inverse-surface: '#dae2fd'
  inverse-on-surface: '#283044'
  outline: '#908fa0'
  outline-variant: '#464554'
  surface-tint: '#c0c1ff'
  primary: '#c0c1ff'
  on-primary: '#1000a9'
  primary-container: '#8083ff'
  on-primary-container: '#0d0096'
  inverse-primary: '#494bd6'
  secondary: '#b9c8de'
  on-secondary: '#233143'
  secondary-container: '#39485a'
  on-secondary-container: '#a7b6cc'
  tertiary: '#bcc7de'
  on-tertiary: '#263143'
  tertiary-container: '#8691a7'
  on-tertiary-container: '#1f2a3c'
  error: '#ffb4ab'
  on-error: '#690005'
  error-container: '#93000a'
  on-error-container: '#ffdad6'
  primary-fixed: '#e1e0ff'
  primary-fixed-dim: '#c0c1ff'
  on-primary-fixed: '#07006c'
  on-primary-fixed-variant: '#2f2ebe'
  secondary-fixed: '#d4e4fa'
  secondary-fixed-dim: '#b9c8de'
  on-secondary-fixed: '#0d1c2d'
  on-secondary-fixed-variant: '#39485a'
  tertiary-fixed: '#d8e3fb'
  tertiary-fixed-dim: '#bcc7de'
  on-tertiary-fixed: '#111c2d'
  on-tertiary-fixed-variant: '#3c475a'
  background: '#0b1326'
  on-background: '#dae2fd'
  surface-variant: '#2d3449'
typography:
  headline-xl:
    fontFamily: Manrope
    fontSize: 40px
    fontWeight: '700'
    lineHeight: 48px
    letterSpacing: -0.02em
  headline-lg:
    fontFamily: Manrope
    fontSize: 32px
    fontWeight: '600'
    lineHeight: 40px
    letterSpacing: -0.01em
  headline-lg-mobile:
    fontFamily: Manrope
    fontSize: 24px
    fontWeight: '600'
    lineHeight: 32px
  headline-md:
    fontFamily: Manrope
    fontSize: 24px
    fontWeight: '600'
    lineHeight: 32px
  body-lg:
    fontFamily: Inter
    fontSize: 18px
    fontWeight: '400'
    lineHeight: 28px
  body-md:
    fontFamily: Inter
    fontSize: 16px
    fontWeight: '400'
    lineHeight: 24px
  body-sm:
    fontFamily: Inter
    fontSize: 14px
    fontWeight: '400'
    lineHeight: 20px
  label-md:
    fontFamily: Inter
    fontSize: 12px
    fontWeight: '600'
    lineHeight: 16px
    letterSpacing: 0.05em
rounded:
  sm: 0.25rem
  DEFAULT: 0.5rem
  md: 0.75rem
  lg: 1rem
  xl: 1.5rem
  full: 9999px
spacing:
  unit: 8px
  container-max: 1280px
  gutter: 24px
  margin-desktop: 64px
  margin-mobile: 20px
---

## Brand & Style

This design system is engineered for high-end service environments, evoking a sense of quiet luxury, precision, and unwavering reliability. The target audience includes discerning travelers, estate managers, and luxury service providers who require clarity and elegance in high-pressure, low-light environments.

The visual style is **Minimalist Glassmorphism**. It prioritizes deep, atmospheric depth through the use of translucent layers and subtle background blurs, creating a digital experience that feels as refined as a concierge desk at midnight. The interface avoids unnecessary decoration, allowing high-quality imagery and typography to lead the user experience.

- **Emotional Response:** Calm, empowered, sophisticated.
- **Visual Strategy:** Utilize heavy whitespace (even in dark mode) to prevent visual clutter and rely on the interplay of charcoal surfaces to define hierarchy.

## Colors

The palette is anchored in a nocturnal spectrum designed to reduce eye strain while maintaining a premium aesthetic.

- **Primary (Vibrant Indigo):** Reserved strictly for primary calls to action, active states, and critical status indicators. It provides a sharp, luminous contrast against the dark background.
- **Secondary (Slate):** Used for supporting text, icons, and secondary information to maintain a clear visual hierarchy.
- **Surfaces:** The foundation is a deep **Charcoal (#0F172A)**. UI containers and cards utilize **Dark Slate (#1E293B)** to create a layered effect without the need for harsh borders.
- **Accents:** Semantic colors (Success, Warning, Error) should be desaturated slightly to prevent "vibrating" against the dark surfaces.

## Typography

Typography in this design system is balanced and highly legible. **Manrope** provides a modern, geometric warmth for headlines, while **Inter** ensures maximum functional clarity for dense information like booking details and schedules.

- **Contrast:** Ensure all body text maintains at least a 4.5:1 contrast ratio against charcoal backgrounds. 
- **Scale:** Use the `headline-xl` sparingly for high-impact hero sections. 
- **Tracking:** Labels use increased letter spacing to enhance readability at small sizes on mobile devices.

## Layout & Spacing

The layout philosophy follows a **Fluid Grid** model with a strict 8px rhythmic scale. This ensures consistency across all components and page structures.

- **Grid:** A 12-column grid for desktop, 8-column for tablet, and 4-column for mobile.
- **Margins:** Generous outer margins are used on desktop to center the content and provide a "gallery" feel.
- **Density:** In management views, spacing can be tightened to 4px increments, but for guest-facing interfaces, use 16px and 24px units to maintain an airy, luxury feel.

## Elevation & Depth

Visual hierarchy is established through **Tonal Layers** and **Glassmorphism**. 

1. **Base (Level 0):** Deep Charcoal (#0F172A). The "ground" of the application.
2. **Surface (Level 1):** Dark Slate (#1E293B). Used for cards and secondary navigation elements.
3. **Overlay (Level 2):** Semi-transparent Slate with a 12px Backdrop Blur. Used for floating headers, modals, and dropdown menus.

**Shadows:** Shadows are rarely used. When necessary, use a large, soft blur (32px+) with a low-opacity Indigo tint (#6366F1 at 15%) to suggest a subtle "glow" rather than a physical shadow cast by a light source.

## Shapes

The shape language is refined and approachable. A consistent **0.5rem (8px)** corner radius is applied to standard components like buttons and input fields.

- **Large Components:** Cards and large containers should use `rounded-lg` (1rem) to soften the overall interface.
- **Interactive Elements:** Small interactive elements like checkboxes use `rounded-sm` (4px) to maintain a sense of precision.

## Components

Components are optimized for visibility and interaction in dark mode.

- **Buttons:** 
    - *Primary:* Solid Indigo with white text. High-contrast.
    - *Secondary:* Ghost style with a Slate border (1px) and subtle hover fill.
- **Input Fields:** Use a Dark Slate fill with a 1px border. The border transitions to Indigo on focus. Error states use a soft coral-red border, never pure red, to avoid harshness.
- **Cards:** No borders by default. Depth is communicated via the Dark Slate surface color against the Charcoal background. On hover, apply a subtle 1px Indigo stroke.
- **Chips/Badges:** Low-opacity Indigo background with high-opacity Indigo text for an "illuminated" effect.
- **Lists:** Separated by thin, low-contrast Slate lines (opacity 10%) to guide the eye without breaking the flow of the page.
- **Specialty Components:** Include a "Night Mode Optimizer"—a toggle that further desaturates images to ensure the UI remains comfortable in pitch-black environments.