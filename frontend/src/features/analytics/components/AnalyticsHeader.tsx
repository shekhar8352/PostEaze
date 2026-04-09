import { Group, Select, Button, Stack, Text, SegmentedControl, Paper, Box, Tooltip } from "@mantine/core";
import { DatePickerInput } from "@mantine/dates";
import dayjs from "dayjs";
import type { BaseChannelDisplay, ChannelProvider } from "@/features/channels/types/base.types";
import type { DateRangeStrings } from "../utils/dateRange";
import { inclusiveRange } from "../utils/dateRange";
import { Icons } from "@/app/theme";
import styles from "./AnalyticsHeader.module.css";

type Props = {
  platform: ChannelProvider;
  onPlatformChange: (p: Extract<ChannelProvider, "instagram" | "facebook">) => void;
  hasInstagram: boolean;
  hasFacebook: boolean;
  channels: BaseChannelDisplay[];
  channelId: number | null;
  onChannelChange: (id: number) => void;
  range: DateRangeStrings;
  onRangeChange: (range: DateRangeStrings) => void;
  /** Pull latest insights from Meta for connected Instagram / Facebook channels */
  onSyncMeta?: () => void;
  syncing?: boolean;
  syncDisabled?: boolean;
};

export function AnalyticsHeader({
  platform,
  onPlatformChange,
  hasInstagram,
  hasFacebook,
  channels,
  channelId,
  onChannelChange,
  range,
  onRangeChange,
  onSyncMeta,
  syncing = false,
  syncDisabled = false,
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

  const rangeLabel = `${dayjs(range.startDate).format("MMM D, YYYY")} – ${dayjs(range.endDate).format("MMM D, YYYY")}`;

  const platformLabel =
    platform === "instagram" ? "Instagram" : platform === "facebook" ? "Facebook Page" : "Channel";

  return (
    <Paper radius="lg" p="xl" withBorder className={styles.hero}>
      <Stack gap="lg">
        <div className={styles.titleRow}>
          <div>
            <Text fw={800} size="xl" className={styles.title} lh={1.2}>
              Channel analytics
            </Text>
            <Text size="sm" mt={6} style={{ color: "var(--pe-text-muted)" }}>
              Performance and reach for <strong style={{ color: "var(--pe-text)" }}>{platformLabel}</strong> — pick an
              account and date range.
            </Text>
          </div>
          <Group gap="sm" wrap="wrap" justify="flex-end">
            {(hasInstagram || hasFacebook) && (
              <SegmentedControl
                size="sm"
                value={platform === "facebook" ? "facebook" : "instagram"}
                onChange={(v) => onPlatformChange(v as "instagram" | "facebook")}
                data={[
                  {
                    value: "instagram",
                    label: (
                      <Group gap={6} wrap="nowrap">
                        <Icons.Instagram size={16} />
                        Instagram
                      </Group>
                    ),
                    disabled: !hasInstagram,
                  },
                  {
                    value: "facebook",
                    label: (
                      <Group gap={6} wrap="nowrap">
                        <Icons.Facebook size={16} />
                        Facebook
                      </Group>
                    ),
                    disabled: !hasFacebook,
                  },
                ]}
                className={styles.segmented}
              />
            )}
            {onSyncMeta && (
              <Tooltip
                label="Fetches the latest Instagram and Facebook insights from Meta. Can take up to a minute."
                multiline
                w={280}
              >
                <Button
                  variant="light"
                  color="blue"
                  size="sm"
                  radius="md"
                  leftSection={<Icons.Refresh size={16} />}
                  loading={syncing}
                  disabled={syncDisabled || syncing}
                  onClick={onSyncMeta}
                >
                  Sync from Meta
                </Button>
              </Tooltip>
            )}
          </Group>
        </div>

        <Group justify="space-between" align="flex-end" wrap="wrap" gap="md">
          <Select
            label="Account"
            placeholder={`Select ${platformLabel.toLowerCase()}`}
            data={selectData}
            value={channelId != null ? String(channelId) : null}
            onChange={(v) => v && onChannelChange(Number(v))}
            searchable
            classNames={{ input: styles.selectInput, label: styles.fieldLabel }}
            w={{ base: "100%", sm: 300 }}
          />

          <Stack gap={6} className={styles.periodBlock}>
            <Text size="xs" fw={600} className={styles.fieldLabel}>
              Reporting period
            </Text>
            <Group gap="xs" wrap="wrap" align="center">
              <Button.Group>
                {([7, 30, 90] as const).map((d) => (
                  <Button
                    key={d}
                    variant={isPresetActive(d) ? "filled" : "default"}
                    size="sm"
                    radius="md"
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
                radius="md"
                placeholder="Custom range"
                w={{ base: "100%", sm: 280 }}
                classNames={{ input: styles.dateInput }}
                styles={{
                  input: {
                    backgroundColor: "var(--pe-bg-elevated, #fff)",
                    color: "var(--pe-text, inherit)",
                    borderColor: "var(--pe-border)",
                  },
                }}
              />
            </Group>
            <Box className={styles.rangePill}>
              <Text size="xs" component="span" className={styles.rangeLabel}>
                Selected:{" "}
              </Text>
              <Text size="xs" fw={600} component="span" className={styles.rangeText}>
                {rangeLabel}
              </Text>
            </Box>
          </Stack>
        </Group>
      </Stack>
    </Paper>
  );
}
