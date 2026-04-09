import { Card, Group, Skeleton, Text, Badge } from "@mantine/core";
import { formatDeltaPct } from "../utils/pctChange";
import styles from "./KPICard.module.css";

type KPICardProps = {
  title: string;
  value: string;
  deltaPct: number | null;
  loading?: boolean;
  /** Optional short hint when delta can’t be computed (e.g. no prior period) */
  deltaHint?: string;
};

export function KPICard({ title, value, deltaPct, loading, deltaHint }: KPICardProps) {
  if (loading) {
    return (
      <Card className={styles.card} padding="lg" radius="md" withBorder>
        <Skeleton height={10} width="42%" mb="sm" />
        <Skeleton height={32} width="55%" mb="md" />
        <Skeleton height={8} width="70%" />
      </Card>
    );
  }

  const positive = deltaPct != null && deltaPct > 0;
  const negative = deltaPct != null && deltaPct < 0;
  const neutral = deltaPct === 0;

  return (
    <Card className={styles.card} padding="lg" radius="md" withBorder>
      <Text size="xs" tt="uppercase" fw={600} className={styles.kicker} lineClamp={2}>
        {title}
      </Text>
      <Group justify="space-between" align="flex-start" wrap="nowrap" gap="sm" mt={8}>
        <Text className={styles.value} component="p" m={0}>
          {value}
        </Text>
        {deltaPct !== null && (
          <Badge
            size="md"
            variant={neutral ? "outline" : "light"}
            className={styles.deltaBadge}
            color={positive ? "teal" : negative ? "red" : "gray"}
          >
            {formatDeltaPct(deltaPct)}
          </Badge>
        )}
      </Group>
      <Group justify="space-between" align="center" mt="md" gap="xs" wrap="nowrap">
        <Text size="xs" className={styles.footerNote}>
          vs previous period
        </Text>
        {deltaPct === null && (
          <Text size="xs" className={styles.naHint}>
            {deltaHint ?? "No comparison"}
          </Text>
        )}
      </Group>
    </Card>
  );
}
