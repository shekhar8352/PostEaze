import { addDays, endOfDay, isAfter, isBefore, parseISO, startOfDay } from 'date-fns';
import { getCalendarPostStatusKey } from '@/features/calendar/utils/postStatus';
import type { ScheduledPostListItem } from '@/features/calendar/types';
import type { DashboardPostMetrics } from '../types';
import { isInRollingDays } from './dateRolling';

export function computeDashboardPostMetrics(posts: ScheduledPostListItem[]): DashboardPostMetrics {
    const now = new Date();
    const horizon = addDays(now, 30);

    const counts = {
        published: 0,
        scheduled: 0,
        failed: 0,
        cancelled: 0,
        other: 0,
    };

    const upcomingAll: ScheduledPostListItem[] = [];
    let publishedLast7dLocal = 0;

    for (const p of posts) {
        const key = getCalendarPostStatusKey(p.status);
        counts[key === 'other' ? 'other' : key] += 1;

        if (key === 'published' && isInRollingDays(p.updated_at, 7)) {
            publishedLast7dLocal += 1;
        }

        if (key === 'scheduled') {
            const at = parseISO(p.scheduled_at);
            if (!isBefore(at, startOfDay(now)) && !isAfter(at, endOfDay(horizon))) {
                upcomingAll.push(p);
            }
        }
    }

    const upcomingRows = [...upcomingAll]
        .sort((a, b) => parseISO(a.scheduled_at).getTime() - parseISO(b.scheduled_at).getTime())
        .slice(0, 5);

    const total = posts.length;
    const healthDen = counts.published + counts.failed;
    const healthPct =
        healthDen > 0 ? Math.round((counts.published / healthDen) * 100) : total === 0 ? 100 : null;

    return {
        counts,
        total,
        upcomingCount: upcomingAll.length,
        upcomingRows,
        publishedLast7dLocal,
        healthPct,
    };
}
