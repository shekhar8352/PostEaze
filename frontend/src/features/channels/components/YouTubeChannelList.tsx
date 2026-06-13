import { ActionIcon, Avatar, Card, Group, SimpleGrid, Stack, Text } from '@mantine/core';
import { Icons } from '@/app/theme';
import { notifications } from '@mantine/notifications';
import { EmptyChannelPanel } from './EmptyChannelPanel';
import { useDeleteYouTubeChannel } from '../services/youtubeChannelQueries';
import type { YouTubeChannelDisplay } from '../types/youtube.types';

interface YouTubeChannelListProps {
  channels: YouTubeChannelDisplay[];
  isLoading?: boolean;
}

export function YouTubeChannelList({ channels, isLoading }: YouTubeChannelListProps) {
  const deleteChannel = useDeleteYouTubeChannel();

  if (isLoading) {
    return (
      <Text c="dimmed" size="sm" ta="center" py="xl">
        Loading channels…
      </Text>
    );
  }

  if (!channels?.length) {
    return (
      <EmptyChannelPanel
        icon={<Icons.YouTube size={56} />}
        title="No YouTube channels connected yet"
        hint='Use "Connect channel" above to get started.'
      />
    );
  }

  return (
    <SimpleGrid cols={{ base: 1, sm: 2, md: 3 }} spacing="md">
      {channels.map((channel) => (
        <Card key={channel.id} padding="lg" radius="md" withBorder>
          <Stack gap="md">
            <Group justify="space-between">
              <Group>
                <Avatar size="lg" radius="xl">
                  <Icons.YouTube size={24} />
                </Avatar>
                <div>
                  <Text fw={600}>{channel.channelName}</Text>
                  <Text size="sm" c="dimmed">
                    {channel.email || channel.channelUrl || 'YouTube'}
                  </Text>
                </div>
              </Group>
              <ActionIcon
                variant="subtle"
                color="red"
                loading={deleteChannel.isPending}
                onClick={async () => {
                  if (!confirm(`Remove "${channel.channelName}"?`)) return;
                  try {
                    await deleteChannel.mutateAsync(channel.id);
                    notifications.show({
                      title: 'Removed',
                      message: 'YouTube channel disconnected',
                      color: 'green',
                    });
                  } catch {
                    notifications.show({
                      title: 'Error',
                      message: 'Failed to remove channel',
                      color: 'red',
                    });
                  }
                }}
              >
                <Icons.Trash size={18} />
              </ActionIcon>
            </Group>
          </Stack>
        </Card>
      ))}
    </SimpleGrid>
  );
}
