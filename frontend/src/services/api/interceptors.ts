import axios, { type AxiosResponse, type AxiosError, type InternalAxiosRequestConfig } from 'axios';
import apiClient from './client';

// Request interceptor - Add auth token
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('auth_token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: AxiosError) => Promise.reject(error)
);

// Response interceptor - Handle token refresh
apiClient.interceptors.response.use(
  (response: AxiosResponse) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    // Only attempt refresh for 401 errors and if not already retried
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = localStorage.getItem('refresh_token');
        
        if (!refreshToken) {
          throw new Error('No refresh token available');
        }

        // Call refresh endpoint
        const response = await axios.post(`${apiClient.defaults.baseURL}/v1/auth/refresh`, {
          refresh_token: refreshToken,
        });

        const newAccessToken = response.data?.data?.access_token || response.data?.access_token;
        const newRefreshToken = response.data?.data?.refresh_token || response.data?.refresh_token;
        
        if (!newAccessToken) {
          throw new Error('No access token in refresh response');
        }

        // Update both tokens (rotation)
        localStorage.setItem('auth_token', newAccessToken);
        if (newRefreshToken) {
          localStorage.setItem('refresh_token', newRefreshToken);
        }

        // Update authorization header and retry
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
        }

        return apiClient(originalRequest);

      } catch (refreshError) {
        // Refresh failed - clear auth and redirect
        console.error('Token refresh failed:', refreshError);
        
        localStorage.removeItem('auth_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user');
        
        if (!window.location.pathname.includes('/login')) {
          window.location.href = '/login';
        }
        
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

export default apiClient;