// src/features/auth/pages/EmailVerification/index.tsx
import { Container, Center, Box } from '@mantine/core';
import { EmailVerificationScreen } from '../../components/EmailVerificationScreen';
import { useNavigate, useLocation } from 'react-router-dom';
import { useEffect } from 'react';

interface LocationState {
  email?: string;
  password?: string;
  name?: string;
}

const EmailVerificationPage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const state = location.state as LocationState;

  useEffect(() => {
    // Redirect to login if no email is provided
    if (!state?.email) {
      navigate('/login');
    }
  }, [state, navigate]);

  const handleBack = () => navigate('/login');
  
  const handleVerified = () => {
    navigate('/dashboard');
  };

  if (!state?.email) {
    return null;
  }

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
            <EmailVerificationScreen
              email={state.email}
              password={state.password || ''}
              name={state.name || ''}
              onBack={handleBack}
              onVerified={handleVerified}
            />
          </Box>
        </Center>
      </Container>
    </Box>
  );
};

export default EmailVerificationPage;
