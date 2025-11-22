import { Container, Title, Text, Button, Stack, Paper, Group, Box } from '@mantine/core';
import { useAuth } from '@/features/auth';
import { useNavigate } from 'react-router-dom';
import { IconLogout, IconUser } from '@tabler/icons-react';
import { notifications } from '@mantine/notifications';

const DashboardPage = () => {
    const { user, logout } = useAuth();
    const navigate = useNavigate();

    const handleLogout = async () => {
        try {
            logout();
            notifications.show({
                title: 'Logged Out',
                message: 'You have been successfully logged out.',
                color: 'blue',
            });
            navigate('/login');
        } catch (error) {
            notifications.show({
                title: 'Logout Failed',
                message: 'An error occurred while logging out.',
                color: 'red',
            });
        }
    };

    return (
        <Box
            style={{
                minHeight: '100vh',
                background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                padding: '2rem',
            }}
        >
            <Container size="lg">
                <Stack gap="xl">
                    <Paper p="xl" radius="md" shadow="md">
                        <Stack gap="lg">
                            <Group justify="space-between" align="center">
                                <div>
                                    <Title order={1}>Welcome to PostEaze!</Title>
                                    <Text size="lg" c="dimmed" mt="xs">
                                        Hello, {user?.name || user?.email || 'User'}! 👋
                                    </Text>
                                </div>
                                <Button
                                    leftSection={<IconLogout size={20} />}
                                    onClick={handleLogout}
                                    variant="light"
                                    color="red"
                                    size="md"
                                >
                                    Logout
                                </Button>
                            </Group>
                        </Stack>
                    </Paper>

                    <Paper p="xl" radius="md" shadow="md">
                        <Stack gap="md">
                            <Group>
                                <IconUser size={24} />
                                <Title order={3}>Your Profile</Title>
                            </Group>
                            <Text>
                                <strong>Email:</strong> {user?.email || 'N/A'}
                            </Text>
                            <Text>
                                <strong>Name:</strong> {user?.name || 'N/A'}
                            </Text>
                            <Text c="dimmed" size="sm">
                                This is a placeholder dashboard. More features coming soon!
                            </Text>
                        </Stack>
                    </Paper>
                </Stack>
            </Container>
        </Box>
    );
};

export default DashboardPage;
