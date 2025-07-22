import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import SignupForm from './index';
import { renderWithStore } from '@/test/utils/renderWithStore';
import { describe, expect, it } from 'vitest';

describe('SignupForm', () => {
  it('renders all form fields', () => {
    renderWithStore(<SignupForm />);

    expect(screen.getByTestId('name-input')).toBeInTheDocument();
    expect(screen.getByTestId('email-input')).toBeInTheDocument();
    expect(screen.getByTestId('password-input')).toBeInTheDocument();
    expect(screen.getByTestId('user-type-select')).toBeInTheDocument();
    expect(screen.getByTestId('signup-button')).toBeInTheDocument();
  });

  it('shows validation errors on submit without values', async () => {
    const user = userEvent.setup();
    renderWithStore(<SignupForm />);
    
    const submitButton = screen.getByTestId('signup-button');
    await user.click(submitButton);

    expect(await screen.findAllByText(/Required/i)).toHaveLength(3);
  });

  it('updates form values when user types', async () => {
    const user = userEvent.setup();
    renderWithStore(<SignupForm />);

    const nameInput = screen.getByTestId('name-input');
    const emailInput = screen.getByTestId('email-input');
    const passwordInput = screen.getByTestId('password-input');
    const userTypeSelect = screen.getByTestId('user-type-select');

    await user.type(nameInput, 'John Doe');
    await user.type(emailInput, 'john@example.com');
    await user.type(passwordInput, 'password123');
    
    await user.click(userTypeSelect);
    await user.click(screen.getByText('Team'));

    expect(nameInput).toHaveValue('John Doe');
    expect(emailInput).toHaveValue('john@example.com');
    expect(passwordInput).toHaveValue('password123');
    expect(userTypeSelect).toHaveValue('Team');
  });
});