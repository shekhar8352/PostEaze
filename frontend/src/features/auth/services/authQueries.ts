import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { authService } from './authService';
import {type LoginRequest,type RegisterRequest } from '../types';


export const useLogin = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: LoginRequest) => authService.loginUser(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['currentUser'] });
    //   toast.success('Login successful!');
    },
    onError: (error: any) => {
      const message = error.response?.data?.message || 'Login failed';
    //   toast.error(message);
    },
  });
};

export const useRegister = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: RegisterRequest) => authService.registerUser(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['currentUser'] });
    //   toast.success('Registration successful!');
    },
    onError: (error: any) => {
      const message = error.response?.data?.message || 'Registration failed';
    //   toast.error(message);
    },
  });
};

export const useLogout = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => authService.logoutUser(),
    onSuccess: () => {
      queryClient.clear();
    //   toast.success('Logged out successfully');
    },
  });
};

export const useCurrentUser = () => {
  return useQuery({
    queryKey: ['currentUser'],
    queryFn: () => authService.getCurrentUser(),
    enabled: authService.isAuthenticated(),
  });
};

export const useForgotPassword = () => {
  return useMutation({
    mutationFn: (email: string) => authService.forgotPassword(email),
    onSuccess: () => {
    //   toast.success('Password reset email sent!');
    },
    onError: (error: any) => {
      const message = error.response?.data?.message || 'Failed to send reset email';
    //   toast.error(message);
    },
  });
};

export const useResetPassword = () => {
  return useMutation({
    mutationFn: ({ token, password }: { token: string; password: string }) => 
      authService.resetPassword(token, password),
    onSuccess: () => {
    //   toast.success('Password reset successful!');
    },
    onError: (error: any) => {
      const message = error.response?.data?.message || 'Failed to reset password';
    //   toast.error(message);
    },
  });
};