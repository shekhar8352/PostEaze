// src/features/auth/pages/SignupPage.tsx
import { Paper, Title } from '@mantine/core';
import SignupForm from '../components/SignupForm';

const SignupPage = () => {
  return (
    <Paper radius="md" withBorder>
      <Title order={2} mb="md">
        Sign Up
      </Title>
      <SignupForm />
    </Paper>
  );
};

export default SignupPage;
