import { Badge, Button, Group, Paper, Stack, Text } from '@mantine/core';
import { IconBrandGoogleDrive } from '@tabler/icons-react';
import { useGoogleOAuth } from '../hooks/useGoogleOAuth';
import {
  useConnectGoogleDrive,
  useDisconnectGoogleDrive,
  useGoogleDriveStatus,
} from '../hooks/useGoogleDriveQueries';

export function GoogleDriveConnectCard() {
  const { data: status, isLoading } = useGoogleDriveStatus();
  const connect = useConnectGoogleDrive();
  const disconnect = useDisconnectGoogleDrive();
  const { openOAuthPopup, isLoading: oauthLoading } = useGoogleOAuth('drive');

  const handleConnect = async () => {
    try {
      const code = await openOAuthPopup();
      await connect.mutateAsync(code);
    } catch {
      /* notifications in hooks */
    }
  };

  const busy = oauthLoading || connect.isPending || disconnect.isPending;

  return (
    <Paper p="md" radius="lg" withBorder>
      <Group justify="space-between" wrap="wrap" gap="md">
        <Group gap="sm">
          <IconBrandGoogleDrive size={28} color="var(--mantine-color-blue-6)" />
          <Stack gap={2}>
            <Text fw={600}>Google Drive</Text>
            <Text size="sm" c="dimmed">
              Import photos and videos from Drive; large videos stay in Drive for YouTube publishing.
            </Text>
          </Stack>
        </Group>
        <Group gap="sm">
          {isLoading ? (
            <Badge variant="light" color="gray">
              Checking…
            </Badge>
          ) : status?.connected ? (
            <>
              <Badge variant="light" color="green">
                {status.email || 'Connected'}
              </Badge>
              <Button
                variant="light"
                color="red"
                size="sm"
                loading={disconnect.isPending}
                onClick={() => disconnect.mutate()}
              >
                Disconnect
              </Button>
            </>
          ) : (
            <Button size="sm" loading={busy} onClick={handleConnect}>
              Connect Drive
            </Button>
          )}
        </Group>
      </Group>
    </Paper>
  );
}
