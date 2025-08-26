// src/features/auth/components/signUp/RegisterForm.tsx
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
  Box,
  Alert,
  Anchor
} from '@mantine/core';
import { IconInfoCircle } from '@tabler/icons-react';
import { notifications } from '@mantine/notifications';
import { 
  useRegister, 
  useGoogleAuth, 
  useFacebookAuth 
} from '../services/authQueries';
import { registerSchema } from '../validation/authSchema';
import type { RegisterFormData } from '../types';

interface RegisterFormProps {
  onToggleMode?: () => void;
  onSuccess?: () => void;
  onEmailSent?: (email: string, password: string, name: string) => void;
}

export const RegisterForm = ({ onToggleMode, onSuccess, onEmailSent }: RegisterFormProps) => {
  const register = useRegister();
  const googleAuth = useGoogleAuth();
  const facebookAuth = useFacebookAuth();

  const initialValues: RegisterFormData = {
    name: '',
    email: '',
    password: '',
    confirmPassword: '',
  };

  const handleRegister = async (values: RegisterFormData) => {
    try {
      const result = await register.mutateAsync(values);
      
      notifications.show({
        title: 'Registration Initiated!',
        message: 'Please check your email and verify your account to complete registration.',
        color: 'blue',
        autoClose: 8000,
      });
      
      // Redirect to email verification screen
      onEmailSent?.(values.email, values.password, values.name);
      
    } catch (error: any) {
      notifications.show({
        title: 'Registration Failed',
        message: error.message || 'Unable to create account. Please try again.',
        color: 'red',
      });
    }
  };

  const handleGoogleRegister = async () => {
    try {
      await googleAuth.mutateAsync();
      
      notifications.show({
        title: 'Account Created!',
        message: 'Successfully created account with Google.',
        color: 'green',
      });
      
      onSuccess?.();
    } catch (error: any) {
      notifications.show({
        title: 'Google Registration Failed',
        message: error.message || 'Unable to create account with Google.',
        color: 'red',
      });
    }
  };

  const handleFacebookRegister = async () => {
    try {
      await facebookAuth.mutateAsync();
      
      notifications.show({
        title: 'Account Created!',
        message: 'Successfully created account with Facebook.',
        color: 'green',
      });
      
      onSuccess?.();
    } catch (error: any) {
      notifications.show({
        title: 'Facebook Registration Failed',
        message: error.message || 'Unable to create account with Facebook.',
        color: 'red',
      });
    }
  };

  return (
    <Paper radius="md" p="xl" withBorder shadow="sm">
      <Title order={2} ta="center" mb="lg">
        Join PostEaze Today
      </Title>

      <Formik
        initialValues={initialValues}
        validationSchema={registerSchema}
        onSubmit={handleRegister}
      >
        {({ errors, touched, isSubmitting, setFieldValue, setFieldTouched }) => (
          <Form>
            <Stack gap="md">
              <Box>
                <Field name="name">
                  {({ field }: any) => (
                    <TextInput
                      {...field}
                      label="Full Name"
                      placeholder="Enter your full name"
                      size="md"
                      error={touched.name && errors.name ? errors.name : null}
                      onChange={(e) => {
                        setFieldValue('name', e.target.value);
                        setFieldTouched('name', true);
                      }}
                    />
                  )}
                </Field>
              </Box>

              <Box>
                <Field name="email">
                  {({ field }: any) => (
                    <TextInput
                      {...field}
                      label="Email Address"
                      placeholder="Enter your email"
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
                      placeholder="Create a password"
                      size="md"
                      description="Must contain uppercase, lowercase, and number"
                      error={touched.password && errors.password ? errors.password : null}
                      onChange={(e) => {
                        setFieldValue('password', e.target.value);
                        setFieldTouched('password', true);
                      }}
                    />
                  )}
                </Field>
              </Box>

              <Box>
                <Field name="confirmPassword">
                  {({ field }: any) => (
                    <PasswordInput
                      {...field}
                      label="Confirm Password"
                      placeholder="Confirm your password"
                      size="md"
                      error={touched.confirmPassword && errors.confirmPassword ? errors.confirmPassword : null}
                      onChange={(e) => {
                        setFieldValue('confirmPassword', e.target.value);
                        setFieldTouched('confirmPassword', true);
                      }}
                    />
                  )}
                </Field>
              </Box>

              <Alert variant="light" color="orange" icon={<IconInfoCircle size="1rem" />}>
                <Text size="sm">
                  <strong>Important:</strong> After clicking "Create Account", we'll send you a verification email. 
                  You must verify your email before you can log in to PostEaze.
                </Text>
              </Alert>

              <Button 
                type="submit" 
                size="md"
                loading={register.isPending || isSubmitting}
                fullWidth
              >
                Create Account & Send Verification Email
              </Button>

              <Divider label="Or continue with" labelPosition="center" my="lg" />

              <Group grow>
                <Button
                  variant="default"
                  size="md"
                  loading={googleAuth.isPending}
                  onClick={handleGoogleRegister}
                  type="button"
                  leftSection={<span>🔍</span>}
                >
                  Google
                </Button>
                <Button
                  variant="default"
                  size="md"
                  loading={facebookAuth.isPending}
                  onClick={handleFacebookRegister}
                  type="button"
                  leftSection={<span>📘</span>}
                >
                  Facebook
                </Button>
              </Group>

              {onToggleMode && (
                <Text c="dimmed" size="sm" ta="center" mt="md">
                  Already have an account?{' '}
                  <Anchor 
                    component="button" 
                    onClick={onToggleMode}
                    type="button"
                  >
                    Sign in here
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
