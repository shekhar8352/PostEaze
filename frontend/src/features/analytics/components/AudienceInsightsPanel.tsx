import { Paper, Skeleton, Stack, Text } from "@mantine/core";
import type { AudienceDashboard } from "../types";
import styles from "./AudienceInsightsPanel.module.css";
import chartStyles from "./ChartCard.module.css";

type Props = {
  audience?: AudienceDashboard | null;
  loading?: boolean;
};

export function AudienceInsightsPanel({ audience, loading }: Props) {
  const metrics = audience?.metrics ?? [];
  const hasData = metrics.length > 0;

  return (
    <Paper withBorder p="md" radius="md" className={`${chartStyles.wrap} ${styles.wrap}`}>
      <div className={styles.header}>
        <Text className={styles.title}>Audience composition</Text>
        {loading ? (
          <Skeleton height={14} width="40%" mt={6} />
        ) : hasData && audience?.snapshot_date ? (
          <Text className={styles.meta}>Snapshot · {audience.snapshot_date}</Text>
        ) : (
          <Text className={styles.meta}>Lifetime breakdown (Meta)</Text>
        )}
      </div>

      {loading ? (
        <Stack gap="sm">
          <Skeleton height={120} />
          <Skeleton height={120} />
        </Stack>
      ) : hasData ? (
        <div className={styles.grid}>
          {metrics.map((block) => (
            <section key={block.key} className={styles.block} aria-labelledby={`aud-${block.key}`}>
              <h3 id={`aud-${block.key}`} className={styles.blockTitle}>
                {block.title}
              </h3>
              <div className={styles.rows}>
                {block.items.map((item) => (
                  <div key={`${block.key}-${item.label}`} className={styles.row}>
                    <span className={styles.label} title={item.label}>
                      {item.label}
                    </span>
                    <span className={styles.value}>
                      {item.value.toLocaleString()} · {item.pct}%
                    </span>
                    <div className={styles.barTrack} role="presentation">
                      <div className={styles.barFill} style={{ width: `${Math.min(100, item.pct)}%` }} />
                    </div>
                  </div>
                ))}
              </div>
            </section>
          ))}
        </div>
      ) : (
        <div className={styles.empty}>
          <Text className={styles.emptyText}>
            No audience breakdown yet. Meta only returns these insights for professional accounts with at least
            about 100 followers, and data can take a sync cycle to appear after you connect the channel.
          </Text>
        </div>
      )}
    </Paper>
  );
}
