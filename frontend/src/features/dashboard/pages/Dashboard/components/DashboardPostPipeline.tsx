import {
    Badge,
    Box,
    Button,
    Group,
    Paper,
    Skeleton,
    Stack,
    Text,
    Title,
    useMantineTheme,
} from '@mantine/core';
import { format, parseISO } from 'date-fns';
import { Icons } from '@/app/theme';
import { formatScheduledPostStatus } from '@/features/calendar/utils/postStatus';
import type { DashboardPostMetrics } from '../types';
import dashStyles from '../DashboardPage.module.css';

interface DashboardPostPipelineProps {
    pipelineFrom: string;
    pipelineTo: string;
    postMetrics: DashboardPostMetrics;
    postsLoading: boolean;
    postsError: boolean;
    onNavigateCalendar: () => void;
}

function segmentColor(theme: ReturnType<typeof useMantineTheme>, key: string, fallback: string): string {
    const map: Record<string, [string, number]> = {
        published: ['green', 6],
        scheduled: ['blue', 6],
        failed: ['red', 6],
        cancelled: ['gray', 6],
        other: ['gray', 4],
    };
    const entry = map[key];
    if (!entry) return fallback;
    const [name, shade] = entry;
    const palette = theme.colors[name as keyof typeof theme.colors];
    const hex = Array.isArray(palette) ? palette[shade] : undefined;
    return typeof hex === 'string' ? hex : fallback;
}

