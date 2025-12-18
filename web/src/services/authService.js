import axios from 'axios'

const API_BASE_URL = '/api'

export const loginUser = async (username) => {
  try {
    const response = await axios.post(`${API_BASE_URL}/user/login`, {
      username: username
    })

    if (response.data.success) {
      return {
        username: username,
        uid: response.data.uid,
        message: response.data.message
      }
    } else {
      throw new Error(response.data.message || 'Login failed')
    }
  } catch (error) {
    if (error.response) {
      // 服务器响应错误
      throw new Error(error.response.data.message || 'Server error')
    } else if (error.request) {
      // 网络错误
      throw new Error('Network error. Please check your connection.')
    } else {
      // 其他错误
      throw new Error(error.message || 'Login failed')
    }
  }
}
