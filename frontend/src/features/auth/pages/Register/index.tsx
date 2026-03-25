// src/features/auth/pages/Register/index.tsx
import { RegisterForm } from '../../components/RegisterForm';
import { AuthPageShell } from '../../components/AuthPageShell';
import { useNavigate } from 'react-router-dom';

const RegisterPage = () => {
  const navigate = useNavigate();

  const switchToLogin = () => navigate('/login');

  const handleAuthSuccess = () => {
    navigate('/dashboard');
  };

  const handleEmailSent = (email: string, password: string) => {
    navigate('/email-verify', { 
      state: { email, password }
    });
  };

  return (
    <AuthPageShell variant="register">
      <RegisterForm 
        onToggleMode={switchToLogin}
        onSuccess={handleAuthSuccess}
        onEmailSent={handleEmailSent}
      />
    </AuthPageShell>
  );
};

export default RegisterPage;
