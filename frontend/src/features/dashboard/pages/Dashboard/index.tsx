import { Container, Title, Text, Button, Stack, Paper, SimpleGrid, Card } from '@mantine/core';
import { useAuth } from '@/features/auth';
import { useNavigate } from 'react-router-dom';
import { IconBrandInstagram, IconBrandFacebook, IconBrandYoutube, IconPlus } from '@tabler/icons-react';

const DashboardPage = () => {
    const { user } = useAuth();
    const navigate = useNavigate();

    const channels = [
        {
            name: 'Instagram',
            icon: IconBrandInstagram,
            color: 'linear-gradient(45deg, #f09433 0%, #e6683c 25%, #dc2743 50%, #cc2366 75%, #bc1888 100%)',
            path: '/channels/instagram',
            connected: false,
        },
        {
            name: 'Facebook',
            icon: IconBrandFacebook,
            color: '#1877F2',
            path: '/channels/facebook',
            connected: false,
        },
        {
            name: 'YouTube',
            icon: IconBrandYoutube,
            color: '#FF0000',
            path: '/channels/youtube',
            connected: false,
        },
    ];

    return (
        <Container size="xl">
            <Stack gap="xl">
                {/* Welcome Section */}
                <Paper p="xl" radius="md" shadow="sm" withBorder>
                    <Stack gap="md">
                        <Title order={1}>Welcome back, {user?.name || user?.email}! 👋</Title>
                        <Text c="dimmed" size="lg">
                            Manage all your social media posts from one place
                        </Text>
                    </Stack>
                </Paper>

                {/* Quick Stats */}
                <SimpleGrid cols={{ base: 1, sm: 3 }} spacing="lg">
                    <Card shadow="sm" padding="lg" radius="md" withBorder>
                        <Stack gap="xs">
                            <Text size="sm" c="dimmed" fw={500}>
                                Connected Accounts
                            </Text>
                            <Text size="xl" fw={700}>
                                0
                            </Text>
                        </Stack>
                    </Card>
                    <Card shadow="sm" padding="lg" radius="md" withBorder>
                        <Stack gap="xs">
                            <Text size="sm" c="dimmed" fw={500}>
                                Scheduled Posts
                            </Text>
                            <Text size="xl" fw={700}>
                                0
                            </Text>
                        </Stack>
                    </Card>
                    <Card shadow="sm" padding="lg" radius="md" withBorder>
                        <Stack gap="xs">
                            <Text size="sm" c="dimmed" fw={500}>
                                Published Today
                            </Text>
                            <Text size="xl" fw={700}>
                                0
                            </Text>
                        </Stack>
                    </Card>
                </SimpleGrid>

                {/* Connect Channels */}
                <Paper p="xl" radius="md" shadow="sm" withBorder>
                    <Stack gap="lg">
                        <Title order={2}>Connect Your Channels</Title>
                        <SimpleGrid cols={{ base: 1, sm: 3 }} spacing="lg">
                            {channels.map((channel) => {
                                const Icon = channel.icon;
                                return (
                                    <Card
                                        key={channel.name}
                                        shadow="sm"
                                        padding="lg"
                                        radius="md"
                                        withBorder
                                        style={{ cursor: 'pointer' }}
                                        onClick={() => navigate(channel.path)}
                                    >
                                        <Stack gap="md" align="center">
                                            <Icon
                                                size={48}
                                                style={{
                                                    color: typeof channel.color === 'string' ? channel.color : undefined,
                                                    background: typeof channel.color !== 'string' ? channel.color : undefined,
                                                    WebkitBackgroundClip: typeof channel.color !== 'string' ? 'text' : undefined,
                                                    WebkitTextFillColor: typeof channel.color !== 'string' ? 'transparent' : undefined,
                                                }}
                                            />
                                            <Text fw={600} size="lg">
                                                {channel.name}
                                            </Text>
                                            <Button
                                                leftSection={<IconPlus size={16} />}
                                                variant="light"
                                                fullWidth
                                            >
                                                Connect
                                            </Button>
                                        </Stack>
                                    </Card>
                                );
                            })}
                        </SimpleGrid>
                    </Stack>
                </Paper>
            </Stack>
        </Container>
    );
};

export default DashboardPage;
