import {
  IconHome,
  IconBrandInstagram,
  IconBrandFacebook,
  IconBrandYoutube,
  IconSettings,
  IconChartBar,
  IconCalendar,
} from '@tabler/icons-react';
import { type NavItem } from './types';

// Sidebar dimensions
export const SIDEBAR_WIDTH = 280;
export const SIDEBAR_COLLAPSED_WIDTH = 80;

// Main navigation items
export const MAIN_NAV_ITEMS: NavItem[] = [
  {
    label: 'Dashboard',
    icon: IconHome,
    path: '/dashboard',
  },
  {
    label: 'Analytics',
    icon: IconChartBar,
    path: '/analytics',
  },
  {
    label: 'Calendar',
    icon: IconCalendar,
    path: '/calendar',
  },
];

// Channel navigation items
export const CHANNEL_NAV_ITEMS: NavItem[] = [
  {
    label: 'Instagram',
    icon: IconBrandInstagram,
    path: '/channels/instagram',
  },
  {
    label: 'Facebook',
    icon: IconBrandFacebook,
    path: '/channels/facebook',
  },
  {
    label: 'YouTube',
    icon: IconBrandYoutube,
    path: '/channels/youtube',
  },
];

// Bottom navigation items
export const BOTTOM_NAV_ITEMS: NavItem[] = [
  {
    label: 'Settings',
    icon: IconSettings,
    path: '/settings',
  },
];

// Channel colors for branding
export const CHANNEL_COLORS = {
  instagram: 'linear-gradient(45deg, #f09433 0%, #e6683c 25%, #dc2743 50%, #cc2366 75%, #bc1888 100%)',
  facebook: '#1877F2',
  youtube: '#FF0000',
};
