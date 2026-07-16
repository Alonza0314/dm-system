import axios from 'axios'
import { AccountApi, CategoryApi, DeviceApi, Configuration } from './api'
import { getErrorMessage } from './utils/getErrorMessage'

const apiBasePath = import.meta.env.VITE_API_BASE_URL || `${window.location.protocol}//${window.location.hostname}:8888`

const configuration = new Configuration({
  basePath: apiBasePath,
  accessToken: () => localStorage.getItem('token') || '',
})

export const accountApi = new AccountApi(configuration)
export const categoryApi = new CategoryApi(configuration)
export const deviceApi = new DeviceApi(configuration)

export const UNAUTHORIZED_MESSAGE_KEY = 'dm-system:unauthorized-message'

axios.interceptors.response.use(
  (response) => response,
  (error) => {
    const isUnauthorized = error?.response?.status === 401
    const isOnLoginPage = window.location.pathname === '/login'

    if (isUnauthorized && !isOnLoginPage) {
      localStorage.removeItem('token')
      sessionStorage.setItem(
        UNAUTHORIZED_MESSAGE_KEY,
        getErrorMessage(error, 'Session expired. Please sign in again.'),
      )
      window.location.href = '/login'
    }

    return Promise.reject(error)
  },
)
