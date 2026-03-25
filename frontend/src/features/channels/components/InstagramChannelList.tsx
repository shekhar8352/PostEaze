import { SimpleGrid, Box, Text } from '@mantine/core';
import { Icons } from '@/app/theme';
import { InstagramChannelCard } from './InstagramChannelCard';
import { EmptyChannelPanel } from './EmptyChannelPanel';
import type { InstagramChannelDisplay } from '../types/instagram.types';

interface InstagramChannelListProps {
    channels: InstagramChannelDisplay[];
    isLoading?: boolean;
}

export const InstagramChannelList = ({ channels, isLoading }: InstagramChannelListProps) => {
    if (isLoading) {
        return (
            <Box
                p="xl"
                style={{
                    textAlign: 'center',
                    borderRadius: 'var(--pe-radius-md)',
                    background: 'var(--pe-bg-subtle)',
                    border: '1px solid var(--pe-border)',
                }}
            >
                <Text c="dimmed" size="sm" fw={500}>
                    Loading channels…
                </Text>
            </Box>
        );
    }

    if (!channels || channels.length === 0) {
        return (
            <EmptyChannelPanel
                icon={<Icons.Instagram size={56} />}
                title="No Instagram accounts connected yet"
                hint='Use "Connect account" above to get started.'
            />
        );
    }

    return (
        <SimpleGrid cols={{ base: 1, sm: 2, md: 3 }} spacing="md">
            {channels.map((channel) => (
                <InstagramChannelCard key={channel.id} channel={channel} />
            ))}
        </SimpleGrid>
    );
};
