import { useMemo } from "react";
import { Line } from "react-chartjs-2";
import { Paper, Text, Skeleton } from "@mantine/core";
import { registerChartJs, useChartTheme } from "../chartSetup";
import type { DailyEngagementPoint } from "../types";
import styles from "./ChartCard.module.css";

registerChartJs();

type Props = {
  series: DailyEngagementPoint[];
  loading?: boolean;
};

/**
 * Engagement intensity: estimated daily interactions per 1k of summed reach proxy
 * (uses stacked total only — reach is account-level; this is a directional ratio for trends).
 * Simplified: rolling 3-day average of (likes+comments+shares+saves) for smoother reading.
 */
export function EngagementRateTrendChart({ series, loading }: Props) {
  const colors = useChartTheme();

  const { labels, smoothed } = useMemo(() => {
    const window = 3;
    const out: number[] = [];
    const totals = series.map((s) => s.total);
    for (let i = 0; i < totals.length; i++) {
      const from = Math.max(0, i - window + 1);
      const slice = totals.slice(from, i + 1);
      const avg = slice.reduce((a, b) => a + b, 0) / slice.length;
      out.push(Math.round(avg * 10) / 10);
    }
    return { labels: series.map((s) => s.date), smoothed: out };
  }, [series]);

  const data = useMemo(
    () => ({
      labels,
      datasets: [
        {
          label: "3-day avg total engagement",
          data: smoothed,
          borderColor: colors.accent5,
          backgroundColor: "rgba(219, 39, 119, 0.12)",
          tension: 0.4,
          fill: true,
          pointRadius: 0,
          pointHoverRadius: 4,
          borderWidth: 2,
        },
      ],
    }),
    [labels, smoothed, colors]
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
      },
      scales: {
        x: {
          grid: { display: false },
          ticks: { color: colors.text, maxTicksLimit: 8 },
        },
        y: {
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
        Engagement momentum
      </Text>
      <Text size="xs" c="dimmed" mb="md">
        Smoothed 3-day average of total estimated daily engagement (likes + comments + shares + saves)
      </Text>
      {loading ? (
        <Skeleton height={220} radius="sm" />
      ) : labels.length === 0 ? (
        <Text c="dimmed" size="sm">
          No data in this range.
        </Text>
      ) : (
        <div className={styles.chart}>
          <Line data={data} options={options} />
        </div>
      )}
    </Paper>
  );
}
