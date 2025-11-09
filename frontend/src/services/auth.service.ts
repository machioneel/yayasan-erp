import { api } from './api';
import { LoginRequest, LoginResponse, User, ApiResponse } from '@/types';

export const authService = {
  login: async (credentials: LoginRequest) => {
    const response = await api.post<ApiResponse<LoginResponse>>('/auth/login', credentials);
    return response.data.data;
  },

  logout: async () => {
    await api.post('/auth/logout');
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },

  getCurrentUser: async () => {
    const response = await api.get<ApiResponse<User>>('/users/me');
    return response.data.data;
  },

  refreshToken: async () => {
    const response = await api.post<ApiResponse<{ token: string }>>('/auth/refresh');
    return response.data.data.token;
  },
};
