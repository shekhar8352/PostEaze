import { useState } from 'react';
import { Formik, Form, Field } from 'formik';
import { TextInput, Button, Text, Stack, Box, Alert } from '@mantine/core';
import { IconMail, IconArrowLeft } from '@tabler/icons-react';
import { notifications } from '@mantine/notifications';
import { useForgotPassword } from '../services/authQueries';
import { forgotPasswordSchema } from '../validation/authSchema';
import authForm from './authForm.module.css';

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
        title: 'Reset email sent',
        message: 'Check your inbox for password reset instructions.',
        color: 'green',
      });
    } catch (error: any) {
      notifications.show({
        title: 'Request failed',
        message: error.message || 'Unable to send reset email. Please try again.',
        color: 'red',
      });
    }
  };

  if (emailSent) {
    return (
      <Box className={`${authForm.surface} ${authForm.centerStack}`}>
        <Stack align="center" gap="md">
          <IconMail size={44} color="var(--auth-surface-accent)" stroke={1.5} />
          
          <Text component="h2" className={authForm.title}>
            Check your email
          </Text>

          <Alert color="green" variant="light" style={{ width: '100%' }}>
            <Text ta="center" size="sm">
              We sent instructions to:
            </Text>
            <Text fw={600} ta="center" mt="xs" size="sm">
              {sentEmail}
            </Text>
          </Alert>

          <Text c="dimmed" size="sm" ta="center">
            Did not receive it? Check spam or contact your IT team.
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
              Try a different email
            </Button>
            
            {onBack && (
              <Button 
                variant="default" 
                fullWidth
                onClick={onBack}
              >
                Back to sign in
              </Button>
            )}
          </Stack>
        </Stack>
      </Box>
    );
  }

  return (
    <Box className={authForm.surface}>
      <Text component="h2" className={authForm.title}>
        Reset password
      </Text>

      <Text className={authForm.lede}>
        Enter your work email. We will send a secure link that expires automatically.
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
                      label="Email address"
                      placeholder="you@company.com"
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
                variant="filled"
                classNames={{ root: authForm.primaryButton }}
              >
                Send reset link
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
                  Back to sign in
                </Button>
              )}
            </Stack>
          </Form>
        )}
      </Formik>
    </Box>
  );
};
