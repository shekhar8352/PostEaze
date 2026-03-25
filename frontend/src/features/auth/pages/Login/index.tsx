// src/features/auth/pages/Login/index.tsx
import { LoginForm } from '../../components/LoginForm';
import { AuthPageShell } from '../../components/AuthPageShell';
import { useNavigate } from 'react-router-dom';
import { getRedirectPath } from '../../utils';

const LoginPage = () => {
  const navigate = useNavigate();

  const switchToRegister = () => navigate('/register');
  const switchToForgotPassword = () => navigate('/forgot-password');

  const handleAuthSuccess = () => {
    const redirectPath = getRedirectPath();
    navigate(redirectPath);
  };

  const handleEmailNotVerified = (email: string, password: string) => {
    navigate('/email-verify', {
      state: { email, password, }
    });
  };

  return (
    <AuthPageShell variant="login">
      <LoginForm
        onToggleMode={switchToRegister}
        onForgotPassword={switchToForgotPassword}
        onSuccess={handleAuthSuccess}
        onEmailNotVerified={handleEmailNotVerified}
      />
    </AuthPageShell>
  );
};

export default LoginPage;
