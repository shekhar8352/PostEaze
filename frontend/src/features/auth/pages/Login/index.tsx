// src/features/auth/pages/Login/index.tsx
import { Container, Center, Box } from '@mantine/core';
import { LoginForm } from '../../components/LoginForm';
import { useNavigate } from 'react-router-dom';

const LoginPage = () => {
  const navigate = useNavigate();

  const switchToRegister = () => navigate('/register');
  const switchToForgotPassword = () => navigate('/forgot-password');

  const handleAuthSuccess = () => {
    navigate('/dashboard');
  };

  const handleEmailNotVerified = (email: string, password: string) => {
    navigate('/email-verification', { 
      state: { email, password, name: '' }
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
            <LoginForm 
              onToggleMode={switchToRegister}
              onForgotPassword={switchToForgotPassword}
              onSuccess={handleAuthSuccess}
              onEmailNotVerified={handleEmailNotVerified}
            />
          </Box>
        </Center>
      </Container>
    </Box>
  );
};

export default LoginPage;
