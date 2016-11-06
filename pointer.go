package minipointer

import (
	"time"

	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

type Pointer struct {
	gaeObj *GaePointerItem
	gaeKey *datastore.Key
	kind   string
}

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

func (obj *Pointer) UpdateMemcache(ctx context.Context) {
	userObjMemSource := obj.ToJson()
	userObjMem := &memcache.Item{
		Key:   obj.gaeKey.StringID(),
		Value: []byte(userObjMemSource), //
	}
	memcache.Set(ctx, userObjMem)
}

func (obj *Pointer) GetName() string {
	return obj.gaeObj.PointerName
}

func (obj *Pointer) GetId() string {
	return obj.gaeObj.PointerId
}
func (obj *Pointer) GetType() string {
	return obj.gaeObj.PointerType
}

func (obj *Pointer) GetValue() string {
	return obj.gaeObj.Value
}

func (obj *Pointer) SetValue(v string) {
	obj.gaeObj.Value = v
}

func (obj *Pointer) GetSign() string {
	return obj.gaeObj.Sign
}

func (obj *Pointer) SetSign(v string) {
	obj.gaeObj.Sign = v
}

func (obj *Pointer) GetOwner() string {
	return obj.gaeObj.Owner
}

func (obj *Pointer) SetOwner(v string) {
	obj.gaeObj.Owner = v
}

func (obj *Pointer) GetInfo() string {
	return obj.gaeObj.Info
}

func (obj *Pointer) GetUpdate() time.Time {
	return obj.gaeObj.Update
}

func (obj *Pointer) Save(ctx context.Context) error {
	Debug(ctx, "SAVE::"+string(obj.ToJson()))
	_, err := datastore.Put(ctx, obj.gaeKey, obj.gaeObj)
	if err == nil {
		obj.UpdateMemcache(ctx)
	}
	return err
}
