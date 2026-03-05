# Frontend Agent Context

## Stack

- React 19 + TypeScript + Vite
- Plain CSS (no Tailwind, no CSS-in-JS)
- Recharts for data visualization
- Path alias: `@/*` → `./src/*`

## Design System

See **[DESIGN_SYSTEM.md](./DESIGN_SYSTEM.md)** for the complete design system reference, including:
- Color palette (CSS variables)
- Typography (Fraunces display + DM Sans body)
- Spacing, border radius, shadows
- Animation tokens and keyframes
- Component patterns (card, button, input)
- Layout conventions and breakpoints
- CSS file architecture

Always use the existing CSS custom properties (`--color-*`, `--space-*`, `--radius-*`, etc.) when adding new styles. Never hardcode color or spacing values.

## Directory Structure

```
src/
├── components/
│   ├── ui/          # Base primitives (Button, Card, Input)
│   ├── energy/      # Energy-specific components
│   └── layout/      # Topbar, nav, logo
├── pages/           # Route-level components
├── styles/          # Feature-scoped CSS files
├── App.css          # Global tokens + base styles
└── index.css        # Minimal reset
```

## Key Conventions

- Feature CSS goes in `src/styles/<feature>.css`, imported where needed
- New components follow existing class naming (BEM-inspired: `.feature-element-modifier`)
- Media queries at end of each CSS file
- Responsive breakpoints: 480px (mobile), 768px (tablet)
