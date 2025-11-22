// Auth-related constants
export const AUTH_STORAGE_KEYS = {
  AUTH_TOKEN: 'auth_token',
  REFRESH_TOKEN: 'refresh_token',
  USER: 'user',
} as const;

export const AUTH_ROUTES = {
  LOGIN: '/login',
  REGISTER: '/register',
  FORGOT_PASSWORD: '/forgot-password',
  EMAIL_VERIFY: '/email-verify',
  DASHBOARD: '/dashboard',
  HOME: '/',
} as const;

export const AUTH_ERRORS = {
  INVALID_CREDENTIALS: 'Invalid email or password',
  EMAIL_NOT_VERIFIED: 'Please verify your email before logging in',
  SESSION_EXPIRED: 'Your session has expired. Please login again',
  NETWORK_ERROR: 'Network error. Please check your connection',
  UNKNOWN_ERROR: 'An unexpected error occurred',
} as const;

export const AUTH_MESSAGES = {
  LOGIN_SUCCESS: 'Welcome back!',
  LOGOUT_SUCCESS: 'You have been logged out',
  REGISTER_SUCCESS: 'Registration successful! Please verify your email',
  PASSWORD_RESET_SENT: 'Password reset link sent to your email',
} as const;
