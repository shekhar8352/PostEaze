// index.test.tsx
import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import SignupForm from './index';
import { renderWithStore } from '@/test/utils/renderWithStore';
import { describe, expect, it } from 'vitest';

describe('SignupForm', () => {
  it('renders all form fields', () => {
    renderWithStore(<SignupForm />);

    expect(screen.getByLabelText(/Name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Password/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/User Type/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Sign Up/i })).toBeInTheDocument();
  });

  it('shows validation errors on submit without values', async () => {
    renderWithStore(<SignupForm />);
    await userEvent.click(screen.getByRole('button', { name: /Sign Up/i }));

    expect(await screen.findAllByText(/Required/i)).toHaveLength(3);
  });
});
