import request from '@/utils/request'

export function getList(params) {
  return request({
    url: './api/consumeConfig/page',
    method: 'get',
    data: { ...params }
  })
}

export function create(params) {
  return request({
    url: './api/consumeConfig/create',
    method: 'get',
    data: {
      ...params,
    }
  })
}

export function update(params) {
  return request({
    url: './api/consumeConfig/update',
    method: 'get',
    data: {
      ...params,
    }
  })
}

export function deleteAction(params) {
  return request({
    url: './api/consumeConfig/delete',
    method: 'get',
    data: {
      ...params
    }
  })
}
