import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
})

// Har request me JWT token automatically attach hoga
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 401 aaye toh auto logout
api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(err)
  }
)

export const authAPI = {
  register: (data) => api.post('/auth/register', data),
  login:    (data) => api.post('/auth/login', data),
  me:       ()     => api.get('/auth/me'),
}

export const marketAPI = {
  search:  (q)              => api.get(`/market/search?q=${encodeURIComponent(q)}`),
  quote:   (symbol)         => api.get(`/market/quote/${symbol}`),
  history: (symbol, period) => api.get(`/market/history/${symbol}?period=${period}`),
}