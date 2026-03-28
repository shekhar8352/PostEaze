import { useEffect, useMemo, useState } from "react";
import { Container, Stack, SimpleGrid, Alert, Button, Center, Loader } from "@mantine/core";
import { useNavigate } from "react-router-dom";
import { useChannels } from "@/features/channels";
import type { InstagramChannelDisplay } from "@/features/channels/types/instagram.types";
import { Icons } from "@/app/theme";
import {
  useAnalyticsDashboard,
  useAnalyticsComparison,
  useProfileTimeSeries,
  useStoriesAnalytics,
  useTopPosts,
} from "../../services/analyticsQueries";
import { AnalyticsHeader } from "../../components/AnalyticsHeader";
import { inclusiveRange } from "../../utils/dateRange";
import { KPICard } from "../../components/KPICard";
import { FollowerGrowthChart } from "../../components/FollowerGrowthChart";
import { ProfileTrendChart } from "../../components/ProfileTrendChart";
import { PostEngagementChart } from "../../components/PostEngagementChart";
import { PostsActivitySummary } from "../../components/PostsActivitySummary";
import { TopPostsTable } from "../../components/TopPostsTable";
import { StoriesPanel } from "../../components/StoriesPanel";
import { AudienceInsightsPanel } from "../../components/AudienceInsightsPanel";
import { computePctDelta, formatCompact } from "../../utils/pctChange";
import type { TopPostsSort } from "../../types";
import styles from "./AnalyticsPage.module.css";

export default function AnalyticsPage() {
  const navigate = useNavigate();
  const { data: channels, isLoading: channelsLoading } = useChannels();
  const [channelId, setChannelId] = useState<number | null>(null);
  const [range, setRange] = useState(() => inclusiveRange(7));
  const [topSort, setTopSort] = useState<TopPostsSort>("engagement");
  const [postType, setPostType] = useState("");

  const instagramChannels = useMemo(
    () => (channels ?? []).filter((c) => c.provider === "instagram") as InstagramChannelDisplay[],
    [channels]
  );

  useEffect(() => {
    if (channelId == null && instagramChannels.length > 0) {
      setChannelId(instagramChannels[0].channel_id);
    }
  }, [instagramChannels, channelId]);

  const rangeParams = useMemo(
    () => ({ startDate: range.startDate, endDate: range.endDate }),
    [range.startDate, range.endDate]
  );

  const dashboard = useAnalyticsDashboard(channelId, rangeParams);
  const comparison = useAnalyticsComparison(channelId, rangeParams);
  const profile = useProfileTimeSeries(channelId, rangeParams);
  const stories = useStoriesAnalytics(channelId, rangeParams);
  const topPosts = useTopPosts(channelId, rangeParams, topSort, postType || undefined);

  const prev = comparison.data?.previous;

  const kpiLoading = dashboard.isLoading || comparison.isLoading;

  const overview = dashboard.data?.overview;

  if (channelsLoading) {
    return (
      <Center className={styles.loader}>
        <Loader size="md" type="bars" />
      </Center>
    );
  }

  if (instagramChannels.length === 0) {
    return (
      <Container size="md" className={styles.container}>
        <Alert
          variant="light"
          color="blue"
          title="Connect Instagram"
          icon={<Icons.Instagram size={20} />}
        >
          Analytics appear after you connect an Instagram channel.
          <Button mt="md" onClick={() => navigate("/channels/instagram")}>
            Go to Instagram channels
          </Button>
        </Alert>
      </Container>
    );
  }

  return (
    <Container size="xl" className={styles.container}>
      <Stack gap="xl">
        <AnalyticsHeader
          channels={instagramChannels}
          channelId={channelId}
          onChannelChange={setChannelId}
          range={range}
          onRangeChange={setRange}
        />

        <SimpleGrid cols={{ base: 1, sm: 2, lg: 3 }} spacing="md">
          <KPICard
            title="Followers (end of range)"
            value={overview ? overview.end_followers.toLocaleString() : "—"}
            deltaPct={computePctDelta(overview?.end_followers ?? 0, prev?.end_followers)}
            loading={kpiLoading}
          />
          <KPICard
            title="Total reach"
            value={overview ? formatCompact(overview.total_reach) : "—"}
            deltaPct={computePctDelta(overview?.total_reach ?? 0, prev?.total_reach)}
            loading={kpiLoading}
          />
          <KPICard
            title="Total impressions"
            value={overview ? formatCompact(overview.total_impressions) : "—"}
            deltaPct={computePctDelta(overview?.total_impressions ?? 0, prev?.total_impressions)}
            loading={kpiLoading}
          />
          <KPICard
            title="Profile views"
            value={overview ? formatCompact(overview.total_profile_views) : "—"}
            deltaPct={computePctDelta(overview?.total_profile_views ?? 0, prev?.total_profile_views)}
            loading={kpiLoading}
          />
          <KPICard
            title="Accounts engaged"
            value={overview ? formatCompact(overview.total_accounts_engaged) : "—"}
            deltaPct={computePctDelta(overview?.total_accounts_engaged ?? 0, prev?.total_accounts_engaged)}
            loading={kpiLoading}
          />
          <KPICard
            title="Total interactions"
            value={overview ? formatCompact(overview.total_interactions) : "—"}
            deltaPct={computePctDelta(overview?.total_interactions ?? 0, prev?.total_interactions)}
            loading={kpiLoading}
          />
          <KPICard
            title="Content views"
            value={overview ? formatCompact(overview.total_views) : "—"}
            deltaPct={computePctDelta(overview?.total_views ?? 0, prev?.total_views)}
            loading={kpiLoading}
          />
          <KPICard
            title="Website clicks"
            value={overview ? formatCompact(overview.total_website_clicks) : "—"}
            deltaPct={computePctDelta(overview?.total_website_clicks ?? 0, prev?.total_website_clicks)}
            loading={kpiLoading}
          />
        </SimpleGrid>

        <AudienceInsightsPanel audience={dashboard.data?.audience} loading={dashboard.isLoading} />

        <SimpleGrid cols={{ base: 1, lg: 2 }} spacing="md">
          <FollowerGrowthChart series={profile.data?.analytics ?? []} loading={profile.isLoading} />
          <PostsActivitySummary overview={dashboard.data?.posts_overview} loading={dashboard.isLoading} />
        </SimpleGrid>

        <ProfileTrendChart series={profile.data?.analytics ?? []} loading={profile.isLoading} />

        <PostEngagementChart overview={dashboard.data?.posts_overview} loading={dashboard.isLoading} />

        <TopPostsTable
          posts={topPosts.data?.top_posts ?? []}
          loading={topPosts.isLoading}
          apiSort={topSort}
          onApiSortChange={setTopSort}
          postType={postType}
          onPostTypeChange={setPostType}
        />

        <StoriesPanel stories={stories.data?.stories ?? []} loading={stories.isLoading} />
      </Stack>
    </Container>
  );
}
