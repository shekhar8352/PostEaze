// User related types
export interface User {
  id: string;
  email: string;
  name: string;
  role: 'admin' | 'user' | 'moderator';
  isActive: boolean;
  avatar?: string;
  createdAt: string;
  updatedAt: string;
}

// Auth request types
export interface LoginRequest {
  email: string;
  password?: string;
  rememberMe?: boolean;
  // Firebase fields (only for verified users)
  firebase_uid: string; // Required - only verified users reach backend
  firebase_token: string; // Required - only verified users reach backend
  display_name?: string;
  email_verified: true; // Always true - only verified users reach backend
  provider: 'email' | 'google.com' | 'facebook.com';
}

export interface RegisterRequest {
  name: string;
  email: string;
  password?: string;
  confirmPassword?: string;
  terms: boolean;
  // Firebase fields (only for verified users)
  firebase_uid: string; // Required - only verified users reach backend
  firebase_token: string; // Required - only verified users reach backend
  display_name: string; // Required for backend user creation
  email_verified: true; // Always true - only verified users reach backend
  provider: 'email' | 'google.com' | 'facebook.com';
}

export interface ForgotPasswordRequest {
  email: string;
}

export interface ResetPasswordRequest {
  token: string;
  password: string;
  confirmPassword: string;
}

export interface ChangePasswordRequest {
  currentPassword: string;
  newPassword: string;
  confirmPassword: string;
}

// Auth response types
export interface AuthResponse {
  user: User;
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface RefreshTokenResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

// Auth state types (for Redux)
export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  loading: boolean;
  error: string | null;
  tokenExpiry: number | null;
}

// Form types (for react-hook-form)
export interface LoginFormData {
  email: string;
  password: string;
}

export interface RegisterFormData {
  name: string;
  email: string;
  password: string;
  confirmPassword: string;
}

export interface ForgotPasswordFormData {
  email: string;
}

export interface EmailVerificationState {
  isEmailSent: boolean;
  email: string;
  isResending: boolean;
  canResend: boolean;
  countdown: number;
}

export interface ResetPasswordFormData {
  password: string;
  confirmPassword: string;
}

// API error types
export interface AuthError {
  message: string;
  code: string;
  field?: string;
}

// Permission types
export type Permission = 'read' | 'write' | 'delete' | 'admin';

export interface UserPermissions {
  userId: string;
  permissions: Permission[];
  role: string;
}