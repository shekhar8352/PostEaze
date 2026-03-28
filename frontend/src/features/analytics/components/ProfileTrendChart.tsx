import { useMemo } from "react";
import { Line } from "react-chartjs-2";
import { Paper, Text, Skeleton } from "@mantine/core";
import { registerChartJs, chartColors } from "../chartSetup";
import type { ProfileAnalyticsItem } from "../types";
import styles from "./ChartCard.module.css";

registerChartJs();

type Props = {
  series: ProfileAnalyticsItem[];
  loading?: boolean;
};

export function ProfileTrendChart({ series, loading }: Props) {
  const data = useMemo(() => {
    const labels = series.map((r) => r.date);
    return {
      labels,
      datasets: [
        {
          label: "Reach",
          data: series.map((r) => r.reach ?? null),
          borderColor: chartColors.accent,
          backgroundColor: chartColors.accent,
          tension: 0.35,
          pointRadius: 0,
          pointHoverRadius: 4,
          spanGaps: true,
        },
        {
          label: "Impressions",
          data: series.map((r) => r.impressions ?? null),
          borderColor: chartColors.accent2,
          backgroundColor: chartColors.accent2,
          tension: 0.35,
          pointRadius: 0,
          pointHoverRadius: 4,
          spanGaps: true,
        },
        {
          label: "Profile views",
          data: series.map((r) => r.profile_views ?? null),
          borderColor: chartColors.accent3,
          backgroundColor: chartColors.accent3,
          tension: 0.35,
          pointRadius: 0,
          pointHoverRadius: 4,
          spanGaps: true,
        },
        {
          label: "Website clicks",
          data: series.map((r) => r.website_clicks ?? null),
          borderColor: chartColors.accent4,
          backgroundColor: chartColors.accent4,
          tension: 0.35,
          pointRadius: 0,
          pointHoverRadius: 4,
          spanGaps: true,
        },
        {
          label: "Accounts engaged",
          data: series.map((r) => r.accounts_engaged ?? null),
          borderColor: chartColors.accent5,
          backgroundColor: chartColors.accent5,
          tension: 0.35,
          pointRadius: 0,
          pointHoverRadius: 4,
          spanGaps: true,
        },
        {
          label: "Total interactions",
          data: series.map((r) => r.total_interactions ?? null),
          borderColor: chartColors.accent6,
          backgroundColor: chartColors.accent6,
          tension: 0.35,
          pointRadius: 0,
          pointHoverRadius: 4,
          spanGaps: true,
        },
      ],
    };
  }, [series]);

  const options = useMemo(
    () => ({
      responsive: true,
      maintainAspectRatio: false,
      interaction: { mode: "index" as const, intersect: false },
      plugins: {
        legend: {
          position: "bottom" as const,
          labels: { color: chartColors.text, boxWidth: 10, usePointStyle: true },
        },
        tooltip: { mode: "index" as const, intersect: false },
      },
      scales: {
        x: {
          grid: { color: chartColors.grid },
          ticks: { color: chartColors.text, maxRotation: 45 },
        },
        y: {
          beginAtZero: true,
          grid: { color: chartColors.grid },
          ticks: { color: chartColors.text },
        },
      },
    }),
    []
  );

  return (
    <Paper className={styles.wrap} p="md" radius="md" withBorder>
      <Text fw={600} size="sm" mb="md" c="var(--pe-text)">
        Account reach, activity & engagement
      </Text>
      {loading ? (
        <Skeleton height={280} radius="sm" />
      ) : series.length === 0 ? (
        <Text c="dimmed" size="sm">
          No profile data in this range.
        </Text>
      ) : (
        <div className={styles.chartTall}>
          <Line data={data} options={options} />
        </div>
      )}
    </Paper>
  );
}
