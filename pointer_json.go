package minipointer

import (
	"time"

	"github.com/firefirestyle/go.miniprop"
)

func (obj *Pointer) ToJson() []byte {
	propObj := miniprop.NewMiniProp()
	propObj.SetString(PropNameRootGroup, obj.gaeObj.RootGroup)
	propObj.SetString(PropNamePointerId, obj.gaeObj.PointerId)
	propObj.SetString(PropNamePointerType, obj.gaeObj.PointerType)
	propObj.SetString(PropNameOwner, obj.gaeObj.Owner)
	propObj.SetInt(PropNamePoint, obj.gaeObj.Point)

	propObj.SetString(PropNameValue, obj.gaeObj.Value)
	propObj.SetString(PropNameInfo, obj.gaeObj.Info)
	propObj.SetTime(PropNameUpdate, obj.gaeObj.Update)
	propObj.SetString(PropNameSign, obj.gaeObj.Sign)
	return propObj.ToJson()
}

func (obj *Pointer) SetValueFromJson(data []byte) {
	propObj := miniprop.NewMiniPropFromJson(data)
	obj.gaeObj.RootGroup = propObj.GetString(PropNameRootGroup, "")
	obj.gaeObj.PointerId = propObj.GetString(PropNamePointerId, "")
	obj.gaeObj.PointerType = propObj.GetString(PropNamePointerType, "")
	obj.gaeObj.Owner = propObj.GetString(PropNameOwner, "")
	obj.gaeObj.Value = propObj.GetString(PropNameValue, "")
	obj.gaeObj.Point = propObj.GetInt(PropNamePoint, 0)
	obj.gaeObj.Info = propObj.GetString(PropNameInfo, "")
	obj.gaeObj.Update = propObj.GetTime(PropNameUpdate, time.Now())
	obj.gaeObj.Sign = propObj.GetString(PropNameSign, "")
}
