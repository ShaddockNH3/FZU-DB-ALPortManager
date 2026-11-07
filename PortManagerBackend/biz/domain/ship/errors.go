package ship

import "errors"

// 领域错误。handler 按类型映射 HTTP 状态码，不再一律 500。
var (
	ErrShipNotFound            = errors.New("舰船不存在")
	ErrEquipmentNotFound       = errors.New("装备不存在")
	ErrUnknownShipType         = errors.New("未知舰种，无法确定装备槽")
	ErrSlotOutOfRange          = errors.New("装备槽位超出范围")
	ErrEquipmentTypeUnknown    = errors.New("未知装备类型")
	ErrEquipmentTypeNotAllowed = errors.New("该槽位不能装备这类装备")
	ErrAugmentLocked           = errors.New("兵装槽未解锁，先把星级补满")
	ErrEmptyShipName           = errors.New("舰船名称不能为空")
)
