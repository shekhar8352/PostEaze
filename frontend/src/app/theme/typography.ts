/**
 * Typography Configuration
 * 
 * Centralized typography settings for the PostEaze application.
 * Update these values to change fonts, sizes, and weights across the entire app.
 */

// Font Families
export const FONT_FAMILIES = {
  body: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif",
  heading: "'Plus Jakarta Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
  mono: "'Fira Code', 'Courier New', monospace",
} as const;

// Font Weights
export const FONT_WEIGHTS = {
  light: '300',
  regular: '400',
  medium: '500',
  semibold: '600',
  bold: '700',
  extrabold: '800',
} as const;

// Font Sizes
export const FONT_SIZES = {
  xs: '12px',
  sm: '14px',
  md: '16px',
  lg: '18px',
  xl: '20px',
} as const;

// Line Heights
export const LINE_HEIGHTS = {
  xs: '1.4',
  sm: '1.45',
  md: '1.55',
  lg: '1.6',
  xl: '1.65',
} as const;

// Heading Configuration
export const HEADING_SIZES = {
  h1: { fontSize: '36px', lineHeight: '1.2', fontWeight: FONT_WEIGHTS.extrabold },
  h2: { fontSize: '30px', lineHeight: '1.3', fontWeight: FONT_WEIGHTS.bold },
  h3: { fontSize: '24px', lineHeight: '1.4', fontWeight: FONT_WEIGHTS.semibold },
  h4: { fontSize: '20px', lineHeight: '1.5', fontWeight: FONT_WEIGHTS.semibold },
  h5: { fontSize: '18px', lineHeight: '1.5', fontWeight: FONT_WEIGHTS.semibold },
  h6: { fontSize: '16px', lineHeight: '1.5', fontWeight: FONT_WEIGHTS.semibold },
} as const;

// Letter Spacing
export const LETTER_SPACING = {
  tight: '-0.02em',
  normal: '0',
  wide: '0.02em',
} as const;
