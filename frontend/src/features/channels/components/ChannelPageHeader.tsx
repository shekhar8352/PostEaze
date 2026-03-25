import { Paper, Group, Title, Text, Button } from '@mantine/core';
import type { ReactNode } from 'react';
import styles from './ChannelPageHeader.module.css';

export type ChannelBrand = 'instagram' | 'facebook' | 'youtube';

const ACCENT: Record<ChannelBrand, { border: string; iconBg: string; iconColor: string }> = {
  instagram: {
    border: '#E1306C',
    iconBg: 'rgba(225, 48, 108, 0.1)',
    iconColor: '#C13584',
  },
  facebook: {
    border: '#1877F2',
    iconBg: 'rgba(24, 119, 242, 0.1)',
    iconColor: '#1877F2',
  },
  youtube: {
    border: '#CC0000',
    iconBg: 'rgba(204, 0, 0, 0.08)',
    iconColor: '#CC0000',
  },
};

export type ChannelPageHeaderProps = {
  brand: ChannelBrand;
  title: string;
  description: string;
  icon: ReactNode;
  actionLabel: string;
  onAction?: () => void;
  actionLeftSection?: ReactNode;
};

export function ChannelPageHeader({
  brand,
  title,
  description,
  icon,
  actionLabel,
  onAction,
  actionLeftSection,
}: ChannelPageHeaderProps) {
  const a = ACCENT[brand];

  return (
    <Paper
      className={styles.paper}
      style={{ borderLeftColor: a.border }}
      shadow="none"
      radius="lg"
    >
      <Group justify="space-between" align="center" wrap="wrap" gap="md">
        <Group gap="md" wrap="nowrap">
          <div
            className={styles.iconWrap}
            style={{ backgroundColor: a.iconBg, color: a.iconColor }}
            aria-hidden
          >
            {icon}
          </div>
          <div>
            <Title order={2} className={styles.title}>
              {title}
            </Title>
            <Text className={styles.subtitle}>{description}</Text>
          </div>
        </Group>
        {onAction && (
          <Button
            leftSection={actionLeftSection}
            size="md"
            radius="md"
            variant="filled"
            color="blue"
            onClick={onAction}
          >
            {actionLabel}
          </Button>
        )}
      </Group>
    </Paper>
  );
}
