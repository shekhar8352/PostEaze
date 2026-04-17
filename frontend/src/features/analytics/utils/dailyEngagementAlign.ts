import type { DailyEngagementPoint } from "../types";
import { enumerateDaysInclusive } from "./dateRange";

/** Ensures every day in the selected range appears (missing API days → zeros). */
export function alignDailyEngagementSeries(
  startDate: string,
  endDate: string,
  series: DailyEngagementPoint[]
): DailyEngagementPoint[] {
  const map = new Map(series.map((p) => [p.date, p]));
  return enumerateDaysInclusive(startDate, endDate).map((date) => {
    const row = map.get(date);
    if (row) return row;
    return { date, likes: 0, comments: 0, shares: 0, saves: 0, total: 0 };
  });
}
