import { Container, Title, Paper, Stack, Group, Badge } from '@mantine/core';
import { Icons } from '@/app/theme';
import { ChannelPageHeader } from '../components/ChannelPageHeader';
import { EmptyChannelPanel } from '../components/EmptyChannelPanel';

const FacebookChannel = () => {
    return (
        <Container size="xl" className="fade-in">
            <Stack gap="xl">
                <ChannelPageHeader
                    brand="facebook"
                    title="Facebook"
                    description="Manage Facebook Pages your organization publishes to."
                    icon={<Icons.Facebook size={28} />}
                    actionLabel="Connect page"
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
                                Connected pages
                            </Title>
                            <Badge size="lg" variant="light" color="gray">
                                0 connected
                            </Badge>
                        </Group>
                        <EmptyChannelPanel
                            icon={<Icons.Facebook size={56} />}
                            title="No Facebook pages connected yet"
                            hint='Use "Connect page" above to get started.'
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
                            title="Connect your Facebook page"
                            hint="Analytics and insights appear after you connect."
                        />
                    </Stack>
                </Paper>
            </Stack>
        </Container>
    );
};

export default FacebookChannel;
