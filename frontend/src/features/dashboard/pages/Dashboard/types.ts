import type { IconComponent } from '@/app/theme/icons';
import type { ScheduledPostListItem } from '@/features/calendar/types';

export interface DashboardKpiStat {
    title: string;
    value: number;
    icon: IconComponent;
    loading: boolean;
    attention: boolean;
    subtitle?: string;
}

export interface DashboardPostMetrics {
    counts: {
        published: number;
        scheduled: number;
        failed: number;
        cancelled: number;
        other: number;
    };
    total: number;
    upcomingCount: number;
    upcomingRows: ScheduledPostListItem[];
    publishedLast7dLocal: number;
    healthPct: number | null;
}

export interface DashboardChannelCard {
    name: string;
    icon: IconComponent;
    provider: 'instagram' | 'facebook' | 'youtube';
    path: string;
    connected: boolean;
}
