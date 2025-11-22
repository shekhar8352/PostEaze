import { useCallback } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { setUser, clearUser, selectUser, selectIsAuthenticated, selectAuthLoading } from '../authSlice';
import { authStorage } from '../utils';
import type { User } from '../types';

/**
 * Custom hook for authentication operations and state
 */
export const useAuth = () => {
  const dispatch = useDispatch();
  const user = useSelector(selectUser);
  const isAuthenticated = useSelector(selectIsAuthenticated);
  const loading = useSelector(selectAuthLoading);

  /**
   * Initialize auth state from localStorage
   * Call this on app mount to hydrate the store
   */
  const initializeAuth = useCallback(() => {
    const storedUser = authStorage.getUser();
    const token = authStorage.getAuthToken();
    
    if (storedUser && token) {
      dispatch(setUser(storedUser));
    }
  }, [dispatch]);

  /**
   * Set user and persist to storage
   */
  const login = useCallback((userData: User, accessToken: string, refreshToken: string) => {
    // Store in localStorage
    authStorage.setAuthToken(accessToken);
    authStorage.setRefreshToken(refreshToken);
    authStorage.setUser(userData);
    
    // Update Redux store
    dispatch(setUser(userData));
  }, [dispatch]);

  /**
   * Clear user and remove from storage
   */
  const logout = useCallback(() => {
    // Clear localStorage
    authStorage.clearAuth();
    
    // Clear Redux store
    dispatch(clearUser());
  }, [dispatch]);

  /**
   * Update user data
   */
  const updateUser = useCallback((userData: User) => {
    authStorage.setUser(userData);
    dispatch(setUser(userData));
  }, [dispatch]);

  return {
    user,
    isAuthenticated,
    loading,
    login,
    logout,
    updateUser,
    initializeAuth,
  };
};
