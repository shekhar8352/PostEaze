import { Container, Title, Paper, Stack, Group, Badge } from '@mantine/core';
import { Icons } from '@/app/theme';
import { ChannelPageHeader } from '../components/ChannelPageHeader';
import { EmptyChannelPanel } from '../components/EmptyChannelPanel';

const YouTubeChannel = () => {
    return (
        <Container size="xl" className="fade-in">
            <Stack gap="xl">
                <ChannelPageHeader
                    brand="youtube"
                    title="YouTube"
                    description="Connect channels for video publishing and metadata sync."
                    icon={<Icons.YouTube size={28} />}
                    actionLabel="Connect channel"
                    actionLeftSection={<Icons.Plus size={20} />}
                />

                <Paper
                    p="xl"
                    radius="lg"
                    withBorder
                    style={{
                        borderColor: 'var(--pe-border)',
                        background: 'var(--pe-bg-elevated)',
                        boxShadow: 'var(--pe-shadow-sm)',
                    }}
                >
                    <Stack gap="md">
                        <Group justify="space-between">
                            <Title order={3} fw={700} c="var(--pe-text)">
                                Connected channels
                            </Title>
                            <Badge size="lg" variant="light" color="gray">
                                0 connected
                            </Badge>
                        </Group>
                        <EmptyChannelPanel
                            icon={<Icons.YouTube size={56} />}
                            title="No YouTube channels connected yet"
                            hint='Use "Connect channel" above to get started.'
                        />
                    </Stack>
                </Paper>

                <Paper
                    p="xl"
                    radius="lg"
                    withBorder
                    style={{
                        borderColor: 'var(--pe-border)',
                        background: 'var(--pe-bg-elevated)',
                        boxShadow: 'var(--pe-shadow-sm)',
                    }}
                >
                    <Stack gap="md">
                        <Title order={3} fw={700} c="var(--pe-text)">
                            Quick stats
                        </Title>
                        <EmptyChannelPanel
                            icon={<Icons.ChartBar size={56} />}
                            title="Connect your YouTube channel"
                            hint="Analytics and insights appear after you connect."
                        />
                    </Stack>
                </Paper>
            </Stack>
        </Container>
    );
};

export default YouTubeChannel;
