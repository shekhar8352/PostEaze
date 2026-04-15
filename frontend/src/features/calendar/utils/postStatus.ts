/** Visual bucket for calendar styling (maps backend `status` text). */
export type CalendarPostStatusKey = "published" | "scheduled" | "failed" | "cancelled" | "other";

export function getCalendarPostStatusKey(status: string): CalendarPostStatusKey {
  const s = (status || "").toLowerCase();
  if (s === "published") return "published";
  if (s === "scheduled" || s === "pending" || s === "submitting") return "scheduled";
  if (s === "failed" || s === "partial_failure") return "failed";
  if (s === "cancelled") return "cancelled";
  return "other";
}

export function formatScheduledPostStatus(status: string): string {
  const s = (status || "").toLowerCase();
  switch (s) {
    case "published":
      return "Published";
    case "scheduled":
      return "Scheduled";
    case "pending":
      return "Pending";
    case "submitting":
      return "Publishing";
    case "failed":
      return "Failed";
    case "partial_failure":
      return "Partial fail";
    case "cancelled":
      return "Cancelled";
    default:
      return status ? status.charAt(0).toUpperCase() + status.slice(1) : "Unknown";
  }
}
