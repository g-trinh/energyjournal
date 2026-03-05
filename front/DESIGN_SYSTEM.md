# Design System

## Overview

Warm, organic aesthetic with deep earthy tones. Key techniques: glassmorphism (`backdrop-filter: blur(20px)`), ambient glow effects, and smooth easing animations. No Tailwind — plain CSS with custom properties.

**Main files:**
- `src/App.css` — all design tokens, base styles, shared animations
- `src/index.css` — minimal reset
- `src/styles/` — feature-scoped stylesheets

---

## Color Palette

All defined as CSS variables in `src/App.css`.

### Backgrounds
| Token | Value | Usage |
|---|---|---|
| `--color-bg-deep` | `#1a1816` | Page background |
| `--color-bg-elevated` | `#242220` | Elevated surfaces |
| `--color-bg-card` | `rgba(38, 35, 32, 0.7)` | Card background |
| `--color-bg-card-hover` | `rgba(48, 44, 40, 0.8)` | Card hover |

### Text
| Token | Value | Usage |
|---|---|---|
| `--color-text-primary` | `#f5f0e8` | Main content |
| `--color-text-secondary` | `#a09a90` | Supporting text |
| `--color-text-muted` | `#6b665c` | Disabled / hint |

### Accents
| Token | Value | Usage |
|---|---|---|
| `--color-accent-amber` | `#e8a445` | Primary CTA, focus, active |
| `--color-accent-amber-glow` | `rgba(232, 164, 69, 0.3)` | Glow / shadow |
| `--color-accent-sage` | `#8fa58b` | Secondary actions, success |
| `--color-accent-terracotta` | `#c4826d` | Errors, destructive |
| `--color-accent-seafoam` | `#7eb8b3` | Tertiary info |

### Borders
| Token | Value |
|---|---|
| `--color-border` | `rgba(160, 154, 144, 0.12)` |
| `--color-border-hover` | `rgba(232, 164, 69, 0.3)` |

---

## Typography

**Fonts:** Google Fonts — `Fraunces` (display/serif) + `DM Sans` (body/sans-serif)

| Token | Value |
|---|---|
| `--font-display` | `'Fraunces', serif` |
| `--font-body` | `'DM Sans', sans-serif` |

### Scale
| Role | Size | Weight | Font |
|---|---|---|---|
| Page title | `clamp(2rem, 5vw, 3.25rem)` | 500 | Fraunces |
| Card title | `2rem` | 500 | Fraunces |
| Section header | `1.5rem` | 500–600 | Fraunces |
| Body | `0.95–1rem` | 400–500 | DM Sans |
| Label / tag | `0.72–0.85rem` | 500–600, uppercase, `letter-spacing: 0.08–0.12em` | DM Sans |

---

## Spacing

8px base unit.

| Token | Value |
|---|---|
| `--space-xs` | `0.25rem` (4px) |
| `--space-sm` | `0.5rem` (8px) |
| `--space-md` | `1rem` (16px) |
| `--space-lg` | `1.5rem` (24px) |
| `--space-xl` | `2rem` (32px) |
| `--space-2xl` | `3rem` (48px) |
| `--space-3xl` | `4rem` (64px) |

**Common patterns:** card padding `--space-xl`, section gap `--space-xl` to `--space-2xl`, field gap `--space-lg`.

---

## Border Radius

| Token | Value | Usage |
|---|---|---|
| `--radius-sm` | `8px` | Buttons, inputs |
| `--radius-md` | `12px` | Cards, dropdowns |
| `--radius-lg` | `20px` | Large cards |
| `--radius-xl` | `28px` | Full-width cards |
| `--radius-infinite` | `999px` | Pills |

---

## Shadows

| Token | Value |
|---|---|
| `--shadow-sm` | `0 2px 8px rgba(0,0,0,0.15)` |
| `--shadow-md` | `0 4px 20px rgba(0,0,0,0.2)` |
| `--shadow-lg` | `0 8px 40px rgba(0,0,0,0.25)` |
| `--shadow-glow` | `0 0 40px var(--color-accent-amber-glow)` |

