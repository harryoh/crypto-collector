import request from '@/utils/request'

export function getPrice () {
  return request({
    url: '/prices',
    method: 'get'
  })
}

export function updateAlarm (data) {
  return request({
    url: '/rule',
    method: 'post',
    data
  })
}
