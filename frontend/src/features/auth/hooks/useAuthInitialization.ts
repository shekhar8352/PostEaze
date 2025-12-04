import { useEffect } from 'react';
import { useDispatch } from 'react-redux';
import { useQuery } from '@tanstack/react-query';
import { authService } from '../services/authService';
import { setUser, clearUser } from '../authSlice';
import { authStorage } from '../utils';

/**
 * Hook that fetches user data on app startup
 * Restores authentication state from backend using /auth/me
 */
export const useAuthInitialization = () => {
  const dispatch = useDispatch();

  const { data: user, error, isLoading } = useQuery({
    queryKey: ['auth', 'me'],
    queryFn: async () => {
      const token = authStorage.getAuthToken();
      if (!token) {
        throw new Error('No authentication token');
      }
      return authService.getCurrentUser();
    },
    retry: false,
    staleTime: 5 * 60 * 1000, // 5 minutes
    enabled: !!authStorage.getAuthToken(),
  });

  useEffect(() => {
    if (user) {
      dispatch(setUser(user));
    } else if (error) {
      console.error('Failed to fetch user:', error);
      dispatch(clearUser());
      authStorage.clearAuth();
    }
  }, [user, error, dispatch]);

  return {
    isLoading,
    isAuthenticated: !!user,
    user,
  };
};

