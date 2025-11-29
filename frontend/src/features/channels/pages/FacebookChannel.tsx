import { Container, Title, Text, Paper, Stack, Button, Group, Badge } from '@mantine/core';
import { Icons } from '@/app/theme';

const FacebookChannel = () => {
    return (
        <Container size="xl">
            <Stack gap="xl">
                {/* Page Header */}
                <Group justify="space-between" align="center">
                    <Group>
                        <Icons.Facebook size={40} style={{ color: '#1877F2' }} />
                        <div>
                            <Title order={1}>Facebook</Title>
                            <Text c="dimmed" size="sm">
                                Manage your Facebook pages and posts
                            </Text>
                        </div>
                    </Group>
                    <Button
                        leftSection={<Icons.Plus size={20} />}
                        color="blue"
                        variant="filled"
                    >
                        Connect Page
                    </Button>
                </Group>

                {/* Connected Pages */}
                <Paper p="xl" radius="md" shadow="sm" withBorder>
                    <Stack gap="md">
                        <Group justify="space-between">
                            <Title order={3}>Connected Pages</Title>
                            <Badge color="gray" variant="light">
                                0 connected
                            </Badge>
                        </Group>
                        <Text c="dimmed" ta="center" py="xl">
                            No Facebook pages connected yet. Click "Connect Page" to get started.
                        </Text>
                    </Stack>
                </Paper>

                {/* Quick Stats - Placeholder */}
                <Paper p="xl" radius="md" shadow="sm" withBorder>
                    <Stack gap="md">
                        <Title order={3}>Quick Stats</Title>
                        <Text c="dimmed">
                            Connect your Facebook page to see analytics and insights.
                        </Text>
                    </Stack>
                </Paper>
            </Stack>
        </Container>
    );
};

export default FacebookChannel;
