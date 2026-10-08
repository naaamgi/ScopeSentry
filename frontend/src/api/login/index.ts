import request from '@/axios'
import type { UserType, Token, changePassword, changePasswordResponse } from './types'

interface RoleParams {
  roleName: string
}

export const loginApi = (data: UserType): Promise<IResponse<Token>> => {
  return request.post({ url: '/api/user/login', data })
}

export const setupStatusApi = (): Promise<IResponse<{ required: boolean }>> => {
  return request.get({ url: '/api/user/setup/status' })
}

export const setupApi = (data: { username: string; password: string }): Promise<IResponse> => {
  return request.post({ url: '/api/user/setup', data })
}

export const changePasswordApi = (
  data: changePassword
): Promise<IResponse<changePasswordResponse>> => {
  return request.post({ url: '/api/user/changePassword', data })
}
export const loginOutApi = (): Promise<IResponse> => {
  return request.get({ url: '/mock/user/loginOut' })
}

export const getUserListApi = ({ params }: AxiosConfig) => {
  return request.get<{
    code: string
    data: {
      list: UserType[]
      total: number
    }
  }>({ url: '/mock/user/list', params })
}

export const getAdminRoleApi = (
  params: RoleParams
): Promise<IResponse<AppCustomRouteRecordRaw[]>> => {
  return request.get({ url: '/mock/role/list', params })
}

export const getTestRoleApi = (params: RoleParams): Promise<IResponse<string[]>> => {
  return request.get({ url: '/mock/role/list2', params })
}
