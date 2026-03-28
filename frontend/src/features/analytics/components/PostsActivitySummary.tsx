import { Paper, SimpleGrid, Text, Skeleton } from "@mantine/core";
import { formatCompact } from "../utils/pctChange";
import type { PostsOverview } from "../types";
import styles from "./PostsActivitySummary.module.css";

type Props = {
  overview: PostsOverview | undefined;
  loading?: boolean;
};

const items = (o: PostsOverview | undefined) => [
  { label: "New posts", value: o?.new_posts ?? 0 },
  { label: "Total posts (range)", value: o?.total_posts ?? 0 },
  { label: "Total reach (posts)", value: formatCompact(o?.total_reach ?? 0) },
  { label: "Total impressions", value: formatCompact(o?.total_impressions ?? 0) },
];

export function PostsActivitySummary({ overview, loading }: Props) {
  return (
    <Paper className={styles.wrap} p="md" radius="md" withBorder>
      <Text fw={600} size="sm" mb="md" c="var(--pe-text)">
        Posts activity
      </Text>
      <SimpleGrid cols={2} spacing="sm">
        {items(overview).map((item) => (
          <div key={item.label} className={styles.cell}>
            <Text size="xs" c="dimmed" tt="uppercase" fw={600}>
              {item.label}
            </Text>
            {loading ? (
              <Skeleton height={22} mt={6} width="50%" />
            ) : (
              <Text fw={700} size="lg" mt={4} className={styles.value}>
                {typeof item.value === "number" ? item.value.toLocaleString() : item.value}
              </Text>
            )}
          </div>
        ))}
      </SimpleGrid>
    </Paper>
  );
}
