import dayjs from "dayjs";

export type DateRangeStrings = { startDate: string; endDate: string };

/** Inclusive calendar-day window ending today (local), matching backend defaults. */
export function inclusiveRange(days: number): DateRangeStrings {
  const end = dayjs().startOf("day");
  const start = end.subtract(days - 1, "day");
  return {
    startDate: start.format("YYYY-MM-DD"),
    endDate: end.format("YYYY-MM-DD"),
  };
}
