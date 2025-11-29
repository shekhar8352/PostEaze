type ChannelType = 'instagram' | 'facebook' | 'youtube';

export const CHANNEL_TYPES = {
  INSTAGRAM: 'instagram' as const,
  FACEBOOK: 'facebook' as const,
  YOUTUBE: 'youtube' as const,
};

export interface NavItem {
  label: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  icon: any;
  path?: string;
  children?: NavItem[];
  badge?: string | number;
  onClick?: () => void;
}

export interface LayoutState {
  sidebarOpened: boolean;
  channelsExpanded: boolean;
}

export interface ConnectedChannel {
  type: ChannelType;
  id: string;
  name: string;
  isConnected: boolean;
}
