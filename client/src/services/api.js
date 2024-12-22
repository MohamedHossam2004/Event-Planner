import axios from 'axios';
import { jwtDecode } from "jwt-decode";

const API_URL = 'http://localhost:8080/v1';

export const login = async (email, password) => {
  try {
    const response = await axios.post(`${API_URL}/login`, { email, password });
    return response.data;
  } catch (error) {
    throw new Error(error.response?.data?.message || 'Login failed');
  }
};

export const signup = async (name, email, password) => {
  try {
    const response = await axios.post(`${API_URL}/register`, { name, email, password });
    return response.data;
  } catch (error) {
    throw new Error(error.response?.data?.message || 'Signup failed');
  }
};

export const logout = async () => {
  // Clear the token from cookies
  document.cookie = 'token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
};

export const verifyAccount = async () => {
  try {
    const response = await axios.post(`${API_URL}/v1/verify`, {}, {
      headers: { Authorization: `Bearer ${getCookie('token')}` }
    });
    return response.data;
  } catch (error) {
    throw new Error(error.response?.data?.message || 'Verification failed');
  }
};

export const createEvent = async (eventData) => {
  try {
    const response = await axios.post(`${API_URL}/events`, eventData, {
      headers: { Authorization: `Bearer ${getCookie('token')}` }
    });
    return response.data;
  } catch (error) {
    throw new Error(error.response?.data?.message || 'Failed to create event');
  }
};

export const getCookie = (name) => {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(';').shift();
};

export const decodeToken = (token) => {
  try {
    return jwtDecode(token);
  } catch (error) {
    console.error('Failed to decode token:', error);
    return null;
  }
};

// Existing getEvents function
export const getEvents = async () => {
  try {
    const response = await axios.get(`${API_URL}/events`);
    return response.data;
  } catch (error) {
    throw new Error(error.response?.data?.message || 'Failed to fetch events');
  }
};


export const getEventsForUser = async()=>{
  try{

    const response = await axios.get(`${API_URL}/events/user`,{
         headers: { Authorization: `Bearer ${getCookie('token')}` }
       });
       return response;


  }
  catch{

    throw new Error(error.response?.data?.message || 'Failed to fetch events for user');

  }
}


export const applyToEvent = async (eventId) => {
  try {
    const token = getCookie('token'); 
    const response = await axios.post(`${API_URL}/events/${eventId}/apply`, {}, {
      headers: {
        Authorization: `Bearer ${token}`, 
      },
    });
    return response.data;
  } catch (error) {
    console.log(error)
    throw new Error(error.response?.data?.error || 'Failed to apply to event');
  }
};
