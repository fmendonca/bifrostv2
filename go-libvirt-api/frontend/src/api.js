import axios from 'axios';

const API = axios.create({
    baseURL: 'http://localhost:8080',
});

export const login = (username, password) =>
    API.post('/login', { username, password });

export const getHosts = (token) =>
    API.get('/hosts', { headers: { Authorization: token } });
