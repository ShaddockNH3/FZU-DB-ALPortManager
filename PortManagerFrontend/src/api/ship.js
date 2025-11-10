import request from './index'

/**
 * 创建舰船
 * @param {Object} data - 舰船数据
 */
export const createShip = (data) => {
  return request({
    url: '/ship',
    method: 'POST',
    data
  })
}

/**
 * 获取舰船详情
 * @param {Number} id - 舰船ID
 */
export const getShipById = (id) => {
  return request({
    url: `/ship/${id}`,
    method: 'GET'
  })
}

/**
 * 获取舰船列表
 * @param {Object} params - 查询参数
 */
export const getShipList = (params) => {
  return request({
    url: '/ships',
    method: 'GET',
    params
  })
}

/**
 * 更新舰船
 * @param {Number} id - 舰船ID
 * @param {Object} data - 舰船数据
 */
export const updateShip = (id, data) => {
  return request({
    url: `/ship/${id}`,
    method: 'PUT',
    data: { ship: data }
  })
}

/**
 * 删除舰船
 * @param {Number} id - 舰船ID
 */
export const deleteShip = (id) => {
  return request({
    url: `/ship/${id}`,
    method: 'DELETE'
  })
}

/**
 * 获取统计数据
 */
export const getStatistics = () => {
  return request({
    url: '/stats',
    method: 'GET'
  })
}
