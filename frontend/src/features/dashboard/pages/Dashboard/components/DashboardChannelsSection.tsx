import { Badge, Button, Card, Group, Paper, SimpleGrid, Stack, Text, ThemeIcon, Title } from '@mantine/core';
import { CHANNEL_COLORS, Icons } from '@/app/theme';
import type { DashboardChannelCard } from '../types';
import dashStyles from '../DashboardPage.module.css';

const CHANNEL_ICON_BG: Record<DashboardChannelCard['provider'], string> = {
    instagram: CHANNEL_COLORS.instagram.solid,
    facebook: CHANNEL_COLORS.facebook.solid,
    youtube: CHANNEL_COLORS.youtube.solid,
};

interface DashboardChannelsSectionProps {
    channels: DashboardChannelCard[];
    onNavigateChannel: (path: string) => void;
}

export function DashboardChannelsSection({ channels, onNavigateChannel }: DashboardChannelsSectionProps) {
    return (
        <Paper p="xl" radius="lg" withBorder={false} className={dashStyles.section}>
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
                                onClick={() => onNavigateChannel(channel.path)}
                                onKeyDown={(e) => {
                                    if (e.key === 'Enter' || e.key === ' ') {
                                        e.preventDefault();
                                        onNavigateChannel(channel.path);
                                    }
                                }}
                                tabIndex={0}
                                role="button"
                                aria-label={`Open ${channel.name}`}
                            >
                                <Stack gap="md">
                                    <Group justify="space-between" wrap="nowrap">
                                        <ThemeIcon size={48} radius="md" color="white" style={{ background: bg }}>
                                            <Icon size={26} />
                                        </ThemeIcon>
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
                                            channel.connected ? <Icons.Settings size={18} /> : <Icons.Plus size={18} />
                                        }
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            onNavigateChannel(channel.path);
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
    );
}
