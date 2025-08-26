// src/features/auth/pages/ForgotPassword/index.tsx
import { Container, Center, Box } from '@mantine/core';
import { ForgotPasswordForm } from '../../components/ForgotPasswordScreen';
import { useNavigate } from 'react-router-dom';

const ForgotPasswordPage = () => {
  const navigate = useNavigate();

  const handleBack = () => navigate('/login');

  return (
    <Box
      style={{
        minHeight: '100vh',
        background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
        display: 'flex',
        alignItems: 'center',
      }}
    >
      <Container size="sm" py="xl">
        <Center>
          <Box style={{ width: '100%', maxWidth: '450px' }}>
            <ForgotPasswordForm onBack={handleBack} />
          </Box>
        </Center>
      </Container>
    </Box>
  );
};

export default ForgotPasswordPage;
