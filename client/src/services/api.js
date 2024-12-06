import axios from 'axios';

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