import { ForgotPasswordForm } from '../../components/ForgotPasswordScreen';
import { AuthPageShell } from '../../components/AuthPageShell';
import { useNavigate } from 'react-router-dom';

const ForgotPasswordPage = () => {
  const navigate = useNavigate();

  const handleBack = () => navigate('/login');

  return (
    <AuthPageShell variant="forgot">
      <ForgotPasswordForm onBack={handleBack} />
    </AuthPageShell>
  );
};

export default ForgotPasswordPage;
