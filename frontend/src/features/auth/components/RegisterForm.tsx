// src/features/auth/components/RegisterForm.tsx
import { Formik, Form, Field } from 'formik';
import { 
  TextInput, 
  PasswordInput, 
  Button, 
  Text, 
  Stack, 
  Group, 
  Divider,
  Box,
  Alert,
  Anchor
} from '@mantine/core';
import { IconBrandFacebook, IconBrandGoogle, IconInfoCircle } from '@tabler/icons-react';
import authForm from './authForm.module.css';
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
  onEmailSent?: (email: string, password: string, name?: string) => void;
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
      await register.mutateAsync(values);
      
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
    <Box className={authForm.surface}>
      <Text component="h2" className={authForm.title}>
        Create your account
      </Text>
      <Text className={authForm.lede}>
        We'll send one verification email—then you're in.
      </Text>

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

              <Alert variant="light" color="blue" icon={<IconInfoCircle size="1rem" />} className={authForm.alert}>
                <Text size="sm" c="inherit">
                  <strong>Next step:</strong> After you register, check your inbox and verify your email before signing in.
                </Text>
              </Alert>

              <Button 
                type="submit" 
                size="md"
                loading={register.isPending || isSubmitting}
                fullWidth
                variant="filled"
                classNames={{ root: authForm.primaryButton }}
              >
                Create account & send link
              </Button>

              <Divider className={authForm.divider} label="Or continue with" labelPosition="center" my="lg" />

              <Group grow>
                <Button
                  variant="default"
                  size="md"
                  loading={googleAuth.isPending}
                  onClick={handleGoogleRegister}
                  type="button"
                  leftSection={<IconBrandGoogle size={18} />}
                  classNames={{ root: authForm.socialButton }}
                >
                  Google
                </Button>
                <Button
                  variant="default"
                  size="md"
                  loading={facebookAuth.isPending}
                  onClick={handleFacebookRegister}
                  type="button"
                  leftSection={<IconBrandFacebook size={18} />}
                  classNames={{ root: authForm.socialButton }}
                >
                  Facebook
                </Button>
              </Group>

              {onToggleMode && (
                <Text className={authForm.footer} mt="md">
                  Already have an account?{' '}
                  <Anchor 
                    component="button" 
                    onClick={onToggleMode}
                    type="button"
                  >
                    Sign in
                  </Anchor>
                </Text>
              )}
            </Stack>
          </Form>
        )}
      </Formik>
    </Box>
  );
};
