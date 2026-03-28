import { Card, Group, Skeleton, Text, Badge } from "@mantine/core";
import { formatDeltaPct } from "../utils/pctChange";
import styles from "./KPICard.module.css";

type KPICardProps = {
  title: string;
  value: string;
  deltaPct: number | null;
  loading?: boolean;
};

export function KPICard({ title, value, deltaPct, loading }: KPICardProps) {
  if (loading) {
    return (
      <Card className={styles.card} padding="lg" radius="md" withBorder>
        <Skeleton height={12} width="45%" mb="sm" />
        <Skeleton height={28} width="60%" />
      </Card>
    );
  }

  const positive = deltaPct != null && deltaPct > 0;
  const negative = deltaPct != null && deltaPct < 0;
  const neutral = deltaPct === 0;

  return (
    <Card className={styles.card} padding="lg" radius="md" withBorder>
      <Text size="xs" tt="uppercase" fw={600} c="dimmed" mb={6}>
        {title}
      </Text>
      <Group justify="space-between" align="flex-end" wrap="nowrap" gap="xs">
        <Text size="xl" fw={700} className={styles.value}>
          {value}
        </Text>
        {deltaPct !== null && (
          <Badge
            size="sm"
            variant="light"
            color={positive ? "teal" : negative ? "red" : neutral ? "gray" : "gray"}
          >
            {formatDeltaPct(deltaPct)}
          </Badge>
        )}
        {deltaPct === null && (
          <Badge size="sm" variant="light" color="gray">
            —
          </Badge>
        )}
      </Group>
      <Text size="xs" c="dimmed" mt={6}>
        vs previous period
      </Text>
    </Card>
  );
}
