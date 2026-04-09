import { useMemo, useState } from 'react';
import {
  Container,
  Title,
  Paper,
  Stack,
  Group,
  Badge,
  Button,
  Text,
  Modal,
  Radio,
  Divider,
  Box,
  Loader,
} from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { Icons } from '@/app/theme';
import { ChannelPageHeader } from '../components/ChannelPageHeader';
import { EmptyChannelPanel } from '../components/EmptyChannelPanel';
import { useFacebookOAuth } from '../hooks/useFacebookOAuth';
import { useChannels, channelKeys } from '../services/channelQueries';
import { useQueryClient } from '@tanstack/react-query';
import {
  fetchMetaPagesFromCode,
  createFacebookChannelFromPageToken,
  type MetaPageRow,
} from '@/services/meta/metaService';
import { syncMetaAnalytics } from '@/services/meta/metaAnalyticsService';
import { getFacebookRedirectUri } from '../utils/facebookOAuth.utils';
import type { BaseChannelDisplay } from '../types/base.types';
import styles from './FacebookChannel.module.css';

const FacebookChannel = () => {
  const queryClient = useQueryClient();
  const { data: channels, isLoading: channelsLoading } = useChannels();
  const { openOAuthPopup, isLoading: oauthLoading } = useFacebookOAuth();

  const [pickerOpen, setPickerOpen] = useState(false);
  const [pages, setPages] = useState<MetaPageRow[]>([]);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [connecting, setConnecting] = useState(false);
  const [syncing, setSyncing] = useState(false);

  const facebookChannels = useMemo(
    () => (channels ?? []).filter((c) => c.provider === 'facebook') as BaseChannelDisplay[],
    [channels]
  );

  const metaChannelIds = useMemo(() => {
    const ids = (channels ?? [])
      .filter((c) => c.provider === 'instagram' || c.provider === 'facebook')
      .map((c) => c.channel_id);
    return ids;
  }, [channels]);

  const handleConnectClick = async () => {
    try {
      const code = await openOAuthPopup();
      const redirectUri = getFacebookRedirectUri();
      setConnecting(true);
      const list = await fetchMetaPagesFromCode(code, redirectUri);
      if (!list.length) {
        notifications.show({
          title: 'No pages found',
          message: 'This account has no Facebook Pages, or required permissions were not granted.',
          color: 'yellow',
        });
        return;
      }
      setPages(list);
      setSelectedId(list[0]?.id ?? null);
      setPickerOpen(true);
    } catch (e) {
      const msg = e instanceof Error ? e.message : 'Facebook login failed';
      if (!msg.includes('cancelled')) {
        notifications.show({ title: 'Connect failed', message: msg, color: 'red' });
      }
    } finally {
      setConnecting(false);
    }
  };

  const handleConfirmPage = async () => {
    const row = pages.find((p) => p.id === selectedId);
    if (!row) return;
    try {
      setConnecting(true);
      await createFacebookChannelFromPageToken({
        page_id: row.id,
        page_access_token: row.access_token,
        channel_name: row.name,
      });
      notifications.show({
        title: 'Page connected',
        message: `${row.name} is now linked to PostEaze.`,
        color: 'teal',
      });
      setPickerOpen(false);
      setPages([]);
      setSelectedId(null);
      await queryClient.invalidateQueries({ queryKey: channelKeys.lists() });
    } catch (e) {
      const msg = e instanceof Error ? e.message : 'Could not save channel';
      notifications.show({ title: 'Could not connect', message: msg, color: 'red' });
    } finally {
      setConnecting(false);
    }
  };

  const handleSyncAnalytics = async () => {
    if (metaChannelIds.length === 0) {
      notifications.show({
        title: 'Nothing to sync',
        message: 'Connect Instagram or Facebook first.',
        color: 'gray',
      });
      return;
    }
    try {
      setSyncing(true);
      const res = await syncMetaAnalytics();
      const ok = res.results.filter((r) => r.status === 'ok').length;
      const fail = res.results.length - ok;
      notifications.show({
        title: 'Meta analytics sync finished',
        message: `${ok} channel(s) updated${fail ? `, ${fail} issue(s)` : ''}.`,
        color: fail ? 'yellow' : 'teal',
      });
      await queryClient.invalidateQueries({ queryKey: channelKeys.lists() });
    } catch (e) {
      const msg = e instanceof Error ? e.message : 'Sync failed';
      notifications.show({ title: 'Sync failed', message: msg, color: 'red' });
    } finally {
      setSyncing(false);
    }
  };

  return (
    <Container size="xl" className={`fade-in ${styles.root}`}>
      <Box className={styles.signalStrip} aria-hidden />

      <Stack gap="xl">
        <Paper radius="lg" p="xl" className={styles.hero}>
          <ChannelPageHeader
            brand="facebook"
            title="Facebook"
            description="Connect Pages, then sync reach and engagement into PostEaze — one ledger for Meta."
            icon={<Icons.Facebook size={28} />}
            actionLabel="Connect page"
            actionLeftSection={
              oauthLoading || connecting ? <Loader size="sm" color="white" /> : <Icons.Plus size={20} />
            }
            onAction={handleConnectClick}
          />
          <Group mt="md" gap="sm" wrap="wrap">
            <Button
              variant="light"
              color="cyan"
              leftSection={syncing ? <Loader size="sm" /> : <Icons.ChartBar size={18} />}
              onClick={handleSyncAnalytics}
              disabled={syncing || channelsLoading}
            >
              Sync Meta analytics
            </Button>
            <Text size="sm" c="dimmed" maw={420}>
              Pulls Instagram + Facebook insights for your connected channels. May take up to a minute.
            </Text>
          </Group>
        </Paper>

        <Paper p="xl" radius="lg" withBorder className={styles.panel}>
          <Stack gap="md">
            <Group justify="space-between">
              <Title order={3} fw={800} style={{ letterSpacing: '-0.03em' }}>
                Connected pages
              </Title>
              <Badge size="lg" variant="light" color="blue">
                {facebookChannels.length} connected
              </Badge>
            </Group>
            {channelsLoading ? (
              <Group justify="center" py="xl">
                <Loader />
              </Group>
            ) : facebookChannels.length === 0 ? (
              <EmptyChannelPanel
                icon={<Icons.Facebook size={56} />}
                title="No Facebook pages connected yet"
                hint="Use Connect page to authorize and choose which Page to link."
              />
            ) : (
              <Stack gap="sm">
                {facebookChannels.map((ch) => (
                  <Paper
                    key={ch.id}
                    p="md"
                    withBorder
                    radius="md"
                    style={{ borderColor: 'var(--pe-border)' }}
                  >
                    <Group justify="space-between" wrap="nowrap" align="flex-start">
                      <div>
                        <Text fw={700}>{ch.channelName}</Text>
                        <Text size="sm" c="dimmed" className={styles.mono}>
                          Page ID {(ch.metadata as { page_id?: string })?.page_id ?? ch.provider_channel_id}
                        </Text>
                        {(ch.metadata as { category?: string })?.category && (
                          <Text size="xs" c="dimmed" mt={4}>
                            {(ch.metadata as { category?: string }).category}
                          </Text>
                        )}
                      </div>
                      <Badge color="blue" variant="outline">
                        Active
                      </Badge>
                    </Group>
                  </Paper>
                ))}
              </Stack>
            )}
          </Stack>
        </Paper>

        <Paper p="xl" radius="lg" withBorder className={styles.panel}>
          <Stack gap="sm">
            <Title order={3} fw={800} style={{ letterSpacing: '-0.03em' }}>
              Ledger
            </Title>
            <Text size="sm" c="dimmed">
              After syncing, daily reach and post metrics are stored for analytics. Run sync after major
              campaigns or weekly for steady reporting.
            </Text>
            <Divider />
            <Group gap="xs">
              <Badge variant="dot" color="cyan">
                Graph API
              </Badge>
              <Badge variant="dot" color="gray">
                Page + post insights
              </Badge>
            </Group>
          </Stack>
        </Paper>
      </Stack>

      <Modal
        opened={pickerOpen}
        onClose={() => !connecting && setPickerOpen(false)}
        title={<Text fw={800}>Choose a Page</Text>}
        centered
        radius="lg"
      >
        <Stack gap="md">
          <Text size="sm" c="dimmed">
            We found {pages.length} Page(s) on this account. Pick one to connect to PostEaze.
          </Text>
          <Radio.Group value={selectedId ?? undefined} onChange={setSelectedId}>
            <Stack gap="xs">
              {pages.map((p) => (
                <Radio
                  key={p.id}
                  value={p.id}
                  label={
                    <div>
                      <Text fw={600}>{p.name}</Text>
                      <Text size="xs" className={styles.mono} c="dimmed">
                        {p.id}
                      </Text>
                    </div>
                  }
                />
              ))}
            </Stack>
          </Radio.Group>
          <Group justify="flex-end" mt="sm">
            <Button variant="default" onClick={() => setPickerOpen(false)} disabled={connecting}>
              Cancel
            </Button>
            <Button onClick={handleConfirmPage} loading={connecting}>
              Connect this Page
            </Button>
          </Group>
        </Stack>
      </Modal>
    </Container>
  );
};

export default FacebookChannel;
