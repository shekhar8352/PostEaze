import apiClient from '../api/client';
import {type ApiResponse,type PaginatedResponse } from '../api/types';

export class BaseService {
  protected endpoint: string;

  constructor(endpoint: string) {
    this.endpoint = endpoint;
  }

  async get<T>(id: string): Promise<T> {
    const response = await apiClient.get<ApiResponse<T>>(`${this.endpoint}/${id}`);
    return response.data.data;
  }

  async getAll<T>(params?: Record<string, any>): Promise<PaginatedResponse<T>> {
    const response = await apiClient.get<PaginatedResponse<T>>(`${this.endpoint}`, { params });
    return response.data;
  }

  async create<T>(data: any): Promise<T> {
    const response = await apiClient.post<ApiResponse<T>>(`${this.endpoint}`, data);
    return response.data.data;
  }

  async update<T>(id: string, data: any): Promise<T> {
    const response = await apiClient.put<ApiResponse<T>>(`${this.endpoint}/${id}`, data);
    return response.data.data;
  }

  async delete(id: string): Promise<void> {
    await apiClient.delete(`${this.endpoint}/${id}`);
  }
}