---
version: alpha
name: Vite App Design System
description: A modern, minimalist, high-contrast design system featuring a striking golden yellow accent, dark grey ink, and friendly rounded corners, built with Tailwind CSS v4 and React.
colors:
  # Light Mode (Default)
  background: "#ffffff"
  foreground: "#0a0a0a"
  card: "#ffffff"
  card-foreground: "#0a0a0a"
  popover: "#ffffff"
  popover-foreground: "#0a0a0a"
  primary: "#fdc700"
  primary-foreground: "#733e0a"
  secondary: "#f4f4f5"
  secondary-foreground: "#18181b"
  muted: "#f5f5f5"
  muted-foreground: "#737373"
  accent: "#f5f5f5"
  accent-foreground: "#171717"
  destructive: "#e7000b"
  border: "#e5e5e5"
  input: "#e5e5e5"
  ring: "#a1a1a1"
  chart-1: "#ffb86a"
  chart-2: "#ff6900"
  chart-3: "#f54900"
  chart-4: "#ca3500"
  chart-5: "#9f2d00"
  sidebar: "#fafafa"
  sidebar-foreground: "#0a0a0a"
  sidebar-primary: "#d08700"
  sidebar-primary-foreground: "#fefce8"
  sidebar-accent: "#f5f5f5"
  sidebar-accent-foreground: "#171717"
  sidebar-border: "#e5e5e5"
  sidebar-ring: "#a1a1a1"

  # Dark Mode Overrides
  dark-background: "#0a0a0a"
  dark-foreground: "#fafafa"
  dark-card: "#171717"
  dark-card-foreground: "#fafafa"
  dark-popover: "#171717"
  dark-popover-foreground: "#fafafa"
  dark-primary: "#f0b100"
  dark-primary-foreground: "#733e0a"
  dark-secondary: "#27272a"
  dark-secondary-foreground: "#fafafa"
  dark-muted: "#262626"
  dark-muted-foreground: "#a1a1a1"
  dark-accent: "#262626"
  dark-accent-foreground: "#fafafa"
  dark-destructive: "#ff6467"
  dark-border: "#ffffff1a"
  dark-input: "#ffffff26"
  dark-ring: "#737373"
  dark-sidebar: "#171717"
  dark-sidebar-foreground: "#fafafa"
  dark-sidebar-primary: "#f0b100"
  dark-sidebar-primary-foreground: "#fefce8"
  dark-sidebar-accent: "#262626"
  dark-sidebar-accent-foreground: "#fafafa"
  dark-sidebar-ring: "#737373"

typography:
  h1:
    fontFamily: Nunito Sans
    fontSize: 24px
    fontWeight: 700
    lineHeight: 1.2
  body-md:
    fontFamily: Nunito Sans
    fontSize: 14px
    fontWeight: 400
    lineHeight: 2.0
  mono-xs:
    fontFamily: monospace
    fontSize: 12px
    fontWeight: 400
    lineHeight: 1.5

spacing:
  xs: 4px
  sm: 8px
  md: 16px
  lg: 24px
  xl: 32px
  xxl: 48px
  xxxl: 64px

rounded:
  sm: 6px
  md: 8px
  lg: 10px
  xl: 14px
  xxl: 18px
  xxxl: 22px
  xxxxl: 26px
  full: 9999px

components:
  button-primary:
    backgroundColor: "{colors.primary}"
    textColor: "{colors.primary-foreground}"
    rounded: "{rounded.xxl}"
    height: "32px"
    padding: "12px"
  button-primary-hover:
    backgroundColor: "#fdc700cc"
  button-secondary:
    backgroundColor: "{colors.secondary}"
    textColor: "{colors.secondary-foreground}"
    rounded: "{rounded.xxl}"
    height: "32px"
    padding: "12px"
  button-outline:
    backgroundColor: "{colors.background}"
    textColor: "{colors.foreground}"
    rounded: "{rounded.xxl}"
    height: "32px"
    padding: "12px"
  button-ghost:
    textColor: "{colors.foreground}"
    rounded: "{rounded.xxl}"
    height: "32px"
    padding: "12px"
  button-destructive:
    backgroundColor: "#e7000b1a"
    textColor: "{colors.destructive}"
    rounded: "{rounded.xxl}"
    height: "32px"
    padding: "12px"
