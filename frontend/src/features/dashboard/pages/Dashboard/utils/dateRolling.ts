import { endOfDay, isAfter, isBefore, parseISO, startOfDay, subDays } from 'date-fns';

export function isInRollingDays(iso: string, days: number): boolean {
    try {
        const d = parseISO(iso);
        const start = startOfDay(subDays(new Date(), days - 1));
        const end = endOfDay(new Date());
        return !isBefore(d, start) && !isAfter(d, end);
    } catch {
        return false;
    }
}
