import { AUTH_STORAGE_KEYS } from './constants';
import type { User } from './types';

/**
 * Storage utility functions for authentication
 */
export const authStorage = {
  // Token management
  getAuthToken: (): string | null => {
    return localStorage.getItem(AUTH_STORAGE_KEYS.AUTH_TOKEN);
  },

  setAuthToken: (token: string): void => {
    localStorage.setItem(AUTH_STORAGE_KEYS.AUTH_TOKEN, token);
  },

  getRefreshToken: (): string | null => {
    return localStorage.getItem(AUTH_STORAGE_KEYS.REFRESH_TOKEN);
  },

  setRefreshToken: (token: string): void => {
    localStorage.setItem(AUTH_STORAGE_KEYS.REFRESH_TOKEN, token);
  },

  // User management
  getUser: (): User | null => {
    const userStr = localStorage.getItem(AUTH_STORAGE_KEYS.USER);
    if (!userStr) return null;
    
    try {
      return JSON.parse(userStr) as User;
    } catch (error) {
      console.error('Failed to parse user data:', error);
      return null;
    }
  },

  setUser: (user: User): void => {
    localStorage.setItem(AUTH_STORAGE_KEYS.USER, JSON.stringify(user));
  },

  // Clear all auth data
  clearAuth: (): void => {
    localStorage.removeItem(AUTH_STORAGE_KEYS.AUTH_TOKEN);
    localStorage.removeItem(AUTH_STORAGE_KEYS.REFRESH_TOKEN);
    localStorage.removeItem(AUTH_STORAGE_KEYS.USER);
  },

  // Check if user is authenticated
  isAuthenticated: (): boolean => {
    return !!authStorage.getAuthToken();
  },
};

/**
 * Validate JWT token (basic check)
 */
export const isTokenValid = (token: string | null): boolean => {
  if (!token) return false;
  
  try {
    // Basic JWT structure check (header.payload.signature)
    const parts = token.split('.');
    if (parts.length !== 3) return false;

    // Decode payload to check expiration
    const payload = JSON.parse(atob(parts[1]));
    if (!payload.exp) return true; // No expiration set
    
    // Check if token is expired
    const expirationTime = payload.exp * 1000; // Convert to milliseconds
    return Date.now() < expirationTime;
  } catch (error) {
    console.error('Token validation error:', error);
    return false;
  }
};

/**
 * Format error messages for display
 */
export const formatAuthError = (error: any): string => {
  if (typeof error === 'string') return error;
  if (error?.message) return error.message;
  if (error?.response?.data?.message) return error.response.data.message;
  return 'An unexpected error occurred';
};

/**
 * Get redirect path after login
 */
export const getRedirectPath = (): string => {
  const params = new URLSearchParams(window.location.search);
  return params.get('redirect') || '/dashboard';
};

/**
 * Create login redirect URL with return path
 */
export const createLoginRedirect = (returnPath?: string): string => {
  const path = returnPath || window.location.pathname;
  if (path === '/login' || path === '/register') return '/login';
  return `/login?redirect=${encodeURIComponent(path)}`;
};
