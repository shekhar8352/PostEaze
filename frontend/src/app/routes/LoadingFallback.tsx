import { Center, Loader, Stack, Text } from '@mantine/core';

export const LoadingFallback = () => {
    return (
        <Center style={{ minHeight: '100vh' }}>
            <Stack align="center" gap="md">
                <Loader size="lg" />
                <Text c="dimmed">Loading...</Text>
            </Stack>
        </Center>
    );
};
