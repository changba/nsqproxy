import axios from 'axios'
import { Message } from 'element-ui'
import qs from 'qs'

// 创建axios实例
const service = axios.create({
  timeout: 30000,
  headers: {
    'X-Requested-With': 'XMLHttpRequest'
  },
  transformRequest: [function(data, headers) {
    if (data instanceof FormData) {
      return data
    }

    if (data?.__json__ === true) {
      headers['Content-Type'] = 'application/json; charset=utf-8'
      data.__json__ = undefined

      return JSON.stringify(data)
    }
    return qs.stringify(data)
  }],
  data: {
    _timestamp: new Date()
  }
})

service.interceptors.request.use(config => {
  if (config.method === 'get') {
    config.params = config.data
  }

  return config
}, error => {
  // Do something with request error
  console.log(error) // for debug
  Promise.reject(error)
})

service.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code !== 200) {
      Message({
        message: res.msg || 'Error',
        type: 'error',
        duration: 3 * 1000
      })
      return Promise.reject(new Error(res.msg || 'Error'))
    } else {
      return res
    }
  },
  error => {
    console.log('err' + error) // for debug
    Message({
      message: error.message,
      type: 'error',
      duration: 5 * 1000
    })
    return Promise.reject(error)
  }
)

export default service
