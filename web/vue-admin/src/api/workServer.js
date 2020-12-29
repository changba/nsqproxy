import request from '@/utils/request'

export function getAll(params) {
  return request({
    url: './api/workServer/all',
    method: 'get'
  })
}

export function getList(params) {
  return request({
    url: './api/workServer/page',
    method: 'get',
    data: {
      ...params
    }
  })
}

export function create(params) {
  return request({
    url: './api/workServer/create',
    method: 'get',
    data: {
      ...params
    }
  })
}

export function update(params) {
  return request({
    url: './api/workServer/update',
    method: 'get',
    data: {
      ...params
    }
  })
}

export function deleteAction(params) {
  return request({
    url: './api/workServer/delete',
    method: 'get',
    data: {
      ...params
    }
  })
}
