import { useMemo } from 'react';
import { addDays, format, startOfDay, subDays } from 'date-fns';
import type { DateRangeParams } from '@/features/analytics/services/analyticsService';

export interface DashboardDateRanges {
    pipelineFrom: string;
    pipelineTo: string;
    analyticsRange: DateRangeParams;
}

export function useDashboardDateRanges(): DashboardDateRanges {
    return useMemo(() => {
        const pipelineEnd = addDays(new Date(), 60);
        const pipelineStart = subDays(new Date(), 30);
        const pipelineFrom = format(startOfDay(pipelineStart), 'yyyy-MM-dd');
        const pipelineTo = format(addDays(startOfDay(pipelineEnd), 1), 'yyyy-MM-dd');

        const analyticsEnd = new Date();
        const analyticsStart = subDays(analyticsEnd, 6);
        const analyticsRange: DateRangeParams = {
            startDate: format(analyticsStart, 'yyyy-MM-dd'),
            endDate: format(analyticsEnd, 'yyyy-MM-dd'),
        };

        return { pipelineFrom, pipelineTo, analyticsRange };
    }, []);
}
