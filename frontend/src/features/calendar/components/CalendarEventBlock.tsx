import type { EventProps } from "react-big-calendar";
import type { CalendarScheduledEvent } from "../types";
import { formatScheduledPostStatus, getCalendarPostStatusKey } from "../utils/postStatus";
import styles from "./CalendarEventBlock.module.css";

export function CalendarEventBlock({ event }: EventProps<CalendarScheduledEvent>) {
  const status = event.resource?.status ?? "";
  const key = getCalendarPostStatusKey(status);
  return (
    <div className={styles.root}>
      <span className={`${styles.badge} cal-ev-badge--${key}`}>{formatScheduledPostStatus(status)}</span>
      <span className={styles.title}>{event.title}</span>
    </div>
  );
}
