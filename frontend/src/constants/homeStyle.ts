import type { HomeStyle } from '@/types'

/** Shipped defaults, matching the built-in homepage appearance. */
export const DEFAULT_HOME_STYLE: HomeStyle = {
  accent_from: '#6366F1',
  accent_to: '#06B6D4',
  hero_title: '',
  hero_desc: '',
  show_providers: true,
  show_pain_points: true,
  show_comparison: true,
  show_terminal: true
}

const HEX_COLOR_PATTERN = /^#[0-9a-fA-F]{6}$/

export function isHexColor(value: string): boolean {
  return HEX_COLOR_PATTERN.test(value.trim())
}

function normalizeHex(value: unknown, fallback: string): string {
  return typeof value === 'string' && isHexColor(value) ? value.trim().toUpperCase() : fallback
}

/**
 * Fills a partial or missing payload with the shipped defaults, so an older
 * server response can never blank the homepage or hide every section.
 */
export function normalizeHomeStyle(value?: Partial<HomeStyle> | null): HomeStyle {
  if (!value || typeof value !== 'object') {
    return { ...DEFAULT_HOME_STYLE }
  }

  return {
    accent_from: normalizeHex(value.accent_from, DEFAULT_HOME_STYLE.accent_from),
    accent_to: normalizeHex(value.accent_to, DEFAULT_HOME_STYLE.accent_to),
    hero_title: typeof value.hero_title === 'string' ? value.hero_title : '',
    hero_desc: typeof value.hero_desc === 'string' ? value.hero_desc : '',
    show_providers: value.show_providers !== false,
    show_pain_points: value.show_pain_points !== false,
    show_comparison: value.show_comparison !== false,
    show_terminal: value.show_terminal !== false
  }
}

function hexToRgbTriplet(hex: string): string {
  const normalized = hex.replace('#', '')
  const r = Number.parseInt(normalized.slice(0, 2), 16)
  const g = Number.parseInt(normalized.slice(2, 4), 16)
  const b = Number.parseInt(normalized.slice(4, 6), 16)
  return `${r} ${g} ${b}`
}

/**
 * Custom properties consumed by HomeView's scoped `.brand-*` classes. Driving
 * the brand color through variables keeps every gradient and tint in sync
 * without needing a Tailwind rebuild for each admin-chosen color.
 */
export function homeStyleCssVars(style: HomeStyle): Record<string, string> {
  return {
    '--brand-from': style.accent_from,
    '--brand-to': style.accent_to,
    '--brand-from-rgb': hexToRgbTriplet(style.accent_from),
    '--brand-to-rgb': hexToRgbTriplet(style.accent_to)
  }
}
