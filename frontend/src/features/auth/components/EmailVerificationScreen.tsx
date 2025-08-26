// src/features/auth/components/signUp/EmailVerificationScreen.tsx
import { useState, useEffect } from 'react';
import { 
  Paper, 
  Title, 
  Text, 
  Stack, 
  Button, 
  Alert,
  Progress,
  Anchor
} from '@mantine/core';
import { IconMail, IconCheck, IconAlertCircle, IconRefresh } from '@tabler/icons-react';
import { notifications } from '@mantine/notifications';
import { 
  useResendVerificationEmail,
  useCheckEmailVerification,
  useCompleteRegistration
} from '../services/authQueries';

interface EmailVerificationScreenProps {
  email: string;
  password: string;
  name: string;
  onBack?: () => void;
  onVerified?: () => void;
}

export const EmailVerificationScreen = ({ 
  email, 
  password, 
  name, 
  onBack, 
  onVerified 
}: EmailVerificationScreenProps) => {
  const [countdown, setCountdown] = useState(60);
  const [canResend, setCanResend] = useState(false);
  const [isChecking, setIsChecking] = useState(false);
  
  const resendEmail = useResendVerificationEmail();
  const checkVerification = useCheckEmailVerification();
  const completeRegistration = useCompleteRegistration();

  // Countdown timer for resend button
  useEffect(() => {
    if (countdown > 0) {
      const timer = setTimeout(() => setCountdown(countdown - 1), 1000);
      return () => clearTimeout(timer);
    } else {
      setCanResend(true);
    }
  }, [countdown]);

  // Auto-check verification status every 10 seconds
  useEffect(() => {
    const interval = setInterval(() => {
      handleCheckVerification(false); // Silent check
    }, 10000);

    return () => clearInterval(interval);
  }, []);

  const handleResendEmail = async () => {
    try {
      await resendEmail.mutateAsync({ email, password });
      
      notifications.show({
        title: 'Email Sent',
        message: 'Verification email has been sent again. Please check your inbox.',
        color: 'green',
      });
      
      setCountdown(60);
      setCanResend(false);
    } catch (error: any) {
      notifications.show({
        title: 'Resend Failed',
        message: error.message || 'Failed to send verification email.',
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
            title: 'Email Verified!',
            message: 'Completing your registration...',
            color: 'green',
          });
        }
        
        // Complete the registration process
        await completeRegistration.mutateAsync({ email, password, name });
        onVerified?.();
      } else {
        if (showNotification) {
          notifications.show({
            title: 'Not Verified Yet',
            message: 'Please check your email and click the verification link.',
            color: 'orange',
          });
        }
      }
    } catch (error: any) {
      if (showNotification) {
        notifications.show({
          title: 'Verification Check Failed',
          message: error.message || 'Unable to check verification status.',
          color: 'red',
        });
      }
    } finally {
      setIsChecking(false);
    }
  };

  return (
    <Paper radius="md" p="xl" withBorder shadow="sm">
      <Stack align="center" gap="md">
        <IconMail size={48} color="var(--mantine-color-blue-6)" />
        
        <Title order={2} ta="center">
          Verify Your Email
        </Title>

        <Text c="dimmed" size="sm" ta="center" maw={400}>
          We've sent a verification email to:
        </Text>
        
        <Text fw={600} ta="center" size="lg">
          {email}
        </Text>

        <Alert color="blue" variant="light" style={{ width: '100%' }}>
          <Text size="sm">
            <strong>Important:</strong> You must verify your email before you can access PostEaze. 
            Click the verification link in your email to complete registration.
          </Text>
        </Alert>

        <Stack gap="sm" style={{ width: '100%' }}>
          <Button 
            leftSection={<IconCheck size="1rem" />}
            onClick={() => handleCheckVerification(true)}
            loading={isChecking || completeRegistration.isPending}
            variant="filled"
            fullWidth
          >
            {isChecking ? 'Checking...' : 'I\'ve Verified My Email'}
          </Button>

          <Button 
            leftSection={<IconRefresh size="1rem" />}
            onClick={handleResendEmail}
            disabled={!canResend}
            loading={resendEmail.isPending}
            variant="light"
            fullWidth
          >
            {canResend ? 'Resend Verification Email' : `Resend in ${countdown}s`}
          </Button>

          {!canResend && (
            <Progress value={((60 - countdown) / 60) * 100} size="xs" />
          )}
        </Stack>

        <Alert color="yellow" variant="light" icon={<IconAlertCircle size="1rem" />}>
          <Text size="sm">
            <strong>Don't see the email?</strong> Check your spam folder. The email might take a few minutes to arrive.
          </Text>
        </Alert>

        {onBack && (
          <Anchor 
            component="button" 
            size="sm"
            onClick={onBack}
          >
            ← Back to Registration
          </Anchor>
        )}
      </Stack>
    </Paper>
  );
};
