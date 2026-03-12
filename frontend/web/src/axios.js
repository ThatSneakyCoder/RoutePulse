import axios from "axios";

const instance = axios.create({
  baseURL: import.meta.env.VITE_BACKEND_API_URL,
  timeout: 10000,
  headers: {"Content-Type": "application/json"}
});

const token = localStorage.getItem("JWT");
if (token) {
  instance.defaults.headers.common.Authorization = `Bearer ${token}`;
}

// instance.interceptors.request.use(
//   (config) => {
//     const token = localStorage.getItem("token");

//     if (token) {
//       config.headers.Authorization = `Bearer ${token}`;
//     }

//     return config;
//   },
//   (error) => Promise.reject(error)
// );

export default instance;
