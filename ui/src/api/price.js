import request from '@/utils/request'

export function getPrice () {
  return request({
    // url: '/prices',
    url: '',
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
