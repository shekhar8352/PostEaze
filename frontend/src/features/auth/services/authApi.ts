import axiosInstance from '@/services/axios';

export const signup = async (data: { name: string; email: string; password: string }) => {
  const response = await axiosInstance.post('/auth/signup', data);
  return response.data;
};