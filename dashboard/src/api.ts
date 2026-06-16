import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080',
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  
  const targetTenant = localStorage.getItem('targetTenant');
  if (targetTenant) {
    config.headers['X-Target-Tenant-ID'] = targetTenant;
  }
  
  return config;
});

export default api;
