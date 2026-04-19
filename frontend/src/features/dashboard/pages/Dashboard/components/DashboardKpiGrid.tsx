import { Card, Group, Skeleton, Text } from '@mantine/core';
import type { DashboardKpiStat } from '../types';
import dashStyles from '../DashboardPage.module.css';

interface DashboardKpiGridProps {
    stats: DashboardKpiStat[];
}

export function DashboardKpiGrid({ stats }: DashboardKpiGridProps) {
    return (
        <div className={dashStyles.kpiGrid}>
            {stats.map((stat) => {
                const Icon = stat.icon;
                return (
                    <Card
                        key={stat.title}
                        padding="lg"
                        radius="md"
                        withBorder={false}
                        className={`${dashStyles.statCard} ${dashStyles.statCardRoot} ${stat.attention ? dashStyles.statCardAttention : ''}`}
                    >
                        <Group justify="space-between" align="flex-start" wrap="wrap" gap="sm">
                            <div style={{ minWidth: 0, flex: '1 1 8rem' }}>
                                <Text size="xs" tt="uppercase" fw={600} c="dimmed" mb={6}>
                                    {stat.title}
                                </Text>
                                {stat.subtitle ? (
                                    <Text size="xs" c="dimmed" mb={4} lineClamp={2}>
                                        {stat.subtitle}
                                    </Text>
                                ) : null}
                                {stat.loading ? (
                                    <Skeleton height={28} width={48} mt={4} />
                                ) : (
                                    <Text size="xl" fw={700} c="var(--pe-text)">
                                        {stat.value}
                                    </Text>
                                )}
                            </div>
                            <div className={dashStyles.statIcon} style={{ flexShrink: 0 }}>
                                <Icon size={22} stroke={1.75} />
                            </div>
                        </Group>
                    </Card>
                );
            })}
        </div>
    );
}
