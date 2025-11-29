import { Container, Title, Text, Paper, Stack, Button, Group, Badge, Box } from '@mantine/core';
import { Icons } from '@/app/theme';

const FacebookChannel = () => {
    return (
        <Container size="xl" className="fade-in">
            <Stack gap="xl">
                {/* Enhanced Page Header */}
                <Paper
                    p="xl"
                    radius="lg"
                    style={{
                        background: 'linear-gradient(135deg, #1877F2 0%, #0C63D4 100%)',
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
                                <Icons.Facebook size={32} color="white" />
                            </Box>
                            <div>
                                <Title order={2}>Facebook</Title>
                                <Text size="md" style={{ opacity: 0.95 }}>
                                    Manage your Facebook pages and posts
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
                            Connect Page
                        </Button>
                    </Group>
                </Paper>

                {/* Connected Pages */}
                <Paper p="xl" radius="lg" shadow="md" withBorder>
                    <Stack gap="md">
                        <Group justify="space-between">
                            <Title order={3}>Connected Pages</Title>
                            <Badge
                                size="lg"
                                variant="filled"
                                color="blue"
                            >
                                0 connected
                            </Badge>
                        </Group>
                        <Box
                            p="xl"
                            style={{
                                textAlign: 'center',
                                borderRadius: '12px',
                                background: 'rgba(24, 119, 242, 0.05)',
                            }}
                        >
                            <Icons.Facebook size={64} style={{ opacity: 0.3, marginBottom: '1rem' }} />
                            <Text c="dimmed" size="lg" fw={500}>
                                No Facebook pages connected yet
                            </Text>
                            <Text c="dimmed" size="sm" mt="xs">
                                Click "Connect Page" above to get started
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
                                background: 'rgba(24, 119, 242, 0.05)',
                            }}
                        >
                            <Icons.ChartBar size={64} style={{ opacity: 0.3, marginBottom: '1rem' }} />
                            <Text c="dimmed" size="lg" fw={500}>
                                Connect your Facebook page
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

export default FacebookChannel;
