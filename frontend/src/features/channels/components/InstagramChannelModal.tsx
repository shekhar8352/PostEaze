import { Modal, Group, Box, Text } from '@mantine/core';
import { Icons } from '@/app/theme';
import { InstagramChannelForm } from './InstagramChannelForm';

interface InstagramChannelModalProps {
    opened: boolean;
    onClose: () => void;
}

export const InstagramChannelModal = ({
    opened,
    onClose,
}: InstagramChannelModalProps) => {
    const handleSuccess = () => {
        onClose();
    };

    return (
        <Modal
            opened={opened}
            onClose={onClose}
            title={
                <Group gap="sm">
                    <Box
                        style={{
                            width: 40,
                            height: 40,
                            borderRadius: '8px',
                            background:
                                'linear-gradient(45deg, #f09433 0%, #e6683c 25%, #dc2743 50%, #cc2366 75%, #bc1888 100%)',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                        }}
                    >
                        <Icons.Instagram size={24} color="white" />
                    </Box>
                    <Text size="lg" fw={600}>
                        Connect Instagram Channel
                    </Text>
                </Group>
            }
            size="md"
            centered
        >
            <InstagramChannelForm onSuccess={handleSuccess} onCancel={onClose} />
        </Modal>
    );
};