---

# Vite App Design System

This document specifies the design system and visual identity rules for the Vite App project. It acts as a contract between designers, developers, and AI coding agents to ensure visual consistency across light and dark modes.

## Overview

The visual identity is defined by a bright, clean, minimalist design style. The system uses **Nunito Sans** as its primary typeface, offering a friendly yet professional voice. A high-contrast, energetic golden yellow (`#fdc700`) accent anchors critical interactions, combined with dark gray ink and soft, generous rounded corners to create an approachable, modern interface.

## Colors

The color palette is built around dynamic light and dark modes, prioritizing high readability and intentional visual focus:

- **Primary (#fdc700 / #f0b100):** Golden Yellow is the main brand accent, used for high-importance visual elements and primary action buttons.
- **Primary Foreground (#733e0a):** Rich dark brown ensures excellent WCAG AA contrast over the yellow background.
- **Background (#ffffff / #0a0a0a):** Pure white in light mode and deep charcoal in dark mode provide the workspace foundation.
- **Secondary (#f4f4f5 / #27272a):** Soft cool-gray tones used for auxiliary content areas and subtle groupings.
- **Muted/Accent (#f5f5f5 / #262626):** Light backgrounds for tooltips, card states, and hover feedback.
- **Destructive (#e7000b / #ff6467):** A sharp, bright red to signal errors, warnings, and dangerous actions.
- **Borders & Inputs:** Neutral frame boundaries to group relative layout items without adding visual noise.

## Typography

Typography is handled using **Nunito Sans** for all text (headings and body).

- **Headlines:** Set in `Nunito Sans` bold/semibold for clear readability and a welcoming character.
- **Body text:** Set in `Nunito Sans` regular at `14px` (`text-sm`) with a generous `leading-loose` line height (2x font-size) to maximize long-form readability.
- **Monospace:** Standard system monospace (`monospace`) is used for technical indicators, keyboard keycaps, code, and configuration output.

## Layout

The layout is built with fluid grid alignments, utilizing a standard 8px-based spacing scale:

- A strict spacing rhythm is maintained using `xs` (4px), `sm` (8px), `md` (16px), `lg` (24px), and `xl` (32px).
- Primary page padding starts at `24px` (`p-6`) for workspace margins.
- Flex layouts with vertical gap hierarchies are used to establish logical information groupings.

## Elevation & Depth

A flat and tonal design style is preferred over heavy box shadows. 

- Depth is achieved via background layer separation (e.g., pure white content cards sitting on off-white surfaces).
- Subtle border framing is used to group relative elements (`#e5e5e5` in light mode, soft semi-transparent white in dark mode).
- Focus states use a ring glow (`focus-visible:ring-3`) to communicate active keyboard focus states.

## Shapes

The shape language is defined by **Friendly Rounded Corners**:

- Standard cards, containers, and inputs use a base corner radius of `10px` (`rounded-lg` or `0.625rem`) to create soft, approachable bounds.
- Small sub-components use `8px` (`rounded-md` or `0.5rem`) and `6px` (`rounded-sm` or `0.375rem`).
- Interactive buttons utilize a softer, more rounded corner footprint of `18px` (`rounded-2xl` or `1.125rem`).

## Components

The system specifies the following rules for component atoms:

- **Buttons:** Built on an `18px` rounded footprint, utilizing standard heights (e.g., `32px` default height) and horizontal padding. 
- **Button Hover States:** Hover states utilize opacity reductions (e.g., primary shifts to 80% opacity) or subtle color-mixes to preserve design consistency while giving clear interactive feedback.
- **Destructive Button:** Uses a subtle red background overlay (`10%` to `20%` opacity) with sharp red text to ensure visibility without overwhelming the layout.

## Do's and Don'ts

- **Do** use the golden yellow accent for primary visual highlights and call-to-actions.
- **Do** ensure primary buttons use the high-contrast brown text (`#733e0a`) for legibility.
- **Don't** mix the `18px` button radius with sharp corners on UI panels.
- **Do** support light and dark theme switching cleanly via the `dark` class selector.
- **Don't** use arbitrary padding values; stick strictly to the 4px/8px-based spacing scale.