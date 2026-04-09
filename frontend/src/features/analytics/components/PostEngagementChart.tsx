import { useMemo } from "react";
import { Bar } from "react-chartjs-2";
import { Paper, Text, Skeleton } from "@mantine/core";
import { registerChartJs, useChartTheme } from "../chartSetup";
import type { PostsOverview } from "../types";
import styles from "./ChartCard.module.css";

registerChartJs();

type Props = {
  overview: PostsOverview | undefined;
  loading?: boolean;
};

export function PostEngagementChart({ overview, loading }: Props) {
  const colors = useChartTheme();

  const data = useMemo(() => {
    const o = overview;
    if (!o) {
      return {
        labels: ["Likes", "Comments", "Saves", "Shares"],
        datasets: [{ label: "Totals", data: [0, 0, 0, 0], backgroundColor: colors.accent }],
      };
    }
    return {
      labels: ["Likes", "Comments", "Saves", "Shares"],
      datasets: [
        {
          label: "Engagement (range)",
          data: [o.total_likes, o.total_comments, o.total_saves, o.total_shares],
          backgroundColor: [
            "rgba(29, 78, 216, 0.85)",
            "rgba(124, 58, 237, 0.85)",
            "rgba(13, 148, 136, 0.85)",
            "rgba(234, 88, 12, 0.85)",
          ],
          borderRadius: 6,
        },
      ],
    };
  }, [overview, colors]);

  const options = useMemo(
    () => ({
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: false },
      },
      scales: {
        x: {
          grid: { display: false },
          ticks: { color: colors.text },
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
      <Text fw={600} size="sm" mb="md" c="var(--pe-text)">
        Post engagement totals
      </Text>
      {loading ? (
        <Skeleton height={220} radius="sm" />
      ) : (
        <div className={styles.chart}>
          <Bar data={data} options={options} />
        </div>
      )}
    </Paper>
  );
}
