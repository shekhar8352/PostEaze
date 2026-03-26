import { Stack, NavLink, Divider, Box, ScrollArea } from '@mantine/core';
import { useLocation, useNavigate } from 'react-router-dom';
import { Icons } from '@/app/theme';
import { MAIN_NAV_ITEMS, BOTTOM_NAV_ITEMS } from '../constants';
import { ChannelsNav } from './ChannelsNav';
import { useAuth } from '@/features/auth';
import { notifications } from '@mantine/notifications';
import sidebarStyles from './MainSidebar.module.css';

export const MainSidebar = () => {
    const location = useLocation();
    const navigate = useNavigate();
    const { logout } = useAuth();

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
        <Box className={sidebarStyles.shell}>
            <ScrollArea
                style={{ flex: 1 }}
                type="auto"
                offsetScrollbars
                scrollbarSize={8}
            >
                <Stack gap="xs" p="md">
                    {/* Main Navigation */}
                    {MAIN_NAV_ITEMS.map((item) => {
                        const Icon = item.icon;
                        const isActive = location.pathname === item.path;

                        return (
                            <NavLink
                                key={item.path}
                                label={item.label}
                                leftSection={<Icon size={20} />}
                                active={isActive}
                                onClick={() => item.path && navigate(item.path)}
                            />
                        );
                    })}

                    <Divider my="sm" />

                    {/* Channels Navigation */}
                    <ChannelsNav />

                    <Divider my="sm" />

                    {/* Bottom Navigation */}
                    {BOTTOM_NAV_ITEMS.map((item) => {
                        const Icon = item.icon;
                        const isActive = location.pathname === item.path;

                        return (
                            <NavLink
                                key={item.path}
                                label={item.label}
                                leftSection={<Icon size={20} />}
                                active={isActive}
                                onClick={() => item.path && navigate(item.path)}
                            />
                        );
                    })}
                </Stack>
            </ScrollArea>

            {/* Logout Button - Fixed at bottom */}
            <Box p="md" className={sidebarStyles.footer}>
                <NavLink
                    label="Logout"
                    leftSection={<Icons.Logout size={20} />}
                    onClick={handleLogout}
                    color="red"
                />
            </Box>
        </Box>
    );
};
