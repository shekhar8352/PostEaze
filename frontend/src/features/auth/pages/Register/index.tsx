// src/features/auth/pages/Register/index.tsx
import { Container, Center, Box } from '@mantine/core';
import { RegisterForm } from '../../components/RegisterForm';
import { useNavigate } from 'react-router-dom';

const RegisterPage = () => {
  const navigate = useNavigate();

  const switchToLogin = () => navigate('/login');

  const handleAuthSuccess = () => {
    navigate('/dashboard');
  };

  const handleEmailSent = (email: string, password: string, name: string) => {
    navigate('/email-verification', { 
      state: { email, password, name }
    });
  };

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
            <RegisterForm 
              onToggleMode={switchToLogin}
              onSuccess={handleAuthSuccess}
              onEmailSent={handleEmailSent}
            />
          </Box>
        </Center>
      </Container>
    </Box>
  );
};

export default RegisterPage;
