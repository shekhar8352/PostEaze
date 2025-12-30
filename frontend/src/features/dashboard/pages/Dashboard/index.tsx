import { Container, Title, Text, Button, Stack, Paper, SimpleGrid, Card, Group, Box, RingProgress, Loader, Center } from '@mantine/core';
import { useAuth } from '@/features/auth';
import { useNavigate } from 'react-router-dom';
import { Icons } from '@/app/theme';
import { useChannels } from '@/features/channels';

const DashboardPage = () => {
    const { user } = useAuth();
    const navigate = useNavigate();
    const { data: connectedChannels, isLoading } = useChannels();

    const channels = [
        {
            name: 'Instagram',
            icon: Icons.Instagram,
            gradient: 'linear-gradient(45deg, #f09433 0%, #e6683c 25%, #dc2743 50%, #cc2366 75%, #bc1888 100%)',
            path: '/channels/instagram',
            connected: connectedChannels?.some(ch => ch.provider === 'instagram') || false,
        },
        {
            name: 'Facebook',
            icon: Icons.Facebook,
            gradient: 'linear-gradient(135deg, #1877F2 0%, #0C63D4 100%)',
            path: '/channels/facebook',
            connected: connectedChannels?.some(ch => ch.provider === 'facebook') || false,
        },
        {
            name: 'YouTube',
            icon: Icons.YouTube,
            gradient: 'linear-gradient(135deg, #FF0000 0%, #CC0000 100%)',
            path: '/channels/youtube',
            connected: connectedChannels?.some(ch => ch.provider === 'youtube') || false,
        },
    ];

    const stats = [
        {
            title: 'Connected Accounts',
            value: connectedChannels?.length || 0,
            icon: Icons.User,
            color: '#667eea',
            gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        },
        {
            title: 'Scheduled Posts',
            value: 0,
            icon: Icons.Calendar,
            color: '#11998e',
            gradient: 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)',
        },
        {
            title: 'Published Today',
            value: 0,
            icon: Icons.CheckCircle,
            color: '#f093fb',
            gradient: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
        },
    ];

    if (isLoading) {
        return (
            <Center style={{ height: '80vh' }}>
                <Loader size="xl" color="blue" variant="bars" />
            </Center>
        );
    }

    return (
        <Container size="xl" className="fade-in">
            <Stack gap="xl">
                {/* Welcome Banner with Gradient */}
                <Paper
                    p="xl"
                    radius="lg"
                    style={{
                        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                        color: 'white',
                        position: 'relative',
                        overflow: 'hidden',
                    }}
                    className="shadow-lift"
                >
                    <Box
                        style={{
                            position: 'absolute',
                            top: -50,
                            right: -50,
                            width: 200,
                            height: 200,
                            borderRadius: '50%',
                            background: 'rgba(255, 255, 255, 0.1)',
                            filter: 'blur(40px)',
                        }}
                    />
                    <Stack gap="md" style={{ position: 'relative', zIndex: 1 }}>
                        <Title order={1} style={{ fontSize: '2.5rem', fontWeight: 800 }}>
                            Welcome back, {user?.name || user?.email?.split('@')[0]}! 👋
                        </Title>
                        <Text size="lg" style={{ opacity: 0.95 }}>
                            Manage all your social media posts from one beautiful dashboard
                        </Text>
                    </Stack>
                </Paper>

                {/* Enhanced Stats Cards */}
                <SimpleGrid cols={{ base: 1, sm: 3 }} spacing="lg">
                    {stats.map((stat) => {
                        const Icon = stat.icon;
                        return (
                            <Card
                                key={stat.title}
                                shadow="md"
                                padding="xl"
                                radius="lg"
                                withBorder
                                className="shadow-lift"
                                style={{
                                    background: 'white',
                                    border: '1px solid rgba(102, 126, 234, 0.1)',
                                }}
                            >
                                <Group justify="space-between" mb="md">
                                    <Box
                                        style={{
                                            width: 60,
                                            height: 60,
                                            borderRadius: '12px',
                                            background: stat.gradient,
                                            display: 'flex',
                                            alignItems: 'center',
                                            justifyContent: 'center',
                                            boxShadow: `0 8px 16px ${stat.color}40`,
                                        }}
                                    >
                                        <Icon size={28} color="white" />
                                    </Box>
                                    <RingProgress
                                        size={60}
                                        thickness={6}
                                        sections={[{ value: 0, color: stat.color }]}
                                        label={
                                            <Text size="xs" ta="center" fw={700}>
                                                0%
                                            </Text>
                                        }
                                    />
                                </Group>
                                <Text size="sm" c="dimmed" fw={600} tt="uppercase" mb={4}>
                                    {stat.title}
                                </Text>
                                <Title order={2} style={{ fontSize: '2.5rem', fontWeight: 800 }}>
                                    {stat.value}
                                </Title>
                            </Card>
                        );
                    })}
                </SimpleGrid>

                {/* Enhanced Channel Cards */}
                <Paper
                    p="xl"
                    radius="lg"
                    shadow="sm"
                    withBorder
                    style={{
                        background: 'white',
                        border: '1px solid rgba(102, 126, 234, 0.1)',
                    }}
                >
                    <Stack gap="xl">
                        <Group justify="space-between" align="center">
                            <div>
                                <Title order={2} mb={4}>
                                    Connect Your Channels
                                </Title>
                                <Text c="dimmed">
                                    Link your social media accounts to start managing posts
                                </Text>
                            </div>
                        </Group>

                        <SimpleGrid cols={{ base: 1, sm: 3 }} spacing="lg">
                            {channels.map((channel, index) => {
                                const Icon = channel.icon;
                                return (
                                    <Card
                                        key={index}
                                        shadow="md"
                                        padding="xl"
                                        radius="lg"
                                        withBorder
                                        className="shadow-lift"
                                        style={{
                                            cursor: 'pointer',
                                            background: 'white',
                                            border: '2px solid transparent',
                                            backgroundImage: `linear-gradient(white, white), ${channel.gradient}`,
                                            backgroundOrigin: 'border-box',
                                            backgroundClip: 'padding-box, border-box',
                                            transition: 'all 0.3s ease',
                                        }}
                                        onClick={() => navigate(channel.path)}
                                        onMouseEnter={(e) => {
                                            e.currentTarget.style.transform = 'translateY(-4px)';
                                            e.currentTarget.style.boxShadow = '0 12px 24px rgba(0,0,0,0.1)';
                                        }}
                                        onMouseLeave={(e) => {
                                            e.currentTarget.style.transform = 'translateY(0)';
                                            e.currentTarget.style.boxShadow = '';
                                        }}
                                    >
                                        <Stack gap="lg" align="center">
                                            <Box
                                                style={{
                                                    width: 80,
                                                    height: 80,
                                                    borderRadius: '20px',
                                                    background: channel.gradient,
                                                    display: 'flex',
                                                    alignItems: 'center',
                                                    justifyContent: 'center',
                                                    boxShadow: '0 12px 24px rgba(0,0,0,0.15)',
                                                }}
                                            >
                                                <Icon size={40} color="white" />
                                            </Box>
                                            <div style={{ textAlign: 'center' }}>
                                                <Text fw={700} size="xl" mb={4}>
                                                    {channel.name}
                                                </Text>
                                                <Text size="sm" c="dimmed">
                                                    {channel.connected ? 'Connected' : 'Not connected'}
                                                </Text>
                                            </div>
                                            <Button
                                                leftSection={channel.connected ? <Icons.Settings size={18} /> : <Icons.Plus size={18} />}
                                                fullWidth
                                                size="md"
                                                radius="md"
                                                style={{
                                                    background: channel.gradient,
                                                    border: 'none',
                                                }}
                                                styles={{
                                                    root: {
                                                        '&:hover': {
                                                            transform: 'scale(1.02)',
                                                        },
                                                    },
                                                }}
                                            >
                                                {channel.connected ? 'Manage' : 'Connect'}
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
