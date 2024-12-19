import axios from 'axios';
import * as jwt_decode from 'jwt-decode';

const api = axios.create({
  baseURL: 'http://localhost:8080'
});

export const getEvents = async () => {
  const response = await api.get('/v1/events');
  return response.data.events;
};

export const getEventById = async (id) => {
  const response = await api.get(`/v1/events/${id}`);
  return response.data.event;
};

export const login = async (email, password) => {
  try {
    const response = await api.post(`/v1/login`, { email, password });
    return response.data;
  } catch (error) {
    throw new Error(error.response?.data?.message || 'Login failed');
  }
};

export const signup = async (name, email, password) => {
  try {
    const response = await api.post(`/v1/signup`, { name, email, password });
    return response.data;
  } catch (error) {
    throw new Error(error.response?.data?.message || 'Signup failed');
  }
};

export const verifyAccount = async (email) => {
  try {
    const response = await api.post(`/v1/verify`, {email}, {
      headers: { Authorization: `Bearer ${getCookie('token')}` }
    });
    return response.data;
  } catch (error) {
    throw new Error(error.response?.data?.message || 'Verification failed');
  }
};

export const getCookie = (name) => {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(';').shift();
};

export const decodeToken = (token) => {
  try {
    return jwt_decode(token);
  } catch (error) {
    console.error('Failed to decode token:', error);
    return null;
  }
};