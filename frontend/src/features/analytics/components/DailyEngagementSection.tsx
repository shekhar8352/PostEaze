import { useMemo } from "react";
import { SimpleGrid, Stack, Text, Alert } from "@mantine/core";
import { Icons } from "@/app/theme";
import type { DailyEngagementPoint, PostsOverview } from "../types";
import { DailyMetricLineChart } from "./DailyMetricLineChart";
import { DailyStackedEngagementChart } from "./DailyStackedEngagementChart";
import { PeriodEngagementMixChart } from "./PeriodEngagementMixChart";
import { EngagementRateTrendChart } from "./EngagementRateTrendChart";
import { alignDailyEngagementSeries } from "../utils/dailyEngagementAlign";

type Props = {
  startDate: string;
  endDate: string;
  series: DailyEngagementPoint[];
  overview: PostsOverview | undefined;
  note?: string;
  loading?: boolean;
};

export function DailyEngagementSection({
  startDate,
  endDate,
  series,
  overview,
  note,
  loading,
}: Props) {
  const aligned = useMemo(
    () => alignDailyEngagementSeries(startDate, endDate, series),
    [startDate, endDate, series]
  );

  const labels = useMemo(() => aligned.map((r) => r.date), [aligned]);

  return (
    <Stack gap="md">
      <div>
        <Text fw={700} size="lg" c="var(--pe-text)">
          Post engagement by day
        </Text>
        <Text size="sm" c="dimmed" maw={900}>
          Separate daily views for reactions and distribution so you can see which behaviors spike together.
        </Text>
      </div>

      {note ? (
        <Alert variant="light" color="gray" icon={<Icons.InfoCircle size={18} />} title="How to read this">
          {note}
        </Alert>
      ) : null}

      <SimpleGrid cols={{ base: 1, md: 3 }} spacing="md">
        <DailyMetricLineChart
          title="Likes per day"
          labels={labels}
          values={aligned.map((r) => r.likes)}
          borderColor="rgba(29, 78, 216, 1)"
          fillColor="rgba(29, 78, 216, 0.15)"
          loading={loading}
        />
        <DailyMetricLineChart
          title="Comments per day"
          labels={labels}
          values={aligned.map((r) => r.comments)}
          borderColor="rgba(124, 58, 237, 1)"
          fillColor="rgba(124, 58, 237, 0.15)"
          loading={loading}
        />
        <DailyMetricLineChart
          title="Shares per day"
          labels={labels}
          values={aligned.map((r) => r.shares)}
          borderColor="rgba(234, 88, 12, 1)"
          fillColor="rgba(234, 88, 12, 0.15)"
          loading={loading}
        />
      </SimpleGrid>

      <SimpleGrid cols={{ base: 1, md: 2 }} spacing="md">
        <DailyMetricLineChart
          title="Saves per day"
          labels={labels}
          values={aligned.map((r) => r.saves)}
          borderColor="rgba(13, 148, 136, 1)"
          fillColor="rgba(13, 148, 136, 0.15)"
          loading={loading}
        />
        <EngagementRateTrendChart series={aligned} loading={loading} />
      </SimpleGrid>

      <SimpleGrid cols={{ base: 1, lg: 2 }} spacing="md">
        <DailyStackedEngagementChart series={aligned} loading={loading} />
        <PeriodEngagementMixChart overview={overview} loading={loading} />
      </SimpleGrid>
    </Stack>
  );
}
