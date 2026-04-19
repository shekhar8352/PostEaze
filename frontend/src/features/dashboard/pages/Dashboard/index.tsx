import { useMemo } from 'react';
import { Container, Stack, Loader, Center } from '@mantine/core';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/features/auth';
import { Icons } from '@/app/theme';
import { useChannels } from '@/features/channels';
import { useScheduledPostsRange } from '@/features/calendar/hooks/useScheduledPostsQueries';
import { useAnalyticsDashboard } from '@/features/analytics/services/analyticsQueries';
import {
    DashboardChannelsSection,
    DashboardHero,
    DashboardInstagramPulse,
    DashboardKpiGrid,
    DashboardPostPipeline,
} from './components';
import { useDashboardDateRanges } from './hooks/useDashboardDateRanges';
import type { DashboardKpiStat } from './types';
import { buildDashboardChannelCards } from './utils/buildDashboardChannels';
import { computeDashboardPostMetrics } from './utils/computePostMetrics';
import dashStyles from './DashboardPage.module.css';

const DashboardPage = () => {
    const { user } = useAuth();
    const navigate = useNavigate();
    const { data: connectedChannels, isLoading: channelsLoading } = useChannels();
    const { pipelineFrom, pipelineTo, analyticsRange } = useDashboardDateRanges();

    const { data: scheduledPayload, isLoading: postsLoading, isError: postsError } =
        useScheduledPostsRange(pipelineFrom, pipelineTo);

    const instagramChannel = useMemo(
        () => connectedChannels?.find((ch) => ch.provider === 'instagram' && ch.isConnected),
        [connectedChannels]
    );

    const { data: analyticsDashboard, isLoading: analyticsLoading } = useAnalyticsDashboard(
        instagramChannel?.channel_id ?? null,
        instagramChannel ? analyticsRange : null
    );

    const displayName = user?.name || user?.email?.split('@')[0] || 'there';

    const channelCards = useMemo(() => buildDashboardChannelCards(connectedChannels), [connectedChannels]);

    const postMetrics = useMemo(
        () => computeDashboardPostMetrics(scheduledPayload?.posts ?? []),
        [scheduledPayload?.posts]
    );

    const publishedLast7d =
        analyticsDashboard?.posts_overview?.new_posts ?? postMetrics.publishedLast7dLocal;

    const stats: DashboardKpiStat[] = [
        {
            title: 'Connected accounts',
            value: connectedChannels?.length ?? 0,
            icon: Icons.User,
            loading: false,
            attention: false,
        },
        {
            title: 'Upcoming (30 days)',
            value: postMetrics.upcomingCount,
            icon: Icons.Calendar,
            loading: postsLoading,
            attention: false,
        },
        {
            title: 'Published (last 7 days)',
            value: publishedLast7d,
            icon: Icons.CheckCircle,
            loading: instagramChannel ? analyticsLoading : postsLoading,
            attention: false,
        },
        {
            title: 'Needs attention',
            subtitle: 'Failed or partial in pipeline window',
            value: postMetrics.counts.failed,
            icon: Icons.AlertTriangle,
            loading: postsLoading,
            attention: postMetrics.counts.failed > 0,
        },
    ];

    if (channelsLoading) {
        return (
            <Center style={{ minHeight: '60vh' }}>
                <Loader size="md" color="blue" type="bars" />
            </Center>
        );
    }

    return (
        <Container size="xl" className="fade-in">
            <Stack gap="xl">
                <DashboardHero displayName={displayName} />
                <DashboardKpiGrid stats={stats} />
                <div className={dashStyles.twoCol}>
                    <DashboardPostPipeline
                        pipelineFrom={pipelineFrom}
                        pipelineTo={pipelineTo}
                        postMetrics={postMetrics}
                        postsLoading={postsLoading}
                        postsError={postsError}
                        onNavigateCalendar={() => navigate('/calendar')}
                    />
                    <DashboardInstagramPulse
                        analyticsRange={analyticsRange}
                        instagramChannel={instagramChannel}
                        analyticsDashboard={analyticsDashboard}
                        analyticsLoading={analyticsLoading}
                    />
                </div>
                <DashboardChannelsSection
                    channels={channelCards}
                    onNavigateChannel={(path) => navigate(path)}
                />
            </Stack>
        </Container>
    );
};

export default DashboardPage;
