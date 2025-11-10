import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as shipApi from '@/api/ship'

export const useShipStore = defineStore('ship', () => {
  const ships = ref([])
  const total = ref(0)
  const loading = ref(false)

  // 获取舰船列表
  const fetchShipList = async (params) => {
    loading.value = true
    try {
      const data = await shipApi.getShipList(params)
      ships.value = data.ships || []
      total.value = data.total || 0
      return data
    } finally {
      loading.value = false
    }
  }

  // 创建舰船
  const createShip = async (shipData) => {
    return await shipApi.createShip(shipData)
  }

  // 更新舰船
  const updateShip = async (id, shipData) => {
    return await shipApi.updateShip(id, shipData)
  }

  // 删除舰船
  const deleteShip = async (id) => {
    return await shipApi.deleteShip(id)
  }

  return {
    ships,
    total,
    loading,
    fetchShipList,
    createShip,
    updateShip,
    deleteShip
  }
})