---

## Animation

### Easing
| Token | Value |
|---|---|
| `--ease-out-expo` | `cubic-bezier(0.16, 1, 0.3, 1)` |
| `--ease-out-back` | `cubic-bezier(0.34, 1.56, 0.64, 1)` |

### Duration
| Token | Value |
|---|---|
| `--duration-fast` | `150ms` |
| `--duration-normal` | `300ms` |
| `--duration-slow` | `500ms` |

### Keyframes (defined in `App.css`)
- `slideDown` — opacity + translateY(-20px → 0)
- `slideUp` — opacity + translateY(20px → 0)
- `fadeIn` — opacity 0 → 1
- `spin` — 360deg, 1s linear
- `pulse` — opacity oscillation, 2s
- `ambientFloat` — scale + translate, 20s (background decorations)

**Reduced motion:** all durations collapse to `0.01ms`, iteration count set to 1.

---

## Component Patterns

### Card
```css
background: var(--color-bg-card);
border: 1px solid var(--color-border);
border-radius: var(--radius-lg);
backdrop-filter: blur(20px);
padding: var(--space-xl);
transition: transform var(--duration-normal), border-color var(--duration-normal);

&:hover {
  transform: translateY(-4px);
  border-color: var(--color-border-hover);
}
```

### Primary Button
```css
background: linear-gradient(135deg, var(--color-accent-amber) 0%, #d4923d 100%);
color: var(--color-bg-deep);
border-radius: var(--radius-sm);
box-shadow: var(--shadow-sm), 0 0 20px var(--color-accent-amber-glow);

&:hover { transform: translateY(-2px); }
&:active { transform: translateY(0); }
&:disabled { opacity: 0.7; cursor: not-allowed; }
```

### Secondary Button
```css
background: var(--color-accent-sage);
/* same hover/active pattern as primary */
```

### Outline / Ghost Button
```css
background: transparent;
border: 1px solid var(--color-border);
color: var(--color-text-secondary);

&:hover { border-color: var(--color-border-hover); }
/* active: amber highlight */
```

### Input
```css
border: 1px solid var(--color-border);
background: var(--color-bg-deep);
border-radius: var(--radius-sm);

&:focus { border-color: var(--color-accent-amber); box-shadow: 0 0 0 3px var(--color-accent-amber-glow); }
&:hover { border-color: var(--color-border-hover); }
/* error: border-color: var(--color-accent-terracotta) */
```

### Focus Ring (universal)
```css
outline: 2px solid var(--color-accent-amber);
outline-offset: 2px;
/* or: box-shadow: 0 0 0 3px var(--color-accent-amber-glow) */
```

---

## Layout

| Context | Max-width |
|---|---|
| Main content | `900px–980px` |
| Auth pages | `1160px` |
| Landing page | `1280px` |

**Topbar height:** `--topbar-height: 64px` (desktop), `--topbar-height-mobile: 56px`

**Responsive breakpoints:**
- Mobile: `max-width: 480px`
- Tablet: `max-width: 767px`
- Desktop: `768px+`

Common responsive changes: 3-col grid → 1-col, padding `--space-xl` → `--space-lg`, flex row → column.

---

## CSS Architecture

```
src/
├── App.css              # Design tokens, global base, keyframes
├── index.css            # Reset only
└── styles/
    ├── topbar.css       # Header & nav
    ├── auth.css         # Login / signup
    ├── energy-levels.css
    ├── energy-edit.css
    ├── energy-context.css
    ├── landing.css
    └── calendar.css
```

**Conventions:** BEM-inspired class names (`.topbar-inner`, `.auth-card-title`), feature-based file split, CSS custom properties for theming, no preprocessor, media queries at end of each file.
