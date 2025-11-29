import { Container, Title, Text, Paper, Stack, Button, Group, Badge, Box } from '@mantine/core';
import { Icons } from '@/app/theme';

const InstagramChannel = () => {
    return (
        <Container size="xl" className="fade-in">
            <Stack gap="xl">
                {/* Enhanced Page Header */}
                <Paper
                    p="xl"
                    radius="lg"
                    style={{
                        background: 'linear-gradient(45deg, #f09433 0%, #e6683c 25%, #dc2743 50%, #cc2366 75%, #bc1888 100%)',
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
                    <Group justify="space-between" align="center" style={{ position: 'relative', zIndex: 1 }}>
                        <Group>
                            <Box
                                style={{
                                    width: 60,
                                    height: 60,
                                    borderRadius: '12px',
                                    background: 'rgba(255, 255, 255, 0.2)',
                                    display: 'flex',
                                    alignItems: 'center',
                                    justifyContent: 'center',
                                }}
                            >
                                <Icons.Instagram size={32} color="white" />
                            </Box>
                            <div>
                                <Title order={2}>Instagram</Title>
                                <Text size="md" style={{ opacity: 0.95 }}>
                                    Manage your Instagram accounts and posts
                                </Text>
                            </div>
                        </Group>
                        <Button
                            leftSection={<Icons.Plus size={20} />}
                            size="lg"
                            variant="white"
                            color="dark"
                            radius="md"
                            className="hover-scale"
                        >
                            Connect Account
                        </Button>
                    </Group>
                </Paper>

                {/* Connected Accounts */}
                <Paper p="xl" radius="lg" shadow="md" withBorder>
                    <Stack gap="md">
                        <Group justify="space-between">
                            <Title order={3}>Connected Accounts</Title>
                            <Badge
                                size="lg"
                                variant="gradient"
                                gradient={{ from: '#f09433', to: '#bc1888', deg: 45 }}
                            >
                                0 connected
                            </Badge>
                        </Group>
                        <Box
                            p="xl"
                            style={{
                                textAlign: 'center',
                                borderRadius: '12px',
                                background: 'linear-gradient(135deg, rgba(240, 148, 51, 0.05) 0%, rgba(188, 24, 136, 0.05) 100%)',
                            }}
                        >
                            <Icons.Instagram size={64} style={{ opacity: 0.3, marginBottom: '1rem' }} />
                            <Text c="dimmed" size="lg" fw={500}>
                                No Instagram accounts connected yet
                            </Text>
                            <Text c="dimmed" size="sm" mt="xs">
                                Click "Connect Account" above to get started
                            </Text>
                        </Box>
                    </Stack>
                </Paper>

                {/* Quick Stats */}
                <Paper p="xl" radius="lg" shadow="md" withBorder className="scale-in">
                    <Stack gap="md">
                        <Title order={3}>Quick Stats</Title>
                        <Box
                            p="xl"
                            style={{
                                textAlign: 'center',
                                borderRadius: '12px',
                                background: 'linear-gradient(135deg, rgba(240, 148, 51, 0.05) 0%, rgba(188, 24, 136, 0.05) 100%)',
                            }}
                        >
                            <Icons.ChartBar size={64} style={{ opacity: 0.3, marginBottom: '1rem' }} />
                            <Text c="dimmed" size="lg" fw={500}>
                                Connect your Instagram account
                            </Text>
                            <Text c="dimmed" size="sm" mt="xs">
                                View analytics and insights once connected
                            </Text>
                        </Box>
                    </Stack>
                </Paper>
            </Stack>
        </Container>
    );
};

export default InstagramChannel;
