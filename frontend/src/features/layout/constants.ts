
import { type NavItem } from './types';
import { LAYOUT_SIZES, Icons } from '@/app/theme';

// Sidebar dimensions (imported from theme)
export const SIDEBAR_WIDTH = LAYOUT_SIZES.sidebarWidth;
export const SIDEBAR_COLLAPSED_WIDTH = LAYOUT_SIZES.sidebarCollapsedWidth;

// Main navigation items
export const MAIN_NAV_ITEMS: NavItem[] = [
  {
    label: 'Dashboard',
    icon: Icons.Home,
    path: '/dashboard',
  },
  {
    label: 'Analytics',
    icon: Icons.ChartBar,
    path: '/analytics',
  },
  {
    label: 'Studio',
    icon: Icons.Kanban,
    path: '/studio',
  },
  {
    label: 'Calendar',
    icon: Icons.Calendar,
    path: '/calendar',
  },
  {
    label: 'Workspace',
    icon: Icons.Photo,
    path: '/workspace',
  },
];

// Channel navigation items
export const CHANNEL_NAV_ITEMS: NavItem[] = [
  {
    label: 'Instagram',
    icon: Icons.Instagram,
    path: '/channels/instagram',
  },
  {
    label: 'Facebook',
    icon: Icons.Facebook,
    path: '/channels/facebook',
  },
  {
    label: 'YouTube',
    icon: Icons.YouTube,
    path: '/channels/youtube',
  },
];

// Bottom navigation items
export const BOTTOM_NAV_ITEMS: NavItem[] = [
  {
    label: 'Settings',
    icon: Icons.Settings,
    path: '/settings',
  },
];
