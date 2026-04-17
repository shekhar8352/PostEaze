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

/** Every calendar day from start through end inclusive (YYYY-MM-DD). */
export function enumerateDaysInclusive(startDate: string, endDate: string): string[] {
  const start = dayjs(startDate);
  const end = dayjs(endDate);
  const out: string[] = [];
  for (let d = start; !d.isAfter(end, "day"); d = d.add(1, "day")) {
    out.push(d.format("YYYY-MM-DD"));
  }
  return out;
}
