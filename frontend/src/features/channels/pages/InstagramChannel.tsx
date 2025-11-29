import { Container, Title, Text, Paper, Stack, Button, Group, Badge } from '@mantine/core';
import { Icons } from '@/app/theme';

const InstagramChannel = () => {
    return (
        <Container size="xl">
            <Stack gap="xl">
                {/* Page Header */}
                <Group justify="space-between" align="center">
                    <Group>
                        <Icons.Instagram size={40} style={{ color: '#E4405F' }} />
                        <div>
                            <Title order={1}>Instagram</Title>
                            <Text c="dimmed" size="sm">
                                Manage your Instagram accounts and posts
                            </Text>
                        </div>
                    </Group>
                    <Button
                        leftSection={<Icons.Plus size={20} />}
                        variant="gradient"
                        gradient={{ from: '#f09433', to: '#bc1888', deg: 45 }}
                    >
                        Connect Account
                    </Button>
                </Group>

                {/* Connected Accounts */}
                <Paper p="xl" radius="md" shadow="sm" withBorder>
                    <Stack gap="md">
                        <Group justify="space-between">
                            <Title order={3}>Connected Accounts</Title>
                            <Badge color="gray" variant="light">
                                0 connected
                            </Badge>
                        </Group>
                        <Text c="dimmed" ta="center" py="xl">
                            No Instagram accounts connected yet. Click "Connect Account" to get started.
                        </Text>
                    </Stack>
                </Paper>

                {/* Quick Stats - Placeholder */}
                <Paper p="xl" radius="md" shadow="sm" withBorder>
                    <Stack gap="md">
                        <Title order={3}>Quick Stats</Title>
                        <Text c="dimmed">
                            Connect your Instagram account to see analytics and insights.
                        </Text>
                    </Stack>
                </Paper>
            </Stack>
        </Container>
    );
};

export default InstagramChannel;
