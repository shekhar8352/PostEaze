// src/features/auth/components/signUp/LoginForm.tsx
import { Formik, Form, Field } from 'formik';
import { 
  Paper, 
  TextInput, 
  PasswordInput, 
  Button, 
  Title, 
  Text, 
  Stack, 
  Group, 
  Divider,
  Anchor,
  Box,
  Alert
} from '@mantine/core';
import { IconAlertTriangle } from '@tabler/icons-react';
import { notifications } from '@mantine/notifications';
import { 
  useLogin, 
  useGoogleAuth, 
  useFacebookAuth,
  useForgotPassword 
} from '../services/authQueries';
import { loginSchema } from '../validation/authSchema'
import type { LoginFormData } from '../types';

interface LoginFormProps {
  onToggleMode?: () => void;
  onForgotPassword?: () => void;
  onSuccess?: () => void;
  onEmailNotVerified?: (email: string, password: string) => void;
}

export const LoginForm = ({ 
  onToggleMode, 
  onForgotPassword,
  onSuccess,
  onEmailNotVerified
}: LoginFormProps) => {
  const login = useLogin();
  const googleAuth = useGoogleAuth();
  const facebookAuth = useFacebookAuth();
  const forgotPassword = useForgotPassword();

  const initialValues: LoginFormData = {
    email: '',
    password: '',
  };

  const handleLogin = async (values: LoginFormData) => {
    try {
      await login.mutateAsync(values);
      
      notifications.show({
        title: 'Welcome to PostEaze!',
        message: 'You have successfully logged in.',
        color: 'green',
      });
      
      onSuccess?.();
    } catch (error: any) {
      const errorMessage = error.message || 'Unable to log in. Please try again.';
      
      // Check if error is about email verification
      if (errorMessage.includes('verify your email') || errorMessage.includes('EMAIL_NOT_VERIFIED')) {
        notifications.show({
          title: 'Email Verification Required',
          message: 'Please verify your email before logging in. Check your inbox for the verification link.',
          color: 'orange',
          autoClose: 8000,
        });
        
        // Redirect to email verification screen
        onEmailNotVerified?.(values.email, values.password);
      } else {
        notifications.show({
          title: 'Login Failed',
          message: errorMessage,
          color: 'red',
        });
      }
    }
  };

  const handleGoogleAuth = async () => {
    try {
      await googleAuth.mutateAsync();
      
      notifications.show({
        title: 'Welcome to PostEaze!',
        message: 'Successfully logged in with Google.',
        color: 'green',
      });
      
      onSuccess?.();
    } catch (error: any) {
      notifications.show({
        title: 'Google Login Failed',
        message: error.message || 'Unable to log in with Google.',
        color: 'red',
      });
    }
  };

  const handleFacebookAuth = async () => {
    try {
      await facebookAuth.mutateAsync();
      
      notifications.show({
        title: 'Welcome to PostEaze!',
        message: 'Successfully logged in with Facebook.',
        color: 'green',
      });
      
      onSuccess?.();
    } catch (error: any) {
      notifications.show({
        title: 'Facebook Login Failed',
        message: error.message || 'Unable to log in with Facebook.',
        color: 'red',
      });
    }
  };

  const handleQuickPasswordReset = async (email: string) => {
    if (!email) {
      notifications.show({
        title: 'Email Required',
        message: 'Please enter your email address first.',
        color: 'orange',
      });
      return;
    }

    try {
      await forgotPassword.mutateAsync(email);
      notifications.show({
        title: 'Email Sent',
        message: 'Password reset link has been sent to your email.',
        color: 'green',
      });
    } catch (error: any) {
      notifications.show({
        title: 'Reset Failed',
        message: error.message || 'Unable to send reset email.',
        color: 'red',
      });
    }
  };

  return (
    <Paper radius="md" p="xl" withBorder shadow="sm">
      <Title order={2} ta="center" mb="lg">
        Welcome back to PostEaze!
      </Title>

      {/* Email Verification Warning */}
      <Alert 
        variant="light" 
        color="blue" 
        icon={<IconAlertTriangle size="1rem" />}
        mb="lg"
      >
        <Text size="sm">
          <strong>Email verification required:</strong> You must verify your email address before you can log in.
        </Text>
      </Alert>

      <Formik
        initialValues={initialValues}
        validationSchema={loginSchema}
        onSubmit={handleLogin}
      >
        {({ values, errors, touched, isSubmitting, setFieldValue, setFieldTouched }) => (
          <Form>
            <Stack gap="md">
              <Box>
                <Field name="email">
                  {({ field }: any) => (
                    <TextInput
                      {...field}
                      label="Email Address"
                      placeholder="Enter your verified email"
                      size="md"
                      error={touched.email && errors.email ? errors.email : null}
                      onChange={(e) => {
                        setFieldValue('email', e.target.value);
                        setFieldTouched('email', true);
                      }}
                    />
                  )}
                </Field>
              </Box>

              <Box>
                <Field name="password">
                  {({ field }: any) => (
                    <PasswordInput
                      {...field}
                      label="Password"
                      placeholder="Enter your password"
                      size="md"
                      error={touched.password && errors.password ? errors.password : null}
                      onChange={(e) => {
                        setFieldValue('password', e.target.value);
                        setFieldTouched('password', true);
                      }}
                    />
                  )}
                </Field>
              </Box>

              <Group justify="flex-end">
                <Anchor 
                  component="button" 
                  size="sm" 
                  type="button"
                  onClick={() => {
                    if (onForgotPassword) {
                      onForgotPassword();
                    } else {
                      handleQuickPasswordReset(values.email);
                    }
                  }}
                >
                  Forgot your password?
                </Anchor>
              </Group>

              <Button 
                type="submit" 
                size="md"
                loading={login.isPending || isSubmitting}
                fullWidth
              >
                Sign In
              </Button>

              <Divider label="Or continue with" labelPosition="center" my="lg" />

              <Group grow>
                <Button
                  variant="default"
                  size="md"
                  loading={googleAuth.isPending}
                  onClick={handleGoogleAuth}
                  type="button"
                  leftSection={<span>🔍</span>}
                >
                  Google
                </Button>
                <Button
                  variant="default"
                  size="md"
                  loading={facebookAuth.isPending}
                  onClick={handleFacebookAuth}
                  type="button"
                  leftSection={<span>📘</span>}
                >
                  Facebook
                </Button>
              </Group>

              {onToggleMode && (
                <Text c="dimmed" size="sm" ta="center" mt="md">
                  New to PostEaze?{' '}
                  <Anchor 
                    component="button" 
                    onClick={onToggleMode}
                    type="button"
                  >
                    Create an account
                  </Anchor>
                </Text>
              )}
            </Stack>
          </Form>
        )}
      </Formik>
    </Paper>
  );
};
