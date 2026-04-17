import { useEffect, useMemo, useState } from "react";
import {
  Container,
  Stack,
  SimpleGrid,
  Alert,
  Button,
  Group,
  Center,
  Loader,
  Box,
  LoadingOverlay,
} from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { useChannels, channelKeys } from "@/features/channels";
import { syncMetaAnalytics } from "@/services/meta/metaAnalyticsService";
import type { BaseChannelDisplay } from "@/features/channels/types/base.types";
import { Icons } from "@/app/theme";
import {
  analyticsKeys,
  useAnalyticsDashboard,
  useAnalyticsComparison,
  useProfileTimeSeries,
  useStoriesAnalytics,
  useTopPosts,
  useDailyEngagement,
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
import { DailyEngagementSection } from "../../components/DailyEngagementSection";
import { computePctDelta, formatCompact } from "../../utils/pctChange";
import type { TopPostsSort } from "../../types";
import styles from "./AnalyticsPage.module.css";

const DELTA_HINT = "No prior period";

export default function AnalyticsPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { data: channels, isLoading: channelsLoading } = useChannels();
  const [syncingMeta, setSyncingMeta] = useState(false);
  const [platform, setPlatform] = useState<"instagram" | "facebook">("instagram");
  const [channelId, setChannelId] = useState<number | null>(null);
  const [range, setRange] = useState(() => inclusiveRange(7));
  const [topSort, setTopSort] = useState<TopPostsSort>("engagement");
  const [postType, setPostType] = useState("");

  const hasInstagram = useMemo(
    () => (channels ?? []).some((c) => c.provider === "instagram"),
    [channels]
  );
  const hasFacebook = useMemo(
    () => (channels ?? []).some((c) => c.provider === "facebook"),
    [channels]
  );

  /** When only one Meta provider exists, lock to it; otherwise use the user’s tab selection. */
  const resolvedPlatform = useMemo<"instagram" | "facebook">(() => {
    if (hasInstagram && hasFacebook) return platform;
    if (hasInstagram) return "instagram";
    if (hasFacebook) return "facebook";
    return platform;
  }, [hasInstagram, hasFacebook, platform]);

  const providerChannels = useMemo(
    () =>
      (channels ?? []).filter((c) => c.provider === resolvedPlatform) as BaseChannelDisplay[],
    [channels, resolvedPlatform]
  );

  /** All Meta-connected channels (API accepts optional filter; we sync every IG + FB account). */
  const metaChannelIds = useMemo(
    () =>
      (channels ?? [])
        .filter((c) => c.provider === "instagram" || c.provider === "facebook")
        .map((c) => c.channel_id),
    [channels]
  );

  const handleSyncMeta = async () => {
    if (metaChannelIds.length === 0 || syncingMeta) return;
    setSyncingMeta(true);
    try {
      const res = await syncMetaAnalytics(metaChannelIds);
      const ok = res.results.filter((r) => r.status === "ok").length;
      const failed = res.results.filter((r) => r.status !== "ok");
      notifications.show({
        title: "Meta sync complete",
        message:
          failed.length === 0
            ? `Updated ${ok} channel(s). Refreshing analytics…`
            : `${ok} succeeded, ${failed.length} had issues. Refreshing…`,
        color: failed.length ? "yellow" : "teal",
      });
      await queryClient.invalidateQueries({ queryKey: analyticsKeys.all });
      await queryClient.invalidateQueries({ queryKey: channelKeys.lists() });
    } catch (e) {
      const msg = e instanceof Error ? e.message : "Sync failed";
      notifications.show({ title: "Meta sync failed", message: msg, color: "red" });
    } finally {
      setSyncingMeta(false);
    }
  };

  useEffect(() => {
    if (providerChannels.length === 0) {
      setChannelId(null);
      return;
    }
    const valid = providerChannels.some((c) => c.channel_id === channelId);
    if (!valid) setChannelId(providerChannels[0].channel_id);
  }, [providerChannels, channelId]);

  const rangeParams = useMemo(
    () => ({ startDate: range.startDate, endDate: range.endDate }),
    [range.startDate, range.endDate]
  );

  const dashboard = useAnalyticsDashboard(channelId, rangeParams);
  const comparison = useAnalyticsComparison(channelId, rangeParams);
  const profile = useProfileTimeSeries(channelId, rangeParams);
  const stories = useStoriesAnalytics(channelId, rangeParams);
  const topPosts = useTopPosts(channelId, rangeParams, topSort, postType || undefined);
  const dailyEngagement = useDailyEngagement(
    resolvedPlatform === "instagram" ? channelId : null,
    resolvedPlatform === "instagram" ? rangeParams : null
  );

  const prev = comparison.data?.previous;

  const kpiLoading =
    dashboard.isLoading ||
    dashboard.isFetching ||
    comparison.isLoading ||
    comparison.isFetching ||
    syncingMeta;

  const overview = dashboard.data?.overview;

  if (channelsLoading) {
    return (
      <Center className={styles.loader}>
        <Loader size="md" type="bars" />
      </Center>
    );
  }

  if (!hasInstagram && !hasFacebook) {
    return (
      <Container size="md" className={styles.container}>
        <Alert variant="light" color="blue" title="Connect a channel" icon={<Icons.ChartBar size={20} />}>
          Analytics appear after you connect Instagram or Facebook Page.
          <Group mt="md" gap="sm">
            <Button leftSection={<Icons.Instagram size={18} />} onClick={() => navigate("/channels/instagram")}>
              Instagram
            </Button>
            <Button leftSection={<Icons.Facebook size={18} />} onClick={() => navigate("/channels/facebook")}>
              Facebook
            </Button>
          </Group>
        </Alert>
      </Container>
    );
  }

  if (providerChannels.length === 0) {
    return (
      <Container size="md" className={styles.container}>
        <Alert
          variant="light"
          color="gray"
          title={`No ${resolvedPlatform === "instagram" ? "Instagram" : "Facebook"} accounts`}
          icon={
            resolvedPlatform === "instagram" ? (
              <Icons.Instagram size={20} />
            ) : (
              <Icons.Facebook size={20} />
            )
          }
        >
          Connect an account to see analytics here, or switch the platform above.
          <Button
            mt="md"
            onClick={() =>
              navigate(resolvedPlatform === "instagram" ? "/channels/instagram" : "/channels/facebook")
            }
          >
            Go to {resolvedPlatform === "instagram" ? "Instagram" : "Facebook"} channels
          </Button>
        </Alert>
      </Container>
    );
  }

  return (
    <Container size="xl" className={styles.container}>
      <Stack gap="xl">
        <AnalyticsHeader
          platform={resolvedPlatform}
          onPlatformChange={setPlatform}
          hasInstagram={hasInstagram}
          hasFacebook={hasFacebook}
          channels={providerChannels}
          channelId={channelId}
          onChannelChange={setChannelId}
          range={range}
          onRangeChange={setRange}
          onSyncMeta={handleSyncMeta}
          syncing={syncingMeta}
          syncDisabled={metaChannelIds.length === 0}
        />

        <Box pos="relative" mih={320}>
          <LoadingOverlay
            visible={syncingMeta}
            overlayProps={{ blur: 2 }}
            loaderProps={{ type: "bars", size: "lg" }}
            zIndex={5}
          />
          <Stack gap="xl">
          <SimpleGrid cols={{ base: 1, sm: 2, lg: 4 }} spacing="md">
          <KPICard
            title="Followers (end of range)"
            value={overview ? overview.end_followers.toLocaleString() : "—"}
            deltaPct={computePctDelta(overview?.end_followers ?? 0, prev?.end_followers)}
            loading={kpiLoading}
            deltaHint={DELTA_HINT}
          />
          <KPICard
            title="Total reach"
            value={overview ? formatCompact(overview.total_reach) : "—"}
            deltaPct={computePctDelta(overview?.total_reach ?? 0, prev?.total_reach)}
            loading={kpiLoading}
            deltaHint={DELTA_HINT}
          />
          <KPICard
            title="Total impressions"
            value={overview ? formatCompact(overview.total_impressions) : "—"}
            deltaPct={computePctDelta(overview?.total_impressions ?? 0, prev?.total_impressions)}
            loading={kpiLoading}
            deltaHint={DELTA_HINT}
          />
          <KPICard
            title="Profile views"
            value={overview ? formatCompact(overview.total_profile_views) : "—"}
            deltaPct={computePctDelta(overview?.total_profile_views ?? 0, prev?.total_profile_views)}
            loading={kpiLoading}
            deltaHint={DELTA_HINT}
          />
          <KPICard
            title="Accounts engaged"
            value={overview ? formatCompact(overview.total_accounts_engaged) : "—"}
            deltaPct={computePctDelta(overview?.total_accounts_engaged ?? 0, prev?.total_accounts_engaged)}
            loading={kpiLoading}
            deltaHint={DELTA_HINT}
          />
          <KPICard
            title="Total interactions"
            value={overview ? formatCompact(overview.total_interactions) : "—"}
            deltaPct={computePctDelta(overview?.total_interactions ?? 0, prev?.total_interactions)}
            loading={kpiLoading}
            deltaHint={DELTA_HINT}
          />
          <KPICard
            title="Content views"
            value={overview ? formatCompact(overview.total_views) : "—"}
            deltaPct={computePctDelta(overview?.total_views ?? 0, prev?.total_views)}
            loading={kpiLoading}
            deltaHint={DELTA_HINT}
          />
          <KPICard
            title="Website clicks"
            value={overview ? formatCompact(overview.total_website_clicks) : "—"}
            deltaPct={computePctDelta(overview?.total_website_clicks ?? 0, prev?.total_website_clicks)}
            loading={kpiLoading}
            deltaHint={DELTA_HINT}
          />
          </SimpleGrid>

        <AudienceInsightsPanel audience={dashboard.data?.audience} loading={dashboard.isLoading} />

        <SimpleGrid cols={{ base: 1, lg: 2 }} spacing="md">
          <FollowerGrowthChart series={profile.data?.analytics ?? []} loading={profile.isLoading} />
          <PostsActivitySummary overview={dashboard.data?.posts_overview} loading={dashboard.isLoading} />
        </SimpleGrid>

        <ProfileTrendChart series={profile.data?.analytics ?? []} loading={profile.isLoading} />

        {resolvedPlatform === "instagram" && (
          <DailyEngagementSection
            startDate={range.startDate}
            endDate={range.endDate}
            series={dailyEngagement.data?.series ?? []}
            overview={dashboard.data?.posts_overview}
            note={dailyEngagement.data?.note}
            loading={dailyEngagement.isLoading}
          />
        )}

        <PostEngagementChart overview={dashboard.data?.posts_overview} loading={dashboard.isLoading} />

        <TopPostsTable
          posts={topPosts.data?.top_posts ?? []}
          loading={topPosts.isLoading}
          apiSort={topSort}
          onApiSortChange={setTopSort}
          postType={postType}
          onPostTypeChange={setPostType}
        />

        {resolvedPlatform === "instagram" && (
          <StoriesPanel stories={stories.data?.stories ?? []} loading={stories.isLoading} />
        )}
          </Stack>
        </Box>
      </Stack>
    </Container>
  );
}
