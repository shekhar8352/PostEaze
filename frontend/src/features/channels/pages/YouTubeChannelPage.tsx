import { useState } from 'react';
import { Badge, Container, Group, Paper, Stack, Title } from '@mantine/core';
import { Icons } from '@/app/theme';
import { ChannelPageHeader } from '../components/ChannelPageHeader';
import { YouTubeChannelList } from '../components/YouTubeChannelList';
import { YouTubeChannelModal } from '../components/YouTubeChannelModal';
import { useYouTubeChannels } from '../services/youtubeChannelQueries';

export default function YouTubeChannelPage() {
  const [modalOpen, setModalOpen] = useState(false);
  const { data: channels, isLoading } = useYouTubeChannels();

  return (
    <Container size="xl" className="fade-in">
      <Stack gap="xl">
        <ChannelPageHeader
          brand="youtube"
          title="YouTube"
          description="Connect channels for video publishing from Media Workspace."
          icon={<Icons.YouTube size={28} />}
          actionLabel="Connect channel"
          onAction={() => setModalOpen(true)}
          actionLeftSection={<Icons.Plus size={20} />}
        />

        <Paper p="xl" radius="lg" withBorder>
          <Stack gap="md">
            <Group justify="space-between">
              <Title order={3} fw={700}>
                Connected channels
              </Title>
              <Badge size="lg" variant="light" color="gray">
                {channels?.length ?? 0} connected
              </Badge>
            </Group>
            <YouTubeChannelList channels={channels ?? []} isLoading={isLoading} />
          </Stack>
        </Paper>
      </Stack>

      <YouTubeChannelModal opened={modalOpen} onClose={() => setModalOpen(false)} />
    </Container>
  );
}
