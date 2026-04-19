import {
    Badge,
    Button,
    Divider,
    Group,
    Paper,
    SimpleGrid,
    Skeleton,
    Stack,
    Text,
    Title,
} from '@mantine/core';
import { Link } from 'react-router-dom';
import { Icons } from '@/app/theme';
import type { DateRangeParams } from '@/features/analytics/services/analyticsService';
import type { AggregatedProfileOverview, DashboardResponse, TopPostItem } from '@/features/analytics/types';
import type { BaseChannelDisplay } from '@/features/channels/types/base.types';
import { formatDashboardInt } from '../utils/formatDashboardNumbers';
import dashStyles from '../DashboardPage.module.css';

interface DashboardInstagramPulseProps {
    analyticsRange: DateRangeParams;
    instagramChannel: BaseChannelDisplay | undefined;
    analyticsDashboard: DashboardResponse | undefined;
    analyticsLoading: boolean;
}

export function DashboardInstagramPulse({
    analyticsRange,
    instagramChannel,
    analyticsDashboard,
    analyticsLoading,
}: DashboardInstagramPulseProps) {
    const overview: AggregatedProfileOverview | undefined = analyticsDashboard?.overview;
    const topPosts: TopPostItem[] = (analyticsDashboard?.top_posts ?? []).slice(0, 3);

    return (
        <Paper
            p="xl"
            radius="lg"
            withBorder={false}
            className={`${dashStyles.section} ${dashStyles.peakStripWrap}`}
        >
            <div className={dashStyles.peakStrip} aria-hidden />
            <Stack gap="md" style={{ position: 'relative' }}>
                <Group justify="space-between" align="flex-start" wrap="wrap">
                    <div>
                        <Title order={2} className={dashStyles.sectionTitle} mb={6}>
                            Instagram pulse
                        </Title>
                        <Text className={dashStyles.sectionDesc}>
                            Last 7 days on your primary Instagram account—reach, engagement, and follower movement.
                            Connect Instagram to unlock this panel.
                        </Text>
                    </div>
                    <Badge variant="light" color="violet">
                        {analyticsRange.startDate} → {analyticsRange.endDate}
                    </Badge>
                </Group>

                {!instagramChannel ? (
                    <Text size="sm" c="dimmed">
                        Connect an Instagram channel to pull Meta insights into the dashboard.
                    </Text>
                ) : analyticsLoading ? (
                    <SimpleGrid cols={2} spacing="sm">
                        {Array.from({ length: 4 }).map((_, i) => (
                            <Skeleton key={i} height={64} />
                        ))}
                    </SimpleGrid>
                ) : !overview ? (
                    <Text size="sm" c="dimmed">
                        No overview data for this range yet. Run a sync or pick another week in Analytics.
                    </Text>
                ) : (
                    <>
                        <div className={dashStyles.analyticsGrid}>
                            <div className={dashStyles.analyticsMetric}>
                                <div className={dashStyles.analyticsMetricLabel}>Reach</div>
                                <div className={dashStyles.analyticsMetricValue}>
                                    {formatDashboardInt(overview.total_reach)}
                                </div>
                            </div>
                            <div className={dashStyles.analyticsMetric}>
                                <div className={dashStyles.analyticsMetricLabel}>Impressions</div>
                                <div className={dashStyles.analyticsMetricValue}>
                                    {formatDashboardInt(overview.total_impressions)}
                                </div>
                            </div>
                            <div className={dashStyles.analyticsMetric}>
                                <div className={dashStyles.analyticsMetricLabel}>Accounts engaged</div>
                                <div className={dashStyles.analyticsMetricValue}>
                                    {formatDashboardInt(overview.total_accounts_engaged)}
                                </div>
                            </div>
                            <div className={dashStyles.analyticsMetric}>
                                <div className={dashStyles.analyticsMetricLabel}>Interactions</div>
                                <div className={dashStyles.analyticsMetricValue}>
                                    {formatDashboardInt(overview.total_interactions)}
                                </div>
                            </div>
                        </div>
                        <Group gap="xl" wrap="wrap">
                            <div>
                                <Text size="xs" c="dimmed" tt="uppercase" fw={600}>
                                    Follower growth
                                </Text>
                                <Text size="lg" fw={700}>
                                    {overview.follower_growth >= 0 ? '+' : ''}
                                    {overview.follower_growth}
                                </Text>
                            </div>
                            <div>
                                <Text size="xs" c="dimmed" tt="uppercase" fw={600}>
                                    Net new posts (range)
                                </Text>
                                <Text size="lg" fw={700}>
                                    {analyticsDashboard?.posts_overview?.new_posts ?? '—'}
                                </Text>
                            </div>
                        </Group>
                    </>
                )}

                {instagramChannel && topPosts.length > 0 ? (
                    <>
                        <Divider label="Top posts" labelPosition="left" />
                        <Stack gap="xs">
                            {topPosts.map((tp) => (
                                <div key={tp.post_id} className={dashStyles.topPostCard}>
                                    <Group justify="space-between" align="flex-start" wrap="nowrap" gap="sm">
                                        <Text size="sm" fw={500} lineClamp={2} style={{ flex: 1 }}>
                                            {tp.caption?.trim() || tp.post_type}
                                        </Text>
                                        <Badge variant="light" color="grape" size="sm">
                                            {formatDashboardInt(tp.engagement)} eng.
                                        </Badge>
                                    </Group>
                                </div>
                            ))}
                        </Stack>
                    </>
                ) : null}

                <Button component={Link} to="/analytics" variant="light" leftSection={<Icons.ChartBar size={18} />}>
                    Explore analytics
                </Button>
            </Stack>
        </Paper>
    );
}
