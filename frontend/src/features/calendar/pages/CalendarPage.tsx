import { useCallback, useEffect, useMemo, useState } from "react";
import { Button, Loader, SegmentedControl } from "@mantine/core";
import {
  IconCalendarEvent,
  IconChevronLeft,
  IconChevronRight,
  IconClock,
  IconPlus,
  IconSend,
  IconAlertTriangle,
  IconCalendarStats,
} from "@tabler/icons-react";
import {
  addDays,
  addHours,
  addMonths,
  addWeeks,
  endOfMonth,
  endOfWeek,
  format,
  getDay,
  isToday,
  setHours,
  setMinutes,
  startOfMonth,
  startOfWeek,
  subMonths,
  subWeeks,
} from "date-fns";
import { enUS } from "date-fns/locale/en-US";
import { Calendar, dateFnsLocalizer, type View } from "react-big-calendar";
import "react-big-calendar/lib/css/react-big-calendar.css";
import { useChannels } from "@/features/channels/services/channelQueries";
import { SchedulePostModal } from "../components/SchedulePostModal";
import { useScheduledPostsRange } from "../hooks/useScheduledPostsQueries";
import type { CalendarViewMode, ScheduledPostListItem } from "../types";
import "./calendar-overrides.css";
import styles from "./CalendarPage.module.css";

const STORAGE_KEY = "posteaze:calendar-prefs";

interface CalendarPrefs {
  view: CalendarViewMode;
}

function loadPrefs(): CalendarPrefs {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) {
      const parsed = JSON.parse(raw) as Partial<CalendarPrefs>;
      if (parsed.view === "month" || parsed.view === "week" || parsed.view === "day") {
        return { view: parsed.view };
      }
    }
  } catch {
    /* ignore corrupted data */
  }
  return { view: "month" };
}

function savePrefs(prefs: CalendarPrefs) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(prefs));
  } catch {
    /* quota exceeded, etc. */
  }
}

const locales = { "en-US": enUS };

const localizer = dateFnsLocalizer({
  format,
  startOfWeek,
  getDay,
  locales,
});

function rangeStrings(view: CalendarViewMode, anchor: Date): { from: string; to: string } {
  if (view === "month") {
    const s = startOfMonth(anchor);
    const e = endOfMonth(anchor);
    return { from: format(s, "yyyy-MM-dd"), to: format(addDays(e, 1), "yyyy-MM-dd") };
  }
  if (view === "week") {
    const s = startOfWeek(anchor, { weekStartsOn: 0 });
    const e = endOfWeek(anchor, { weekStartsOn: 0 });
    return { from: format(s, "yyyy-MM-dd"), to: format(addDays(e, 1), "yyyy-MM-dd") };
  }
  return { from: format(anchor, "yyyy-MM-dd"), to: format(addDays(anchor, 1), "yyyy-MM-dd") };
}

function getNavLabel(view: CalendarViewMode, anchor: Date): string {
  if (view === "month") return format(anchor, "MMMM yyyy");
  if (view === "week") {
    const s = startOfWeek(anchor, { weekStartsOn: 0 });
    const e = endOfWeek(anchor, { weekStartsOn: 0 });
    return `${format(s, "MMM d")} – ${format(e, "MMM d, yyyy")}`;
  }
  return format(anchor, "EEEE, MMMM d, yyyy");
}

type CalEvent = {
  id: number;
  title: string;
  start: Date;
  end: Date;
  resource: ScheduledPostListItem;
};

const SEG_CLASS_NAMES = {
  root: styles.segRoot,
  indicator: styles.segIndicator,
  control: styles.segControl,
  label: styles.segLabel,
} as const;

const SCROLL_TO_TIME = setMinutes(setHours(new Date(), 7), 0);
const MIN_TIME = setMinutes(setHours(new Date(), 0), 0);
const MAX_TIME = setMinutes(setHours(new Date(), 23), 59);

