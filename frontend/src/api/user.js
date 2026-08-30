import request from '@/utils/request'

export function login (data) {
  return request({
    url: '/login',
    method: 'post',
    data,
  })
}

export function current () {
  return request({
    url: '/user/current',
    method: 'get',
  })
}

export function list (params) {
  return request({
    url: '/user/list',
    params,
  })
}

export function detail (id) {
  return request({
    url: `/user/detail/${id}`,
  })
}

export function create (data) {
  return request({
    url: '/user/create',
    method: 'post',
    data,
  })
}

export function update (data) {
  return request({
    url: '/user/update',
    method: 'post',
    data,
  })
}

export function remove (data) {
  return request({
    url: '/user/delete',
    method: 'post',
    data,
  })
}

export function changePwd (data) {
  return request({
    url: '/user/changePwd',
    method: 'post',
    data,
  })
}

export function changeCurPwd (data) {
  return request({
    url: '/user/changeCurPwd',
    method: 'post',
    data,
  })
}

export function otpStatus () {
  return request({ url: '/user/otp/status' })
}

export function otpSetup () {
  return request({ url: '/user/otp/setup', method: 'post' })
}

export function otpConfirm (code) {
  return request({ url: '/user/otp/confirm', method: 'post', data: { code } })
}

export function otpDisable (code) {
  return request({ url: '/user/otp/disable', method: 'post', data: { code } })
}

export function otpReset (id) {
  return request({ url: '/user/otp/reset', method: 'post', data: { id } })
}

export function myOauth () {
  return request({
    url: '/user/myOauth',
    method: 'post',
  })
}

export function myPeer (params) {
  return request({
    url: '/user/myPeer',
    params,
  })
}

export function groupUsers (data) {
  return request({
    url: '/user/groupUsers',
    method: 'post',
    data,
  })
}

export function register (data) {
  return request({
    url: '/user/register',
    method: 'post',
    data,
  })
}
