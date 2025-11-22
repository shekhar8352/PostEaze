import { Container, Title, Text, Button, Stack, Center, Box } from '@mantine/core';
import { useNavigate } from 'react-router-dom';
import { IconHome } from '@tabler/icons-react';

export const NotFound = () => {
    const navigate = useNavigate();

    return (
        <Box
            style={{
                minHeight: '100vh',
                background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                display: 'flex',
                alignItems: 'center',
            }}
        >
            <Container size="sm">
                <Center>
                    <Stack align="center" gap="lg">
                        <Title
                            order={1}
                            style={{
                                fontSize: '120px',
                                fontWeight: 900,
                                color: 'white',
                                textShadow: '0 4px 6px rgba(0,0,0,0.1)',
                            }}
                        >
                            404
                        </Title>
                        <Title order={2} c="white" ta="center">
                            Page Not Found
                        </Title>
                        <Text size="lg" c="white" ta="center" maw={400}>
                            The page you're looking for doesn't exist or has been moved.
                        </Text>
                        <Button
                            size="lg"
                            leftSection={<IconHome size={20} />}
                            onClick={() => navigate('/')}
                            variant="white"
                            color="violet"
                        >
                            Go Home
                        </Button>
                    </Stack>
                </Center>
            </Container>
        </Box>
    );
};
