/**
 * Icon Management
 * 
 * Centralized icon exports for the PostEaze application.
 * All icons are imported from @tabler/icons-react and re-exported here.
 * 
 * Usage:
 * import { Icons } from '@/app/theme';
 * <Icons.Home size={20} />
 */

import {
  // Navigation Icons
  IconHome,
  IconChartBar,
  IconCalendar,
  IconSettings,
  
  // Social Media Icons
  IconBrandInstagram,
  IconBrandFacebook,
  IconBrandYoutube,
  IconBrandTwitter,
  IconBrandLinkedin,
  
  // Action Icons
  IconPlus,
  IconLogout,
  IconRefresh,
  IconCheck,
  IconX,
  IconEdit,
  IconTrash,
  IconDownload,
  IconUpload,
  IconSearch,
  
  // UI Icons
  IconChevronDown,
  IconChevronRight,
  IconChevronLeft,
  IconChevronUp,
  IconArrowLeft,
  IconArrowRight,
  IconMenu2,
  
  // Status Icons
  IconAlertCircle,
  IconAlertTriangle,
  IconInfoCircle,
  IconCircleCheck,
  IconCircleX,
  
  // User Icons
  IconUser,
  IconMail,
  IconLock,
  IconEye,
  IconEyeOff,
  
  // Content Icons
  IconPhoto,
  IconVideo,
  IconFile,
  IconFileText,
  
  // Type re-export for icon component type
  type Icon,
} from '@tabler/icons-react';

// Organized icon exports
export const Icons = {
  // Navigation
  Home: IconHome,
  ChartBar: IconChartBar,
  Calendar: IconCalendar,
  Settings: IconSettings,
  
  // Social Media
  Instagram: IconBrandInstagram,
  Facebook: IconBrandFacebook,
  YouTube: IconBrandYoutube,
  Twitter: IconBrandTwitter,
  LinkedIn: IconBrandLinkedin,
  
  // Actions
  Plus: IconPlus,
  Logout: IconLogout,
  Refresh: IconRefresh,
  Check: IconCheck,
  X: IconX,
  Edit: IconEdit,
  Trash: IconTrash,
  Download: IconDownload,
  Upload: IconUpload,
  Search: IconSearch,
  
  // UI
  ChevronDown: IconChevronDown,
  ChevronRight: IconChevronRight,
  ChevronLeft: IconChevronLeft,
  ChevronUp: IconChevronUp,
  ArrowLeft: IconArrowLeft,
  ArrowRight: IconArrowRight,
  Menu: IconMenu2,
  
  // Status
  AlertCircle: IconAlertCircle,
  AlertTriangle: IconAlertTriangle,
  InfoCircle: IconInfoCircle,
  CheckCircle: IconCircleCheck,
  XCircle: IconCircleX,
  
  // User
  User: IconUser,
  Mail: IconMail,
  Lock: IconLock,
  Eye: IconEye,
  EyeOff: IconEyeOff,
  
  // Content
  Photo: IconPhoto,
  Video: IconVideo,
  File: IconFile,
  FileText: IconFileText,
} as const;

// Export icon component type for component typing
export type IconComponent = Icon;

// Default icon size
export const DEFAULT_ICON_SIZE = 20;

// Icon size presets
export const ICON_SIZES = {
  xs: 14,
  sm: 16,
  md: 20,
  lg: 24,
  xl: 28,
} as const;
