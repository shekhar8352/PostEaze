/**
 * Spacing and Sizing Configuration
 * 
 * Centralized spacing, sizing, and layout constants.
 * Update these values to change spacing and sizes across the entire app.
 */

// Spacing Scale
export const SPACING = {
  xs: '8px',
  sm: '12px',
  md: '16px',
  lg: '24px',
  xl: '32px',
  xxl: '48px',
} as const;

// Layout Sizes
export const LAYOUT_SIZES = {
  headerHeight: 60,
  sidebarWidth: 280,
  sidebarCollapsedWidth: 80,
  maxContentWidth: 1200,
} as const;

// Border Radius
export const BORDER_RADIUS = {
  xs: '4px',
  sm: '8px',
  md: '12px',
  lg: '16px',
  xl: '24px',
  full: '9999px',
} as const;

// Shadows
export const SHADOWS = {
  xs: '0 1px 3px rgba(0, 0, 0, 0.05)',
  sm: '0 1px 3px rgba(0, 0, 0, 0.1), 0 1px 2px rgba(0, 0, 0, 0.06)',
  md: '0 4px 6px rgba(0, 0, 0, 0.1), 0 2px 4px rgba(0, 0, 0, 0.06)',
  lg: '0 10px 15px rgba(0, 0, 0, 0.1), 0 4px 6px rgba(0, 0, 0, 0.05)',
  xl: '0 20px 25px rgba(0, 0, 0, 0.1), 0 10px 10px rgba(0, 0, 0, 0.04)',
} as const;

// Z-Index Scale
export const Z_INDEX = {
  dropdown: 1000,
  sticky: 1020,
  fixed: 1030,
  modalBackdrop: 1040,
  modal: 1050,
  popover: 1060,
  tooltip: 1070,
} as const;

// Breakpoints (matches Mantine defaults)
export const BREAKPOINTS = {
  xs: '36em',  // 576px
  sm: '48em',  // 768px
  md: '62em',  // 992px
  lg: '75em',  // 1200px
  xl: '88em',  // 1408px
} as const;
