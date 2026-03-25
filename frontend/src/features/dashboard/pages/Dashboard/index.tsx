import {
    Container,
    Title,
    Text,
    Button,
    Stack,
    Paper,
    SimpleGrid,
    Card,
    Group,
    Box,
    Badge,
    Loader,
    Center,
} from '@mantine/core';
import { useAuth } from '@/features/auth';
import { useNavigate } from 'react-router-dom';
import { Icons } from '@/app/theme';
import { CHANNEL_COLORS } from '@/app/theme';
import { useChannels } from '@/features/channels';
import dashStyles from './DashboardPage.module.css';

const CHANNEL_ICON_BG: Record<string, string> = {
    instagram: CHANNEL_COLORS.instagram.solid,
    facebook: CHANNEL_COLORS.facebook.solid,
    youtube: CHANNEL_COLORS.youtube.solid,
};

const DashboardPage = () => {
    const { user } = useAuth();
    const navigate = useNavigate();
    const { data: connectedChannels, isLoading } = useChannels();

    const displayName = user?.name || user?.email?.split('@')[0] || 'there';

    const channels = [
        {
            name: 'Instagram',
            icon: Icons.Instagram,
            provider: 'instagram' as const,
            path: '/channels/instagram',
            connected: connectedChannels?.some((ch) => ch.provider === 'instagram') || false,
        },
        {
            name: 'Facebook',
            icon: Icons.Facebook,
            provider: 'facebook' as const,
            path: '/channels/facebook',
            connected: connectedChannels?.some((ch) => ch.provider === 'facebook') || false,
        },
        {
            name: 'YouTube',
            icon: Icons.YouTube,
            provider: 'youtube' as const,
            path: '/channels/youtube',
            connected: connectedChannels?.some((ch) => ch.provider === 'youtube') || false,
        },
    ];

    const stats = [
        {
            title: 'Connected accounts',
            value: connectedChannels?.length ?? 0,
            icon: Icons.User,
        },
        {
            title: 'Scheduled posts',
            value: 0,
            icon: Icons.Calendar,
        },
        {
            title: 'Published today',
            value: 0,
            icon: Icons.CheckCircle,
        },
    ];

    if (isLoading) {
        return (
            <Center style={{ minHeight: '60vh' }}>
                <Loader size="md" color="blue" type="bars" />
            </Center>
        );
    }

    return (
        <Container size="xl" className="fade-in">
            <Stack gap="xl">
                <Paper className={dashStyles.hero} shadow="none" radius="lg" withBorder={false}>
                    <Box className={dashStyles.heroInner}>
                        <Title order={1} className={dashStyles.heroTitle}>
                            Welcome back, {displayName}
                        </Title>
                        <Text className={dashStyles.heroSubtitle}>
                            Overview of connected channels and publishing activity. Connect accounts to
                            start scheduling from one workspace.
                        </Text>
                    </Box>
                </Paper>

                <SimpleGrid cols={{ base: 1, sm: 3 }} spacing="md">
                    {stats.map((stat) => {
                        const Icon = stat.icon;
                        return (
                            <Card
                                key={stat.title}
                                padding="lg"
                                radius="md"
                                withBorder={false}
                                className={dashStyles.statCard}
                            >
                                <Group justify="space-between" align="flex-start" wrap="nowrap">
                                    <div>
                                        <Text size="xs" tt="uppercase" fw={600} c="dimmed" mb={6}>
                                            {stat.title}
                                        </Text>
                                        <Text size="xl" fw={700} c="var(--pe-text)">
                                            {stat.value}
                                        </Text>
                                    </div>
                                    <div className={dashStyles.statIcon}>
                                        <Icon size={22} stroke={1.75} />
                                    </div>
                                </Group>
                            </Card>
                        );
                    })}
                </SimpleGrid>

                <Paper
                    p="xl"
                    radius="lg"
                    withBorder={false}
                    className={dashStyles.section}
                >
                    <Stack gap="lg">
                        <div>
                            <Title order={2} className={dashStyles.sectionTitle} mb={6}>
                                Channels
                            </Title>
                            <Text className={dashStyles.sectionDesc}>
                                Link social accounts your team is approved to manage.
                            </Text>
                        </div>

                        <SimpleGrid cols={{ base: 1, sm: 3 }} spacing="md">
                            {channels.map((channel) => {
                                const Icon = channel.icon;
                                const bg = CHANNEL_ICON_BG[channel.provider];
                                return (
                                    <Card
                                        key={channel.path}
                                        padding="lg"
                                        radius="md"
                                        withBorder={false}
                                        className={dashStyles.channelCard}
                                        onClick={() => navigate(channel.path)}
                                        onKeyDown={(e) => {
                                            if (e.key === 'Enter' || e.key === ' ') {
                                                e.preventDefault();
                                                navigate(channel.path);
                                            }
                                        }}
                                        tabIndex={0}
                                        role="button"
                                        aria-label={`Open ${channel.name}`}
                                    >
                                        <Stack gap="md">
                                            <Group justify="space-between" wrap="nowrap">
                                                <div
                                                    className={dashStyles.channelIcon}
                                                    style={{ background: bg }}
                                                >
                                                    <Icon size={26} />
                                                </div>
                                                <Badge
                                                    size="sm"
                                                    variant="light"
                                                    color={channel.connected ? 'green' : 'gray'}
                                                >
                                                    {channel.connected ? 'Connected' : 'Not connected'}
                                                </Badge>
                                            </Group>
                                            <div>
                                                <Text fw={600} size="lg" c="var(--pe-text)">
                                                    {channel.name}
                                                </Text>
                                                <Text size="sm" c="dimmed" mt={4}>
                                                    Manage integration and content
                                                </Text>
                                            </div>
                                            <Button
                                                fullWidth
                                                variant="light"
                                                color="blue"
                                                leftSection={
                                                    channel.connected ? (
                                                        <Icons.Settings size={18} />
                                                    ) : (
                                                        <Icons.Plus size={18} />
                                                    )
                                                }
                                                onClick={(e) => {
                                                    e.stopPropagation();
                                                    navigate(channel.path);
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