export default function CalendarPage() {
  const [viewMode, setViewMode] = useState<CalendarViewMode>(() => loadPrefs().view);
  const [anchorDate, setAnchorDate] = useState(() => new Date());
  const [modalOpen, setModalOpen] = useState(false);
  const [slotStart, setSlotStart] = useState<Date | null>(null);

  useEffect(() => {
    savePrefs({ view: viewMode });
  }, [viewMode]);

  const { from, to } = useMemo(() => rangeStrings(viewMode, anchorDate), [viewMode, anchorDate]);

  const { data, isLoading, isError } = useScheduledPostsRange(from, to);
  const { data: channels = [], isLoading: chLoading } = useChannels();

  const rbcView: View = viewMode === "month" ? "month" : viewMode === "week" ? "week" : "day";

  const events: CalEvent[] = useMemo(() => {
    const posts = data?.posts ?? [];
    return posts.map((p) => {
      const start = new Date(p.scheduled_at);
      return {
        id: p.id,
        title: (p.caption?.trim() || p.post_type).slice(0, 48),
        start,
        end: addHours(start, 1),
        resource: p,
      };
    });
  }, [data?.posts]);

  const stats = useMemo(() => {
    const posts = data?.posts ?? [];
    const total = posts.length;
    const scheduled = posts.filter((p) => p.status === "scheduled").length;
    const published = posts.filter((p) => p.status === "published").length;
    const failed = posts.filter((p) => p.status === "failed" || p.status === "partial_failure").length;
    return { total, scheduled, published, failed };
  }, [data?.posts]);

  const navLabel = useMemo(() => getNavLabel(viewMode, anchorDate), [viewMode, anchorDate]);

  const goBack = useCallback(() => {
    setAnchorDate((d) => {
      if (viewMode === "month") return subMonths(d, 1);
      if (viewMode === "week") return subWeeks(d, 1);
      return addDays(d, -1);
    });
  }, [viewMode]);

  const goForward = useCallback(() => {
    setAnchorDate((d) => {
      if (viewMode === "month") return addMonths(d, 1);
      if (viewMode === "week") return addWeeks(d, 1);
      return addDays(d, 1);
    });
  }, [viewMode]);

  const goToday = useCallback(() => setAnchorDate(new Date()), []);

  const onSelectSlot = useCallback(({ start }: { start: Date }) => {
    setSlotStart(start);
    setModalOpen(true);
  }, []);

  const openNewPost = useCallback(() => {
    setSlotStart(null);
    setModalOpen(true);
  }, []);

  const isTimeView = viewMode === "week" || viewMode === "day";

  const isTodayAnchor =
    viewMode === "day"
      ? isToday(anchorDate)
      : viewMode === "week"
        ? isToday(startOfWeek(anchorDate, { weekStartsOn: 0 }))
        : anchorDate.getMonth() === new Date().getMonth() &&
          anchorDate.getFullYear() === new Date().getFullYear();

  return (
    <div className={styles.page}>
      {/* ─── Header ─── */}
      <div className={styles.header}>
        <div className={styles.headerLeft}>
          <span className={styles.kicker}>
            <span className={styles.kickerDot} />
            Content Studio
          </span>
          <h1 className={styles.pageTitle}>Calendar</h1>
          <p className={styles.pageSubtitle}>
            Plan, schedule, and visualize your content pipeline across all channels.
          </p>
        </div>
        <div className={styles.headerRight}>
          <SegmentedControl
            value={viewMode}
            onChange={(v) => setViewMode(v as CalendarViewMode)}
            classNames={SEG_CLASS_NAMES}
            data={[
              { label: "Month", value: "month" },
              { label: "Week", value: "week" },
              { label: "Day", value: "day" },
            ]}
          />
          <Button
            leftSection={<IconPlus size={16} />}
            onClick={openNewPost}
            className={styles.newPostBtn}
          >
            New post
          </Button>
        </div>
      </div>

      {/* ─── Stats ribbon ─── */}
      <div className={styles.statsRow}>
        <div className={styles.statCard}>
          <div className={styles.statIconBlue}>
            <IconCalendarEvent size={20} />
          </div>
          <div className={styles.statBody}>
            <span className={styles.statValue}>{stats.total}</span>
            <span className={styles.statLabel}>Total Posts</span>
          </div>
        </div>
        <div className={styles.statCard}>
          <div className={styles.statIconAmber}>
            <IconClock size={20} />
          </div>
          <div className={styles.statBody}>
            <span className={styles.statValue}>{stats.scheduled}</span>
            <span className={styles.statLabel}>Scheduled</span>
          </div>
        </div>
        <div className={styles.statCard}>
          <div className={styles.statIconGreen}>
            <IconSend size={20} />
          </div>
          <div className={styles.statBody}>
            <span className={styles.statValue}>{stats.published}</span>
            <span className={styles.statLabel}>Published</span>
          </div>
        </div>
        <div className={styles.statCard}>
          <div className={styles.statIconPurple}>
            <IconAlertTriangle size={20} />
          </div>
          <div className={styles.statBody}>
            <span className={styles.statValue}>{stats.failed}</span>
            <span className={styles.statLabel}>Failed</span>
          </div>
        </div>
      </div>

      {/* ─── Navigation bar ─── */}
      <div className={styles.navRow}>
        <div className={styles.navCenter}>
          <button type="button" className={styles.navBtn} onClick={goBack} aria-label="Previous">
            <IconChevronLeft size={18} />
          </button>
          <span className={styles.navLabel}>{navLabel}</span>
          <button type="button" className={styles.navBtn} onClick={goForward} aria-label="Next">
            <IconChevronRight size={18} />
          </button>
        </div>
        <div className={styles.navRight}>
          {!isTodayAnchor && (
            <button type="button" className={styles.todayBtn} onClick={goToday}>
              Today
            </button>
          )}
        </div>
      </div>

      {/* ─── Calendar grid ─── */}
      <div className={isTimeView ? styles.calendarPanelTime : styles.calendarPanel}>
        {isLoading || chLoading ? (
          <div className={styles.loadingWrap}>
            <Loader size="md" color="var(--pe-accent)" />
            <span className={styles.loadingText}>Loading your content…</span>
          </div>
        ) : isError ? (
          <div className={styles.errorWrap}>
            <IconCalendarStats size={40} color="var(--pe-text-muted)" />
            <span className={styles.errorTitle}>Unable to load posts</span>
            <span className={styles.errorDesc}>
              Something went wrong while fetching your scheduled content. Try refreshing.
            </span>
          </div>
        ) : (
          <div className={`posteaze-rbc-wrap ${isTimeView ? styles.calendarBodyTime : styles.calendarBody}`}>
            <Calendar
              culture="en-US"
              localizer={localizer}
              events={events}
              startAccessor="start"
              endAccessor="end"
              view={rbcView}
              date={anchorDate}
              onNavigate={setAnchorDate}
              onView={(v: View) => {
                if (v === "month") setViewMode("month");
                else if (v === "week") setViewMode("week");
                else if (v === "day") setViewMode("day");
              }}
              selectable
              onSelectSlot={onSelectSlot}
              views={["month", "week", "day"]}
              components={{ toolbar: () => null }}
              style={{ height: viewMode === "month" ? 560 : 680 }}
              step={30}
              timeslots={2}
              min={MIN_TIME}
              max={MAX_TIME}
              scrollToTime={SCROLL_TO_TIME}
              dayLayoutAlgorithm="overlap"
              eventPropGetter={() => ({
                style: {
                  borderRadius: viewMode === "month" ? 5 : 4,
                  border: "none",
                  fontSize: viewMode === "month" ? 11 : 11,
                  fontWeight: 500,
                  padding: viewMode === "month" ? "2px 6px" : "4px 8px 4px 6px",
                },
              })}
            />
          </div>
        )}
      </div>

      <SchedulePostModal
        opened={modalOpen}
        onClose={() => setModalOpen(false)}
        initialStart={slotStart}
        channels={channels}
      />
    </div>
  );
}
