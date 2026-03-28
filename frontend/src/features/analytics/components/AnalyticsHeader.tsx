import { Group, Select, Button, Stack, Text } from "@mantine/core";
import { DatePickerInput } from "@mantine/dates";
import dayjs from "dayjs";
import type { InstagramChannelDisplay } from "@/features/channels/types/instagram.types";
import type { DateRangeStrings } from "../utils/dateRange";
import { inclusiveRange } from "../utils/dateRange";
import styles from "./AnalyticsHeader.module.css";

type Props = {
  channels: InstagramChannelDisplay[];
  channelId: number | null;
  onChannelChange: (id: number) => void;
  range: DateRangeStrings;
  onRangeChange: (range: DateRangeStrings) => void;
};

export function AnalyticsHeader({
  channels,
  channelId,
  onChannelChange,
  range,
  onRangeChange,
}: Props) {
  const selectData = channels.map((c) => ({
    value: String(c.channel_id),
    label: c.channelName || c.username || `Channel ${c.channel_id}`,
  }));

  const pickerValue: [Date | null, Date | null] = [
    dayjs(range.startDate).isValid() ? dayjs(range.startDate).toDate() : null,
    dayjs(range.endDate).isValid() ? dayjs(range.endDate).toDate() : null,
  ];

  const isPresetActive = (days: number) => {
    const p = inclusiveRange(days);
    return p.startDate === range.startDate && p.endDate === range.endDate;
  };

  return (
    <Stack gap="md" className={styles.header}>
      <div>
        <Text fw={700} size="xl" className={styles.title}>
          Instagram analytics
        </Text>
        <Text size="sm" c="dimmed">
          Channel and post performance for the selected period.
        </Text>
      </div>
      <Group justify="space-between" align="flex-end" wrap="wrap" gap="md">
        <Select
          label="Channel"
          placeholder="Select Instagram account"
          data={selectData}
          value={channelId != null ? String(channelId) : null}
          onChange={(v) => v && onChannelChange(Number(v))}
          searchable
          classNames={{ input: styles.selectInput }}
          w={{ base: "100%", sm: 280 }}
        />
        <div>
          <Text size="xs" fw={600} c="dimmed" mb={6}>
            Period
          </Text>
          <Group gap="xs" wrap="wrap">
            <Button.Group>
              {([7, 30, 90] as const).map((d) => (
                <Button
                  key={d}
                  variant={isPresetActive(d) ? "filled" : "default"}
                  size="sm"
                  onClick={() => onRangeChange(inclusiveRange(d))}
                >
                  {d}d
                </Button>
              ))}
            </Button.Group>
            <DatePickerInput
              type="range"
              value={pickerValue}
              onChange={(dates) => {
                const [a, b] = dates;
                if (a && b) {
                  const s = dayjs(a).startOf("day");
                  const e = dayjs(b).startOf("day");
                  const [start, end] = s.isAfter(e) ? [e, s] : [s, e];
                  onRangeChange({
                    startDate: start.format("YYYY-MM-DD"),
                    endDate: end.format("YYYY-MM-DD"),
                  });
                }
              }}
              maxDate={new Date()}
              size="sm"
              w={{ base: "100%", sm: 260 }}
              placeholder="Custom range"
            />
          </Group>
        </div>
      </Group>
    </Stack>
  );
}
