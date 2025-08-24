import axios, { type AxiosResponse, AxiosError } from 'axios';
import apiClient from './client';
// import { store } from '@/app/store/store';
// import { logout } from '../features/auth/store/authSlice';

// Request interceptor - Add auth token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
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
    const originalRequest = error.config as any;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      try {
        const refreshToken = localStorage.getItem('refresh_token');
        
        if (!refreshToken) {
          throw new Error('No refresh token');
        }

        // Call refresh endpoint
        const response = await axios.post(`${apiClient.defaults.baseURL}/auth/refresh`, {
          refresh_token: refreshToken,
        });

        const { access_token, refresh_token: newRefreshToken } = response.data;
        
        // Update tokens
        localStorage.setItem('auth_token', access_token);
        localStorage.setItem('refresh_token', newRefreshToken);

        // Retry original request
        originalRequest.headers.Authorization = `Bearer ${access_token}`;
        return apiClient(originalRequest);

      } catch (refreshError) {
        // Refresh failed - logout user
        localStorage.clear();
        // store.dispatch(logout());
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

export default apiClient;