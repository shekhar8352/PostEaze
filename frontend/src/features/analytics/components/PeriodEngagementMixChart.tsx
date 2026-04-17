import { useMemo } from "react";
import { Doughnut } from "react-chartjs-2";
import { Paper, Text, Skeleton } from "@mantine/core";
import { registerChartJs, useChartTheme } from "../chartSetup";
import type { PostsOverview } from "../types";
import styles from "./ChartCard.module.css";

registerChartJs();

type Props = {
  overview: PostsOverview | undefined;
  loading?: boolean;
};

function clamp0(n: number) {
  return Math.max(0, n);
}

export function PeriodEngagementMixChart({ overview, loading }: Props) {
  const colors = useChartTheme();

  const { data, empty } = useMemo(() => {
    const o = overview;
    if (!o) {
      return {
        empty: true,
        data: {
          labels: ["Likes", "Comments", "Shares", "Saves"],
          datasets: [{ data: [0, 0, 0, 0], backgroundColor: ["#ccc", "#ccc", "#ccc", "#ccc"] }],
        },
      };
    }
    const likes = clamp0(o.new_likes);
    const comments = clamp0(o.new_comments);
    const shares = clamp0(o.new_shares);
    const saves = clamp0(o.new_saves);
    const sum = likes + comments + shares + saves;
    return {
      empty: sum === 0,
      data: {
        labels: ["Net likes", "Net comments", "Net shares", "Net saves"],
        datasets: [
          {
            data: [likes, comments, shares, saves],
            backgroundColor: [
              "rgba(29, 78, 216, 0.9)",
              "rgba(124, 58, 237, 0.9)",
              "rgba(234, 88, 12, 0.9)",
              "rgba(13, 148, 136, 0.9)",
            ],
            borderWidth: 0,
            hoverOffset: 6,
          },
        ],
      },
    };
  }, [overview]);

  const options = useMemo(
    () => ({
      responsive: true,
      maintainAspectRatio: false,
      cutout: "58%",
      plugins: {
        legend: {
          position: "right" as const,
          labels: { color: colors.text, boxWidth: 10, usePointStyle: true },
        },
        tooltip: {
          callbacks: {
            label: (ctx: { label?: string; parsed: number }) => {
              const v = ctx.parsed;
              const label = ctx.label ?? "";
              return `${label}: ${v.toLocaleString()}`;
            },
          },
        },
      },
    }),
    [colors]
  );

  return (
    <Paper className={styles.wrap} p="md" radius="md" withBorder>
      <Text fw={600} size="sm" mb={4} c="var(--pe-text)">
        Period mix — net engagement
      </Text>
      <Text size="xs" c="dimmed" mb="md">
        Share of likes, comments, shares, and saves gained vs the prior period of equal length
      </Text>
      {loading ? (
        <Skeleton height={240} radius="sm" />
      ) : empty ? (
        <Text c="dimmed" size="sm">
          Not enough change between periods to chart — try a longer range or sync latest data.
        </Text>
      ) : (
        <div className={styles.chartTall}>
          <Doughnut data={data} options={options} />
        </div>
      )}
    </Paper>
  );
}
