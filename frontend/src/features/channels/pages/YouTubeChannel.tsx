import { Container, Title, Text, Paper, Stack, Button, Group, Badge } from '@mantine/core';
import { Icons } from '@/app/theme';

const YouTubeChannel = () => {
    return (
        <Container size="xl">
            <Stack gap="xl">
                {/* Page Header */}
                <Group justify="space-between" align="center">
                    <Group>
                        <Icons.YouTube size={40} style={{ color: '#FF0000' }} />
                        <div>
                            <Title order={1}>YouTube</Title>
                            <Text c="dimmed" size="sm">
                                Manage your YouTube channels and videos
                            </Text>
                        </div>
                    </Group>
                    <Button
                        leftSection={<Icons.Plus size={20} />}
                        color="red"
                        variant="filled"
                    >
                        Connect Channel
                    </Button>
                </Group>

                {/* Connected Channels */}
                <Paper p="xl" radius="md" shadow="sm" withBorder>
                    <Stack gap="md">
                        <Group justify="space-between">
                            <Title order={3}>Connected Channels</Title>
                            <Badge color="gray" variant="light">
                                0 connected
                            </Badge>
                        </Group>
                        <Text c="dimmed" ta="center" py="xl">
                            No YouTube channels connected yet. Click "Connect Channel" to get started.
                        </Text>
                    </Stack>
                </Paper>

                {/* Quick Stats - Placeholder */}
                <Paper p="xl" radius="md" shadow="sm" withBorder>
                    <Stack gap="md">
                        <Title order={3}>Quick Stats</Title>
                        <Text c="dimmed">
                            Connect your YouTube channel to see analytics and insights.
                        </Text>
                    </Stack>
                </Paper>
            </Stack>
        </Container>
    );
};

export default YouTubeChannel;
