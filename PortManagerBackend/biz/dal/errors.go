package dal

import "errors"

// ErrNotFound 数据访问层的「查不到」。
//
// 服务层把它翻成对应的领域错误（ErrShipNotFound / ErrEquipmentNotFound），
// 这样 service 不用 import gorm，也不用处理 gorm.ErrRecordNotFound。
var ErrNotFound = errors.New("记录不存在")
