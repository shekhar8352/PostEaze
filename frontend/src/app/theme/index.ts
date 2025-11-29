/**
 * Theme Configuration
 * 
 * Main theme export that composes all theme modules.
 * Import this in MantineProvider to apply the theme.
 */

import { createTheme } from '@mantine/core';
import { FONT_FAMILIES, FONT_SIZES, LINE_HEIGHTS, HEADING_SIZES } from './typography';
import { BRAND_COLORS } from './colors';
import { COMPONENT_OVERRIDES } from './components';

// Export all theme modules for direct access
export * from './typography';
export * from './colors';
export * from './spacing';
export * from './components';
export * from './icons';

// Composed Mantine Theme
export const theme = createTheme({
  primaryColor: BRAND_COLORS.primary,
  
  // Typography
  fontFamily: FONT_FAMILIES.body,
  fontSizes: FONT_SIZES,
  lineHeights: LINE_HEIGHTS,
  
  // Headings
  headings: {
    fontFamily: FONT_FAMILIES.heading,
    fontWeight: '700',
    sizes: HEADING_SIZES,
  },
  
  // Component Overrides
  components: COMPONENT_OVERRIDES,
});
