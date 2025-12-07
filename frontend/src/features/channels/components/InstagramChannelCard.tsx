import { Card, Stack, Group, Avatar, Text, Badge, ActionIcon } from '@mantine/core';
import { Icons } from '@/app/theme';
import { useDeleteInstagramChannel } from '../services/instagramChannelQueries';
import { notifications } from '@mantine/notifications';
import type { InstagramChannelDisplay } from '../types/instagram.types';

interface InstagramChannelCardProps {
    channel: InstagramChannelDisplay;
}

export const InstagramChannelCard = ({ channel }: InstagramChannelCardProps) => {
    const deleteChannel = useDeleteInstagramChannel();

    const handleDelete = async () => {
        if (!confirm(`Are you sure you want to remove "${channel.channelName}"?`)) {
            return;
        }

        try {
            await deleteChannel.mutateAsync(channel.id);
            notifications.show({
                title: 'Channel Removed',
                message: 'Instagram channel has been successfully removed',
                color: 'green',
                icon: <Icons.CheckCircle size={18} />,
            });
        } catch (error) {
            notifications.show({
                title: 'Error',
                message: 'Failed to remove channel',
                color: 'red',
                icon: <Icons.XCircle size={18} />,
            });
        }
    };

    return (
        <Card
            shadow="sm"
            padding="lg"
            radius="md"
            withBorder
            className="shadow-lift"
        >
            <Stack gap="md">
                <Group justify="space-between">
                    <Group>
                        <Avatar
                            src={channel.metadata.profile_picture_url}
                            size="lg"
                            radius="xl"
                            style={{
                                border: '2px solid',
                                borderImage: 'linear-gradient(45deg, #f09433 0%, #e6683c 25%, #dc2743 50%, #cc2366 75%, #bc1888 100%) 1',
                            }}
                        >
                            <Icons.Instagram size={24} />
                        </Avatar>
                        <div>
                            <Text fw={600} size="lg">
                                {channel.channelName}
                            </Text>
                            <Text size="sm" c="dimmed">
                                {channel.username || channel.email}
                            </Text>
                        </div>
                    </Group>
                    <ActionIcon
                        variant="subtle"
                        color="red"
                        onClick={handleDelete}
                        loading={deleteChannel.isPending}
                    >
                        <Icons.Trash size={18} />
                    </ActionIcon>
                </Group>

                {/* Stats */}
                {(channel.followersCount !== undefined || channel.mediaCount !== undefined) && (
                    <Group gap="lg">
                        {channel.followersCount !== undefined && (
                            <Group gap="xs">
                                <Icons.User size={16} />
                                <Text size="sm">
                                    {channel.followersCount.toLocaleString()} followers
                                </Text>
                            </Group>
                        )}
                        {channel.mediaCount !== undefined && (
                            <Group gap="xs">
                                <Icons.FileText size={16} />
                                <Text size="sm">
                                    {channel.mediaCount.toLocaleString()} posts
                                </Text>
                            </Group>
                        )}
                    </Group>
                )}

                {/* Status Badge */}
                <Group justify="space-between" mt="auto">
                    <Badge
                        color={channel.isConnected ? 'green' : 'gray'}
                        variant="light"
                        leftSection={
                            channel.isConnected ? (
                                <Icons.CheckCircle size={12} />
                            ) : (
                                <Icons.XCircle size={12} />
                            )
                        }
                    >
                        {channel.isConnected ? 'Active' : 'Inactive'}
                    </Badge>

                    {channel.lastSyncedAt && (
                        <Text size="xs" c="dimmed">
                            Synced {new Date(channel.lastSyncedAt).toLocaleDateString()}
                        </Text>
                    )}
                </Group>
            </Stack>
        </Card>
    );
};
