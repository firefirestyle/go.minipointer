package minipointer

import (
	"time"

	"github.com/firefirestyle/go.miniprop"
)

func (obj *Pointer) ToJson() []byte {
	propObj := miniprop.NewMiniProp()
	propObj.SetString(TypeRootGroup, obj.gaeObj.RootGroup)
	propObj.SetString(TypePointerName, obj.gaeObj.PointerName)
	propObj.SetString(TypePointerId, obj.gaeObj.PointerId)
	propObj.SetString(TypePointerType, obj.gaeObj.PointerType)
	propObj.SetString(TypeOwner, obj.gaeObj.Owner)

	propObj.SetString(TypeValue, obj.gaeObj.Value)
	propObj.SetString(TypeInfo, obj.gaeObj.Info)
	propObj.SetTime(TypeUpdate, obj.gaeObj.Update)
	propObj.SetString(TypeSign, obj.gaeObj.Sign)
	return propObj.ToJson()
}

func (obj *Pointer) SetValueFromJson(data []byte) {
	propObj := miniprop.NewMiniPropFromJson(data)
	obj.gaeObj.RootGroup = propObj.GetString(TypeRootGroup, "")
	obj.gaeObj.PointerName = propObj.GetString(TypePointerName, "")
	obj.gaeObj.PointerId = propObj.GetString(TypePointerId, "")
	obj.gaeObj.PointerType = propObj.GetString(TypePointerType, "")
	obj.gaeObj.Owner = propObj.GetString(TypeOwner, "")
	obj.gaeObj.Value = propObj.GetString(TypeValue, "")
	obj.gaeObj.Info = propObj.GetString(TypeInfo, "")
	obj.gaeObj.Update = propObj.GetTime(TypeUpdate, time.Now())
	obj.gaeObj.Sign = propObj.GetString(TypeSign, "")
}
