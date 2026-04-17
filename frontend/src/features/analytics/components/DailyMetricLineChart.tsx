import { useMemo } from "react";
import { Line } from "react-chartjs-2";
import { Paper, Text, Skeleton } from "@mantine/core";
import { registerChartJs, useChartTheme } from "../chartSetup";
import styles from "./ChartCard.module.css";

registerChartJs();

type Props = {
  title: string;
  subtitle?: string;
  labels: string[];
  values: number[];
  borderColor: string;
  fillColor?: string;
  loading?: boolean;
};

export function DailyMetricLineChart({
  title,
  subtitle,
  labels,
  values,
  borderColor,
  fillColor,
  loading,
}: Props) {
  const colors = useChartTheme();
  const fill = fillColor ?? `${borderColor}33`;

  const data = useMemo(
    () => ({
      labels,
      datasets: [
        {
          label: title,
          data: values,
          borderColor,
          backgroundColor: fill,
          tension: 0.35,
          fill: true,
          pointRadius: 0,
          pointHoverRadius: 4,
          borderWidth: 2,
        },
      ],
    }),
    [labels, values, title, borderColor, fill]
  );

  const options = useMemo(
    () => ({
      responsive: true,
      maintainAspectRatio: false,
      interaction: { mode: "index" as const, intersect: false },
      plugins: {
        legend: { display: false },
        tooltip: { mode: "index" as const, intersect: false },
      },
      scales: {
        x: {
          grid: { display: false },
          ticks: { color: colors.text, maxTicksLimit: 8, maxRotation: 0 },
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
      <Text fw={600} size="sm" c="var(--pe-text)">
        {title}
      </Text>
      {subtitle ? (
        <Text size="xs" c="dimmed" mb="sm" mt={4}>
          {subtitle}
        </Text>
      ) : (
        <Text size="xs" c="dimmed" mb="sm" mt={4}>
          Estimated net change vs prior snapshot
        </Text>
      )}
      {loading ? (
        <Skeleton height={200} radius="sm" />
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
