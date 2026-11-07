package api

import (
	"errors"

	"FZU-DB-ALPortManager/biz/domain/ship"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"FZU-DB-ALPortManager/biz/model/api"
)

// statusOf 把领域错误翻成 HTTP 状态码。
//
// 原来所有错误一律 500 + 原始 error 文本，前端分不清"没找到"和"服务器炸了"。
// 这里加一张小表，handler 只需要调 mapError。
func statusOf(err error) int {
	switch {
	case errors.Is(err, ship.ErrShipNotFound), errors.Is(err, ship.ErrEquipmentNotFound):
		return consts.StatusNotFound
	case errors.Is(err, ship.ErrSlotOutOfRange),
		errors.Is(err, ship.ErrEquipmentTypeNotAllowed),
		errors.Is(err, ship.ErrAugmentLocked),
		errors.Is(err, ship.ErrEmptyShipName),
		errors.Is(err, ship.ErrUnknownShipType),
		errors.Is(err, ship.ErrEquipmentTypeUnknown):
		return consts.StatusBadRequest
	default:
		return consts.StatusInternalServerError
	}
}

// mapError 输出统一错误体。
func mapError(c *app.RequestContext, err error) {
	status := statusOf(err)
	c.JSON(status, &api.BaseResp{
		Code: int32(status),
		Msg:  err.Error(),
	})
}
