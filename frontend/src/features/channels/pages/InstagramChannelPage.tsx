import { Container, Title, Paper, Stack, Group, Badge } from '@mantine/core';
import { useAppDispatch, useAppSelector } from '@/app/store/hooks';
import { Icons } from '@/app/theme';
import { openCreateModal, closeCreateModal } from '../store/instagramChannelSlice';
import { useInstagramChannels } from '../services/instagramChannelQueries';
import { InstagramChannelList } from '../components/InstagramChannelList';
import { InstagramChannelModal } from '../components/InstagramChannelModal';
import { ChannelPageHeader } from '../components/ChannelPageHeader';

export const InstagramChannelPage = () => {
    const dispatch = useAppDispatch();
    const { isCreateModalOpen } = useAppSelector((state) => state.instagramChannel);
    const { data: channels, isLoading } = useInstagramChannels();

    const handleOpenModal = () => {
        dispatch(openCreateModal());
    };

    const handleCloseModal = () => {
        dispatch(closeCreateModal());
    };

    return (
        <Container size="xl" className="fade-in">
            <Stack gap="xl">
                <ChannelPageHeader
                    brand="instagram"
                    title="Instagram"
                    description="Connect Business or Creator accounts to schedule and sync content."
                    icon={<Icons.Instagram size={28} />}
                    actionLabel="Connect account"
                    onAction={handleOpenModal}
                    actionLeftSection={<Icons.Plus size={20} />}
                />

                <Paper
                    p="xl"
                    radius="lg"
                    withBorder
                    style={{
                        borderColor: 'var(--pe-border)',
                        background: 'var(--pe-bg-elevated)',
                        boxShadow: 'var(--pe-shadow-sm)',
                    }}
                >
                    <Stack gap="md">
                        <Group justify="space-between">
                            <Title order={3} fw={700} c="var(--pe-text)">
                                Connected accounts
                            </Title>
                            <Badge size="lg" variant="light" color="gray">
                                {channels?.length ?? 0} connected
                            </Badge>
                        </Group>

                        <InstagramChannelList channels={channels || []} isLoading={isLoading} />
                    </Stack>
                </Paper>
            </Stack>

            <InstagramChannelModal opened={isCreateModalOpen} onClose={handleCloseModal} />
        </Container>
    );
};

export default InstagramChannelPage;
