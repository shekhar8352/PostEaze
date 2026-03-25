import { useState, useEffect } from 'react';
import { 
  Text, 
  Stack, 
  Button, 
  Alert,
  Progress,
  Anchor,
  Box,
} from '@mantine/core';
import { IconMail, IconCheck, IconAlertCircle, IconRefresh } from '@tabler/icons-react';
import { notifications } from '@mantine/notifications';
import { 
  useResendVerificationEmail,
  useCheckEmailVerification,
  useCompleteRegistration
} from '../services/authQueries';
import authForm from './authForm.module.css';

interface EmailVerificationScreenProps {
  email: string;
  password: string;
  onBack?: () => void;
  onVerified?: () => void;
}

export const EmailVerificationScreen = ({ 
  email, 
  password,  
  onBack, 
  onVerified 
}: EmailVerificationScreenProps) => {
  const [countdown, setCountdown] = useState(60);
  const [canResend, setCanResend] = useState(false);
  const [isChecking, setIsChecking] = useState(false);
  
  const resendEmail = useResendVerificationEmail();
  const checkVerification = useCheckEmailVerification();
  const completeRegistration = useCompleteRegistration();

  useEffect(() => {
    if (countdown > 0) {
      const timer = setTimeout(() => setCountdown(countdown - 1), 1000);
      return () => clearTimeout(timer);
    } else {
      setCanResend(true);
    }
  }, [countdown]);

  useEffect(() => {
    const interval = setInterval(() => {
      handleCheckVerification(false);
    }, 10000);

    return () => clearInterval(interval);
  }, []);

  const handleResendEmail = async () => {
    try {
      await resendEmail.mutateAsync({ email, password });
      
      notifications.show({
        title: 'Email sent',
        message: 'Verification message sent again. Please check your inbox.',
        color: 'green',
      });
      
      setCountdown(60);
      setCanResend(false);
    } catch (error: any) {
      notifications.show({
        title: 'Resend failed',
        message: error.message || 'Could not send verification email.',
        color: 'red',
      });
    }
  };

  const handleCheckVerification = async (showNotification = true) => {
    try {
      setIsChecking(true);
      const result = await checkVerification.mutateAsync({ email, password });
      
      if (result.isVerified) {
        if (showNotification) {
          notifications.show({
            title: 'Verified',
            message: 'Completing registration…',
            color: 'green',
          });
        }
        
        await completeRegistration.mutateAsync({ email, password });
        onVerified?.();
      } else {
        if (showNotification) {
          notifications.show({
            title: 'Not verified yet',
            message: 'Open the link in your email to continue.',
            color: 'orange',
          });
        }
      }
    } catch (error: any) {
      if (showNotification) {
        notifications.show({
          title: 'Check failed',
          message: error.message || 'Unable to verify status.',
          color: 'red',
        });
      }
    } finally {
      setIsChecking(false);
    }
  };

  return (
    <Box className={`${authForm.surface} ${authForm.centerStack}`}>
      <Stack align="center" gap="md">
        <IconMail size={44} color="var(--auth-surface-accent)" stroke={1.5} />
        
        <Text component="h2" className={authForm.title}>
          Verify your email
        </Text>

        <Text c="dimmed" size="sm" ta="center" maw={400}>
          We sent a message to:
        </Text>
        
        <Text fw={600} ta="center" size="md" c="var(--auth-surface-ink)">
          {email}
        </Text>

        <Alert color="blue" variant="light" style={{ width: '100%' }} className={authForm.alert}>
          <Text size="sm">
            <strong>Required:</strong> Open the verification link in that email before you can sign in.
          </Text>
        </Alert>

        <Stack gap="sm" style={{ width: '100%' }}>
          <Button 
            leftSection={<IconCheck size="1rem" />}
            onClick={() => handleCheckVerification(true)}
            loading={isChecking || completeRegistration.isPending}
            variant="filled"
            fullWidth
            classNames={{ root: authForm.primaryButton }}
          >
            {isChecking ? 'Checking…' : 'I have verified my email'}
          </Button>

          <Button 
            leftSection={<IconRefresh size="1rem" />}
            onClick={handleResendEmail}
            disabled={!canResend}
            loading={resendEmail.isPending}
            variant="light"
            fullWidth
          >
            {canResend ? 'Resend email' : `Resend in ${countdown}s`}
          </Button>

          {!canResend && (
            <Progress value={((60 - countdown) / 60) * 100} size="xs" color="blue" />
          )}
        </Stack>

        <Alert color="yellow" variant="light" icon={<IconAlertCircle size="1rem" />}>
          <Text size="sm">
            <strong>No message?</strong> Check spam and allow a few minutes for delivery.
          </Text>
        </Alert>

        {onBack && (
          <Anchor 
            component="button" 
            size="sm"
            onClick={onBack}
            className={authForm.footer}
          >
            ← Back to sign in
          </Anchor>
        )}
      </Stack>
    </Box>
  );
};
