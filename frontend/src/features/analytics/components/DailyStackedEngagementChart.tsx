import { useMemo } from "react";
import { Bar } from "react-chartjs-2";
import { Paper, Text, Skeleton } from "@mantine/core";
import { registerChartJs, useChartTheme } from "../chartSetup";
import type { DailyEngagementPoint } from "../types";
import styles from "./ChartCard.module.css";

registerChartJs();

type Props = {
  series: DailyEngagementPoint[];
  loading?: boolean;
};

export function DailyStackedEngagementChart({ series, loading }: Props) {
  const colors = useChartTheme();

  const data = useMemo(
    () => ({
      labels: series.map((r) => r.date),
      datasets: [
        {
          label: "Likes",
          data: series.map((r) => r.likes),
          backgroundColor: "rgba(29, 78, 216, 0.88)",
          stack: "e",
          borderRadius: 4,
        },
        {
          label: "Comments",
          data: series.map((r) => r.comments),
          backgroundColor: "rgba(124, 58, 237, 0.88)",
          stack: "e",
          borderRadius: 4,
        },
        {
          label: "Shares",
          data: series.map((r) => r.shares),
          backgroundColor: "rgba(234, 88, 12, 0.88)",
          stack: "e",
          borderRadius: 4,
        },
        {
          label: "Saves",
          data: series.map((r) => r.saves),
          backgroundColor: "rgba(13, 148, 136, 0.88)",
          stack: "e",
          borderRadius: 4,
        },
      ],
    }),
    [series]
  );

  const options = useMemo(
    () => ({
      responsive: true,
      maintainAspectRatio: false,
      interaction: { mode: "index" as const, intersect: false },
      plugins: {
        legend: {
          position: "bottom" as const,
          labels: { color: colors.text, boxWidth: 10, usePointStyle: true },
        },
        tooltip: { mode: "index" as const, intersect: false },
      },
      scales: {
        x: {
          stacked: true,
          grid: { display: false },
          ticks: { color: colors.text, maxTicksLimit: 10, maxRotation: 45 },
        },
        y: {
          stacked: true,
          beginAtZero: true,
          grid: { color: colors.grid },
          ticks: { color: colors.text },
        },
      },
    }),
    [colors]
  );

  return (
    <Paper className={styles.wrap} p="md" radius="md" withBorder>
      <Text fw={600} size="sm" mb={4} c="var(--pe-text)">
        Daily engagement stack
      </Text>
      <Text size="xs" c="dimmed" mb="md">
        Sum of estimated likes, comments, shares, and saves gained each day
      </Text>
      {loading ? (
        <Skeleton height={240} radius="sm" />
      ) : series.length === 0 ? (
        <Text c="dimmed" size="sm">
          No data in this range.
        </Text>
      ) : (
        <div className={styles.chartTall}>
          <Bar data={data} options={options} />
        </div>
      )}
    </Paper>
  );
}
