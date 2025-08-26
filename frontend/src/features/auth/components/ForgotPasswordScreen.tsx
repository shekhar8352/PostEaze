// src/features/auth/components/signUp/ForgotPasswordForm.tsx
import { useState } from 'react';
import { Formik, Form, Field } from 'formik';
import { 
  Paper, 
  TextInput, 
  Button, 
  Title, 
  Text, 
  Stack,
  Box,
  Alert,
} from '@mantine/core';
import { IconMail, IconArrowLeft } from '@tabler/icons-react';
import { notifications } from '@mantine/notifications';
import { useForgotPassword } from '../services/authQueries';
import { forgotPasswordSchema } from '../validation/authSchema';

interface ForgotPasswordFormProps {
  onBack?: () => void;
}

export const ForgotPasswordForm = ({ onBack }: ForgotPasswordFormProps) => {
  const forgotPassword = useForgotPassword();
  const [emailSent, setEmailSent] = useState(false);
  const [sentEmail, setSentEmail] = useState('');

  const initialValues = {
    email: '',
  };

  const handlePasswordReset = async (values: typeof initialValues) => {
    try {
      await forgotPassword.mutateAsync(values.email);
      setSentEmail(values.email);
      setEmailSent(true);
      
      notifications.show({
        title: 'Reset Email Sent',
        message: 'Please check your email for password reset instructions.',
        color: 'green',
      });
    } catch (error: any) {
      notifications.show({
        title: 'Reset Failed',
        message: error.message || 'Unable to send reset email. Please try again.',
        color: 'red',
      });
    }
  };

  if (emailSent) {
    return (
      <Paper radius="md" p="xl" withBorder shadow="sm">
        <Stack align="center" gap="md">
          <IconMail size={48} color="var(--mantine-color-blue-6)" />
          
          <Title order={2} ta="center">
            Check Your Email
          </Title>

          <Alert color="green" variant="light" style={{ width: '100%' }}>
            <Text ta="center">
              We've sent password reset instructions to:
            </Text>
            <Text fw={600} ta="center" mt="xs">
              {sentEmail}
            </Text>
          </Alert>

          <Text c="dimmed" size="sm" ta="center">
            Didn't receive the email? Check your spam folder or try again.
          </Text>

          <Stack gap="sm" style={{ width: '100%' }}>
            <Button 
              variant="light" 
              fullWidth
              leftSection={<IconArrowLeft size="1rem" />}
              onClick={() => {
                setEmailSent(false);
                setSentEmail('');
              }}
            >
              Try Different Email
            </Button>
            
            {onBack && (
              <Button 
                variant="outline" 
                fullWidth
                onClick={onBack}
              >
                Back to Login
              </Button>
            )}
          </Stack>
        </Stack>
      </Paper>
    );
  }

  return (
    <Paper radius="md" p="xl" withBorder shadow="sm">
      <Title order={2} ta="center" mb="md">
        Reset Your Password
      </Title>

      <Text c="dimmed" size="sm" ta="center" mb="lg">
        Enter your email address and we'll send you instructions to reset your password.
      </Text>

      <Formik
        initialValues={initialValues}
        validationSchema={forgotPasswordSchema}
        onSubmit={handlePasswordReset}
      >
        {({ errors, touched, isSubmitting, setFieldValue, setFieldTouched }) => (
          <Form>
            <Stack gap="md">
              <Box>
                <Field name="email">
                  {({ field }: any) => (
                    <TextInput
                      {...field}
                      label="Email Address"
                      placeholder="Enter your email"
                      size="md"
                      leftSection={<IconMail size="1rem" />}
                      error={touched.email && errors.email ? errors.email : null}
                      onChange={(e) => {
                        setFieldValue('email', e.target.value);
                        setFieldTouched('email', true);
                      }}
                    />
                  )}
                </Field>
              </Box>

              <Button 
                type="submit" 
                size="md"
                loading={forgotPassword.isPending || isSubmitting}
                fullWidth
              >
                Send Reset Instructions
              </Button>

              {onBack && (
                <Button 
                  variant="subtle" 
                  size="md"
                  fullWidth
                  leftSection={<IconArrowLeft size="1rem" />}
                  onClick={onBack}
                  type="button"
                >
                  Back to Login
                </Button>
              )}
            </Stack>
          </Form>
        )}
      </Formik>
    </Paper>
  );
};
