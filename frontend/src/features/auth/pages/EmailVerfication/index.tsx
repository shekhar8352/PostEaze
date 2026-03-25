import { EmailVerificationScreen } from '../../components/EmailVerificationScreen';
import { AuthPageShell } from '../../components/AuthPageShell';
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
    <AuthPageShell variant="verify">
      <EmailVerificationScreen
        email={state.email}
        password={state.password || ''}
        onBack={handleBack}
        onVerified={handleVerified}
      />
    </AuthPageShell>
  );
};

export default EmailVerificationPage;
