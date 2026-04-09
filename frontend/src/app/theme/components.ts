/**
 * Component Style Overrides
 * 
 * Mantine component default props and style overrides.
 * Customize component behavior and appearance here.
 */

import type { MantineTheme, MantineThemeComponents } from '@mantine/core';

export const COMPONENT_OVERRIDES: MantineThemeComponents = {
  Text: {
    defaultProps: {
      size: 'md',
    },
  },
  Title: {
    defaultProps: {
      fw: 600,
    },
  },
  Button: {
    defaultProps: {
      radius: 'md',
    },
  },
  Input: {
    defaultProps: {
      radius: 'md',
    },
  },
  Card: {
    defaultProps: {
      radius: 'md',
      shadow: 'sm',
      padding: 'lg',
    },
  },
  Paper: {
    defaultProps: {
      radius: 'md',
      shadow: 'xs',
    },
  },
  Modal: {
    defaultProps: {
      radius: 'md',
      centered: true,
    },
  },
  NavLink: {
    styles: (_theme: MantineTheme, props: { active?: boolean }) =>
      props.active
        ? {}
        : {
            label: { color: 'var(--pe-text-secondary)' },
            section: { color: 'var(--pe-text-muted)' },
          },
  },
  Notification: {
    defaultProps: {
      radius: 'md',
    },
  },
};
