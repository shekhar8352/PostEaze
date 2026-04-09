import { Group, Burger, Text, Menu, Avatar, UnstyledButton, ActionIcon, rem, useMantineColorScheme, useComputedColorScheme } from '@mantine/core';
import { Icons } from '@/app/theme';
import { useAuth } from '@/features/auth';
import { useNavigate } from 'react-router-dom';
import { notifications } from '@mantine/notifications';
import headerStyles from './AppHeader.module.css';

interface AppHeaderProps {
    mobileOpened: boolean;
    toggleMobile: () => void;
    desktopOpened: boolean;
    toggleDesktop: () => void;
}

export const AppHeader = ({ mobileOpened, toggleMobile }: AppHeaderProps) => {
    const { user, logout } = useAuth();
    const navigate = useNavigate();
    const { setColorScheme } = useMantineColorScheme();
    const computedScheme = useComputedColorScheme('light', { getInitialValueInEffect: true });

    const toggleColorScheme = () => {
        setColorScheme(computedScheme === 'dark' ? 'light' : 'dark');
    };

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
            className={headerStyles.bar}
        >
            <Group>
                <Burger
                    opened={mobileOpened}
                    onClick={toggleMobile}
                    hiddenFrom="sm"
                    size="sm"
                />
                <UnstyledButton
                    type="button"
                    className={headerStyles.brand}
                    onClick={() => navigate('/dashboard')}
                    aria-label="PostEaze home"
                >
                    Post<span className={headerStyles.brandAccent}>Eaze</span>
                </UnstyledButton>
            </Group>

            <Group gap="sm">
                <ActionIcon
                    onClick={toggleColorScheme}
                    variant="subtle"
                    color="gray"
                    size="lg"
                    radius="md"
                    aria-label="Toggle color scheme"
                >
                    {computedScheme === 'dark'
                        ? <Icons.Sun size={18} stroke={1.5} />
                        : <Icons.Moon size={18} stroke={1.5} />
                    }
                </ActionIcon>

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
                                <Text fw={600} size="sm" className={headerStyles.userLabel} style={{ lineHeight: 1 }} mr={3}>
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
        </Group>
    );
};
