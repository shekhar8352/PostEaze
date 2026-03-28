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

export function FollowerGrowthChart({ series, loading }: Props) {
  const data = useMemo(() => {
    const labels = series.map((r) => r.date);
    const followers = series.map((r) => r.follower_count ?? null);
    return {
      labels,
      datasets: [
        {
          label: "Followers",
          data: followers as (number | null)[],
          borderColor: chartColors.accent,
          backgroundColor: "rgba(29, 78, 216, 0.12)",
          fill: true,
          tension: 0.35,
          spanGaps: true,
          pointRadius: 2,
          pointHoverRadius: 4,
        },
      ],
    };
  }, [series]);

  const options = useMemo(
    () => ({
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: false },
        tooltip: { mode: "index" as const, intersect: false },
      },
      scales: {
        x: {
          grid: { color: chartColors.grid },
          ticks: { color: chartColors.text, maxRotation: 45, minRotation: 0 },
        },
        y: {
          beginAtZero: false,
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
        Follower growth
      </Text>
      {loading ? (
        <Skeleton height={220} radius="sm" />
      ) : series.length === 0 ? (
        <Text c="dimmed" size="sm">
          No profile data in this range.
        </Text>
      ) : (
        <div className={styles.chart}>
          <Line data={data} options={options} />
        </div>
      )}
    </Paper>
  );
}
