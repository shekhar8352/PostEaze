import { BaseService } from '@/services/base/BaseService';
import apiClient from '@/services/api/client';
import { type ApiResponse } from '@/services/api/types';
import {type User,type LoginRequest, type RegisterRequest } from '../types';

class AuthService extends BaseService {
  constructor() {
    super('/auth');
  }

  // Auth-specific methods that don't follow CRUD pattern
  async login(data: LoginRequest): Promise<{ user: User; access_token: string; refresh_token: string }> {
    const response = await apiClient.post<ApiResponse<any>>(`${this.endpoint}/login`, data);
    return response.data.data;
  }

  async register(data: RegisterRequest): Promise<{ user: User; access_token: string; refresh_token: string }> {
    const response = await apiClient.post<ApiResponse<any>>(`${this.endpoint}/register`, data);
    return response.data.data;
  }

  async logout(): Promise<void> {
    await apiClient.post(`${this.endpoint}/logout`);
  }

  async getCurrentUser(): Promise<User> {
    const response = await apiClient.get<ApiResponse<User>>(`${this.endpoint}/me`);
    return response.data.data;
  }

  async forgotPassword(email: string): Promise<void> {
    await apiClient.post(`${this.endpoint}/forgot-password`, { email });
  }

  async resetPassword(token: string, password: string): Promise<void> {
    await apiClient.post(`${this.endpoint}/reset-password`, { token, password });
  }

  // Business logic methods
  async loginUser(credentials: LoginRequest) {
    const response = await this.login(credentials);
    
    // Store tokens
    localStorage.setItem('auth_token', response.access_token);
    localStorage.setItem('refresh_token', response.refresh_token);
    localStorage.setItem('user', JSON.stringify(response.user));
    
    return response;
  }

  async registerUser(userData: RegisterRequest) {
    const response = await this.register(userData);
    
    // Store tokens
    localStorage.setItem('auth_token', response.access_token);
    localStorage.setItem('refresh_token', response.refresh_token);
    localStorage.setItem('user', JSON.stringify(response.user));
    
    return response;
  }

  async logoutUser() {
    try {
      await this.logout();
    } finally {
      // Always clear local data
      localStorage.clear();
    }
  }

  isAuthenticated(): boolean {
    return !!localStorage.getItem('auth_token');
  }
}

export const authService = new AuthService();