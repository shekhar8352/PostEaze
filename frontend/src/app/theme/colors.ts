/**
 * Color Configuration
 * 
 * Centralized color palette for the PostEaze application.
 * Update these values to change colors across the entire app.
 */

// Brand Colors — primary matches `--pe-accent` (blue) so Mantine buttons/links match workspace tokens
export const BRAND_COLORS = {
  primary: 'blue',
  secondary: 'cyan',
  accent: 'blue',
} as const;

// Semantic Colors
export const SEMANTIC_COLORS = {
  success: 'green',
  warning: 'yellow',
  error: 'red',
  info: 'blue',
} as const;

// Channel Colors (for branding)
export const CHANNEL_COLORS = {
  instagram: {
    gradient: 'linear-gradient(45deg, #f09433 0%, #e6683c 25%, #dc2743 50%, #cc2366 75%, #bc1888 100%)',
    solid: '#E1306C',
  },
  facebook: {
    solid: '#1877F2',
  },
  youtube: {
    solid: '#FF0000',
  },
  twitter: {
    solid: '#1DA1F2',
  },
  linkedin: {
    solid: '#0A66C2',
  },
} as const;

// Gradient Definitions
export const GRADIENTS = {
  primary: { from: 'blue', to: 'cyan', deg: 45 },
  success: { from: 'teal', to: 'green', deg: 45 },
  warning: { from: 'yellow', to: 'orange', deg: 45 },
  error: { from: 'red', to: 'pink', deg: 45 },
} as const;

// Neutral Colors (for future dark mode support)
export const NEUTRAL_COLORS = {
  gray: 'gray',
  dark: 'dark',
} as const;
