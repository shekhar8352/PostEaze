import { Group, Burger, Text, Menu, Avatar, UnstyledButton, rem } from '@mantine/core';
import { Icons } from '@/app/theme';
import { useAuth } from '@/features/auth';
import { useNavigate } from 'react-router-dom';
import { notifications } from '@mantine/notifications';

interface AppHeaderProps {
    mobileOpened: boolean;
    toggleMobile: () => void;
    desktopOpened: boolean;
    toggleDesktop: () => void;
}

export const AppHeader = ({ mobileOpened, toggleMobile }: AppHeaderProps) => {
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
        <Group
            h="100%"
            px="md"
            justify="space-between"
            style={{
                background: 'rgba(255, 255, 255, 0.8)',
                backdropFilter: 'blur(10px)',
                borderBottom: '1px solid rgba(102, 126, 234, 0.1)',
            }}
        >
            <Group>
                <Burger
                    opened={mobileOpened}
                    onClick={toggleMobile}
                    hiddenFrom="sm"
                    size="sm"
                />
                <Text
                    size="24px"
                    fw={800}
                    variant="gradient"
                    gradient={{ from: 'blue', to: 'cyan', deg: 45 }}
                    style={{
                        letterSpacing: '-0.5px',
                        cursor: 'pointer',
                        transition: 'transform 0.3s ease',
                    }}
                    onMouseEnter={(e) => {
                        e.currentTarget.style.transform = 'scale(1.05)';
                    }}
                    onMouseLeave={(e) => {
                        e.currentTarget.style.transform = 'scale(1)';
                    }}
                >
                    PostEaze
                </Text>
            </Group>

            <Menu shadow="md" width={200} position="bottom-end">
                <Menu.Target>
                    <UnstyledButton>
                        <Group gap={7}>
                            <Avatar
                                src={user?.avatar}
                                alt={user?.name || user?.email}
                                radius="xl"
                                size={32}
                            />
                            <Text fw={600} size="sm" style={{ lineHeight: 1 }} mr={3}>
                                {user?.name || user?.email}
                            </Text>
                            <Icons.ChevronDown size={12} stroke={1.5} />
                        </Group>
                    </UnstyledButton>
                </Menu.Target>

                <Menu.Dropdown>
                    <Menu.Label>Account</Menu.Label>
                    <Menu.Item
                        leftSection={<Icons.User style={{ width: rem(14), height: rem(14) }} />}
                        onClick={() => navigate('/profile')}
                    >
                        Profile
                    </Menu.Item>
                    <Menu.Item
                        leftSection={<Icons.Settings style={{ width: rem(14), height: rem(14) }} />}
                        onClick={() => navigate('/settings')}
                    >
                        Settings
                    </Menu.Item>

                    <Menu.Divider />

                    <Menu.Item
                        color="red"
                        leftSection={<Icons.Logout style={{ width: rem(14), height: rem(14) }} />}
                        onClick={handleLogout}
                    >
                        Logout
                    </Menu.Item>
                </Menu.Dropdown>
            </Menu>
        </Group>
    );
};
