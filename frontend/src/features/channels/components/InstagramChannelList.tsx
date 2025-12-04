import { SimpleGrid, Box, Text, Stack } from '@mantine/core';
import { Icons } from '@/app/theme';
import { InstagramChannelCard } from './InstagramChannelCard';
import type { InstagramChannelDisplay } from '../types/instagram.types';

interface InstagramChannelListProps {
    channels: InstagramChannelDisplay[];
    isLoading?: boolean;
}

export const InstagramChannelList = ({
    channels,
    isLoading,
}: InstagramChannelListProps) => {
    // Loading State
    if (isLoading) {
        return (
            <Box
                p="xl"
                style={{
                    textAlign: 'center',
                    borderRadius: '12px',
                    background:
                        'linear-gradient(135deg, rgba(240, 148, 51, 0.05) 0%, rgba(188, 24, 136, 0.05) 100%)',
                }}
            >
                <Text c="dimmed" size="lg" fw={500}>
                    Loading channels...
                </Text>
            </Box>
        );
    }

    // Empty State
    if (!channels || channels.length === 0) {
        return (
            <Box
                p="xl"
                style={{
                    textAlign: 'center',
                    borderRadius: '12px',
                    background:
                        'linear-gradient(135deg, rgba(240, 148, 51, 0.05) 0%, rgba(188, 24, 136, 0.05) 100%)',
                }}
            >
                <Stack align="center" gap="md">
                    <Icons.Instagram
                        size={64}
                        style={{ opacity: 0.3 }}
                    />
                    <div>
                        <Text c="dimmed" size="lg" fw={500}>
                            No Instagram accounts connected yet
                        </Text>
                        <Text c="dimmed" size="sm" mt="xs">
                            Click "Connect Account" above to get started
                        </Text>
                    </div>
                </Stack>
            </Box>
        );
    }

    // Channel List
    return (
        <SimpleGrid cols={{ base: 1, sm: 2, md: 3 }} spacing="lg">
            {channels.map((channel) => (
                <InstagramChannelCard key={channel.id} channel={channel} />
            ))}
        </SimpleGrid>
    );
};
