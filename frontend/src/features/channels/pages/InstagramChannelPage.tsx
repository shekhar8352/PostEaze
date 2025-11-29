import { Container, Title, Text, Paper, Stack, Button, Group, Badge, Box } from '@mantine/core';
import { useAppDispatch, useAppSelector } from '@/app/store/hooks';
import { Icons } from '@/app/theme';
import { openCreateModal, closeCreateModal } from '../store/instagramChannelSlice';
import { useInstagramChannels } from '../services/instagramChannelQueries';
import { InstagramChannelList } from '../components/InstagramChannelList';
import { InstagramChannelModal } from '../components/InstagramChannelModal';

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
                {/* Enhanced Page Header */}
                <Paper
                    p="xl"
                    radius="lg"
                    style={{
                        background:
                            'linear-gradient(45deg, #f09433 0%, #e6683c 25%, #dc2743 50%, #cc2366 75%, #bc1888 100%)',
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
                                <Icons.Instagram size={32} color="white" />
                            </Box>
                            <div>
                                <Title order={2}>Instagram</Title>
                                <Text size="md" style={{ opacity: 0.95 }}>
                                    Manage your Instagram accounts and posts
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
                            onClick={handleOpenModal}
                        >
                            Connect Account
                        </Button>
                    </Group>
                </Paper>

                {/* Connected Accounts */}
                <Paper p="xl" radius="lg" shadow="md" withBorder>
                    <Stack gap="md">
                        <Group justify="space-between">
                            <Title order={3}>Connected Accounts</Title>
                            <Badge
                                size="lg"
                                variant="gradient"
                                gradient={{ from: '#f09433', to: '#bc1888', deg: 45 }}
                            >
                                {channels?.length || 0} connected
                            </Badge>
                        </Group>

                        <InstagramChannelList channels={channels || []} isLoading={isLoading} />
                    </Stack>
                </Paper>
            </Stack>

            {/* Create Channel Modal */}
            <InstagramChannelModal opened={isCreateModalOpen} onClose={handleCloseModal} />
        </Container>
    );
};

export default InstagramChannelPage;
