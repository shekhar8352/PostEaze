import { useCallback, useMemo, useState } from "react";
import {
  Box,
  Button,
  Group,
  Loader,
  Paper,
  SegmentedControl,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import {
  addDays,
  addHours,
  endOfMonth,
  endOfWeek,
  format,
  getDay,
  startOfMonth,
  startOfWeek,
} from "date-fns";
import { enUS } from "date-fns/locale/en-US";
import { Calendar, dateFnsLocalizer, type View } from "react-big-calendar";
import "react-big-calendar/lib/css/react-big-calendar.css";
import { useChannels } from "@/features/channels/services/channelQueries";
import { SchedulePostModal } from "../components/SchedulePostModal";
import { useScheduledPostsRange } from "../hooks/useScheduledPostsQueries";
import type { CalendarViewMode, ScheduledPostListItem } from "../types";
import "./calendar-overrides.css";

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

type CalEvent = {
  id: number;
  title: string;
  start: Date;
  end: Date;
  resource: ScheduledPostListItem;
};

export default function CalendarPage() {
  const [viewMode, setViewMode] = useState<CalendarViewMode>("month");
  const [anchorDate, setAnchorDate] = useState(() => new Date());
  const [modalOpen, setModalOpen] = useState(false);
  const [slotStart, setSlotStart] = useState<Date | null>(null);

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

  const onSelectSlot = useCallback(({ start }: { start: Date }) => {
    setSlotStart(start);
    setModalOpen(true);
  }, []);

  const openNewPost = useCallback(() => {
    setSlotStart(null);
    setModalOpen(true);
  }, []);

  return (
    <Stack gap="lg" p={{ base: "sm", sm: "md" }} maw={1400} mx="auto">
      <Group justify="space-between" align="flex-start" wrap="wrap">
        <div>
          <Title order={2}>Calendar</Title>
          <Text size="sm" c="dimmed" mt={4}>
            Schedule Instagram posts. Month, week, and day views.
          </Text>
        </div>
        <Group gap="sm">
          <SegmentedControl
            value={viewMode}
            onChange={(v) => setViewMode(v as CalendarViewMode)}
            data={[
              { label: "Month", value: "month" },
              { label: "Week", value: "week" },
              { label: "Day", value: "day" },
            ]}
          />
          <Button onClick={openNewPost}>New post</Button>
        </Group>
      </Group>

      <Paper withBorder radius="md" p="xs" style={{ minHeight: 560 }}>
        {isLoading || chLoading ? (
          <Group justify="center" p="xl">
            <Loader />
          </Group>
        ) : isError ? (
          <Text c="red" p="md">
            Could not load scheduled posts.
          </Text>
        ) : (
          <Box className="posteaze-rbc-wrap">
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
              style={{ height: 520 }}
              eventPropGetter={() => ({
                style: {
                  borderRadius: 6,
                  border: "none",
                  boxShadow: "0 1px 2px rgba(0,0,0,0.06)",
                },
              })}
            />
          </Box>
        )}
      </Paper>

      <SchedulePostModal
        opened={modalOpen}
        onClose={() => setModalOpen(false)}
        initialStart={slotStart}
        channels={channels}
      />
    </Stack>
  );
}
