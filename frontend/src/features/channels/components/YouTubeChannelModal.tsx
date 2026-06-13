import { Button, Group, Modal, Stack, Text } from '@mantine/core';
import { Icons } from '@/app/theme';
import { useGoogleOAuth } from '@/features/integrations/hooks/useGoogleOAuth';
import { useCreateYouTubeChannel } from '../services/youtubeChannelQueries';

interface YouTubeChannelModalProps {
  opened: boolean;
  onClose: () => void;
}

export function YouTubeChannelModal({ opened, onClose }: YouTubeChannelModalProps) {
  const { openOAuthPopup, isLoading: oauthLoading } = useGoogleOAuth('youtube');
  const createChannel = useCreateYouTubeChannel();

  const handleConnect = async () => {
    try {
      const code = await openOAuthPopup();
      await createChannel.mutateAsync(code);
      onClose();
    } catch {
      /* notifications in hooks */
    }
  };

  return (
    <Modal opened={opened} onClose={onClose} title="Connect YouTube channel" centered>
      <Stack gap="md">
        <Group gap="sm">
          <Icons.YouTube size={32} />
          <Text size="sm" c="dimmed">
            Authorize PostEaze to upload videos to your YouTube channel. Large Drive-backed
            videos stream directly without copying to blob storage.
          </Text>
        </Group>
        <Button
          fullWidth
          leftSection={<Icons.YouTube size={18} />}
          loading={oauthLoading || createChannel.isPending}
          onClick={handleConnect}
        >
          Sign in with Google
        </Button>
      </Stack>
    </Modal>
  );
}
