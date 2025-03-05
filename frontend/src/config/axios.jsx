import axios from "axios";

const API_BASE_URL = import.meta.env.VITE_DEVELOPMENT_API;

const axiosInstance = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
});

export default axiosInstance;