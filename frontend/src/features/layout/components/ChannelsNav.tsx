import { useState } from 'react';
import { NavLink, Stack, Group, Text, Collapse, UnstyledButton, Box } from '@mantine/core';
import { Icons } from '@/app/theme';
import { useLocation, useNavigate } from 'react-router-dom';
import { CHANNEL_NAV_ITEMS } from '../constants';
import classes from './ChannelsNav.module.css';

export const ChannelsNav = () => {
    const [opened, setOpened] = useState(true);
    const location = useLocation();
    const navigate = useNavigate();

    const isChannelActive = CHANNEL_NAV_ITEMS.some(item =>
        location.pathname.startsWith(item.path || '')
    );

    return (
        <Box>
            <UnstyledButton
                onClick={() => setOpened((o) => !o)}
                className={classes.control}
                data-active={isChannelActive || undefined}
            >
                <Group justify="space-between" gap={0}>
                    <Text fw={500} size="sm">
                        Channels
                    </Text>
                    <Icons.ChevronRight
                        size={16}
                        style={{
                            transform: opened ? 'rotate(90deg)' : 'none',
                            transition: 'transform 200ms ease',
                        }}
                    />
                </Group>
            </UnstyledButton>

            <Collapse in={opened}>
                <Stack gap={4} mt={4}>
                    {CHANNEL_NAV_ITEMS.map((item) => {
                        const Icon = item.icon;
                        const isActive = location.pathname === item.path;

                        return (
                            <NavLink
                                key={item.path}
                                label={item.label}
                                leftSection={<Icon size={20} />}
                                active={isActive}
                                onClick={() => item.path && navigate(item.path)}
                                className={classes.channelLink}
                            />
                        );
                    })}
                </Stack>
            </Collapse>
        </Box>
    );
};