export function DashboardPostPipeline({
    pipelineFrom,
    pipelineTo,
    postMetrics,
    postsLoading,
    postsError,
    onNavigateCalendar,
}: DashboardPostPipelineProps) {
    const theme = useMantineTheme();
    const pipelineTotal = postMetrics.total || 1;
    const segPublished = (postMetrics.counts.published / pipelineTotal) * 100;
    const segScheduled = (postMetrics.counts.scheduled / pipelineTotal) * 100;
    const segFailed = (postMetrics.counts.failed / pipelineTotal) * 100;
    const segCancelled = (postMetrics.counts.cancelled / pipelineTotal) * 100;
    const segOther = (postMetrics.counts.other / pipelineTotal) * 100;

    const legendItems = [
        { key: 'published', label: 'Published', count: postMetrics.counts.published },
        { key: 'scheduled', label: 'Scheduled', count: postMetrics.counts.scheduled },
        { key: 'failed', label: 'Failed', count: postMetrics.counts.failed },
        { key: 'cancelled', label: 'Cancelled', count: postMetrics.counts.cancelled },
        { key: 'other', label: 'Other', count: postMetrics.counts.other },
    ];

    const pipelineViz = (
        <Stack gap="md" w="100%" miw={0}>
            <Title order={3} className={dashStyles.subsectionHeading}>
                Status mix
            </Title>
            <Box w="100%">
                <div
                    className={dashStyles.pipelineBar}
                    role="img"
                    aria-label="Post status distribution"
                >
                    <div
                        className={dashStyles.pipelineSegment}
                        style={{
                            width: `${segPublished}%`,
                            background: segmentColor(theme, 'published', '#40c057'),
                        }}
                    />
                    <div
                        className={dashStyles.pipelineSegment}
                        style={{
                            width: `${segScheduled}%`,
                            background: segmentColor(theme, 'scheduled', '#228be6'),
                        }}
                    />
                    <div
                        className={dashStyles.pipelineSegment}
                        style={{
                            width: `${segFailed}%`,
                            background: segmentColor(theme, 'failed', '#fa5252'),
                        }}
                    />
                    <div
                        className={dashStyles.pipelineSegment}
                        style={{
                            width: `${segCancelled}%`,
                            background: segmentColor(theme, 'cancelled', '#868e96'),
                        }}
                    />
                    <div
                        className={dashStyles.pipelineSegment}
                        style={{
                            width: `${segOther}%`,
                            background: segmentColor(theme, 'other', '#adb5bd'),
                        }}
                    />
                </div>
            </Box>
            <div className={dashStyles.pipelineLegend}>
                {legendItems.map((item) => (
                    <div key={item.key} className={dashStyles.legendItem}>
                        <span
                            className={dashStyles.legendSwatch}
                            style={{ background: segmentColor(theme, item.key, '#adb5bd') }}
                        />
                        <Text size="sm" c="dimmed" span>
                            {item.label}{' '}
                            <Text span fw={600} c="var(--pe-text)">
                                {item.count}
                            </Text>
                        </Text>
                    </div>
                ))}
            </div>
            {postMetrics.healthPct !== null ? (
                <Text size="sm" c="dimmed">
                    Delivery health (published vs failed in window):{' '}
                    <Text span fw={600} c="var(--pe-text)">
                        {postMetrics.healthPct}%
                    </Text>
                </Text>
            ) : null}
        </Stack>
    );

    const nextUpBlock = (
        <Stack gap="md" w="100%" miw={0}>
            <Title order={3} className={dashStyles.subsectionHeading}>
                Next up
            </Title>
            {postsLoading ? (
                <Stack gap="xs">
                    <Skeleton height={40} radius="md" />
                    <Skeleton height={40} radius="md" />
                </Stack>
            ) : postMetrics.upcomingRows.length === 0 ? (
                <Text size="sm" c="dimmed" py="xs">
                    Nothing scheduled in the next 30 days.
                </Text>
            ) : (
                <div className={dashStyles.nextUpList}>
                    <Stack gap={0}>
                        {postMetrics.upcomingRows.map((p) => (
                            <div key={p.id} className={dashStyles.miniPostRow}>
                                <div style={{ minWidth: 0 }}>
                                    <Text size="sm" fw={600} lineClamp={1} c="var(--pe-text)">
                                        {p.caption?.trim() || p.post_type}
                                    </Text>
                                    <Text size="xs" c="dimmed">
                                        {format(parseISO(p.scheduled_at), 'MMM d, h:mm a')} ·{' '}
                                        {p.platforms.join(', ')}
                                    </Text>
                                </div>
                                <Badge size="sm" variant="light" color="blue" style={{ flexShrink: 0 }}>
                                    {formatScheduledPostStatus(p.status)}
                                </Badge>
                            </div>
                        ))}
                    </Stack>
                </div>
            )}
            <Button variant="light" rightSection={<Icons.ArrowRight size={16} />} onClick={onNavigateCalendar}>
                Manage in calendar
            </Button>
        </Stack>
    );

    const loadingBody = (
        <Stack gap="xl" w="100%" className={dashStyles.pipelineBody}>
            <Stack gap="md" w="100%">
                <Title order={3} className={dashStyles.subsectionHeading}>
                    Status mix
                </Title>
                <Skeleton height={14} radius="xl" w="100%" />
                <Skeleton height={56} radius="md" w="100%" />
                <Skeleton height={16} width="70%" />
            </Stack>
            <Stack gap="md" w="100%">
                <Title order={3} className={dashStyles.subsectionHeading}>
                    Next up
                </Title>
                <Skeleton height={40} radius="md" />
                <Skeleton height={40} radius="md" />
                <Button variant="light" rightSection={<Icons.ArrowRight size={16} />} onClick={onNavigateCalendar}>
                    Manage in calendar
                </Button>
            </Stack>
        </Stack>
    );

    const mainColumn =
        !postsError && !postsLoading && postMetrics.total > 0 ? (
            <Stack gap="xl" w="100%" className={dashStyles.pipelineBody}>
                {pipelineViz}
                {nextUpBlock}
            </Stack>
        ) : null;

    return (
        <Paper p="xl" radius="lg" withBorder={false} className={dashStyles.section}>
            <Stack gap="lg" align="stretch">
                <Group justify="space-between" align="flex-start" wrap="wrap" gap="md">
                    <div style={{ minWidth: 0, flex: '1 1 16rem' }}>
                        <Title order={2} className={dashStyles.sectionTitle} mb={6}>
                            Post pipeline
                        </Title>
                        <Text className={dashStyles.sectionDesc}>
                            Status mix for posts scheduled between {pipelineFrom} and {pipelineTo} (exclusive end).
                            Tracks what is still queued versus what shipped.
                        </Text>
                    </div>
                    <Badge variant="light" color="gray" style={{ flexShrink: 0 }}>
                        {postsLoading ? 'Loading…' : `${postMetrics.total} posts`}
                    </Badge>
                </Group>

                {postsError ? (
                    <Stack gap="xl" w="100%" className={dashStyles.pipelineBody}>
                        <Text size="sm" c="red">
                            Could not load scheduled posts. Try again from the calendar view.
                        </Text>
                        <Stack gap="md" w="100%">
                            <Title order={3} className={dashStyles.subsectionHeading}>
                                Next up
                            </Title>
                            <Text size="sm" c="dimmed">
                                Open the calendar to retry or create a new schedule.
                            </Text>
                            <Button
                                variant="light"
                                rightSection={<Icons.ArrowRight size={16} />}
                                onClick={onNavigateCalendar}
                            >
                                Manage in calendar
                            </Button>
                        </Stack>
                    </Stack>
                ) : postsLoading ? (
                    loadingBody
                ) : postMetrics.total === 0 ? (
                    <Stack gap="xl" w="100%" className={dashStyles.pipelineBody}>
                        <Text size="sm" c="dimmed">
                            No posts in this window yet. Schedule from the calendar to see your pipeline here.
                        </Text>
                        <Stack gap="md" w="100%">
                            <Title order={3} className={dashStyles.subsectionHeading}>
                                Next up
                            </Title>
                            <Text size="sm" c="dimmed">
                                Nothing scheduled in the next 30 days.
                            </Text>
                            <Button
                                variant="light"
                                rightSection={<Icons.ArrowRight size={16} />}
                                onClick={onNavigateCalendar}
                            >
                                Manage in calendar
                            </Button>
                        </Stack>
                    </Stack>
                ) : (
                    mainColumn
                )}
            </Stack>
        </Paper>
    );
}
