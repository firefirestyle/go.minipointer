package minipointer

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

const (
	PropNameRootGroup   = "RootGroup"
	PropNamePointerId   = "IdentifyId"
	PropNamePointerType = "PointerType"
	PropNameValue       = "UserName"
	PropNameInfo        = "Info"
	PropNameUpdate      = "Update"
	PropNameSign        = "Sign"
	PropNameOwner       = "Owner"
	PropNamePoint       = "Point"
)

type GaePointerItem struct {
	RootGroup   string
	PointerId   string
	PointerType string
	Value       string
	Info        string
	Update      time.Time
	Owner       string
	Sign        string
	Point       int
}

type Pointer struct {
	gaeObj *GaePointerItem
	gaeKey *datastore.Key
	kind   string
}

func (obj *Pointer) UpdateMemcache(ctx context.Context) {
	userObjMemSource := obj.ToJson()
	userObjMem := &memcache.Item{
		Key:   obj.gaeKey.StringID(),
		Value: []byte(userObjMemSource), //
	}
	memcache.Set(ctx, userObjMem)
}
func (obj *PointerManager) DeleteMemcache(ctx context.Context, stringId string) {
	memcache.Delete(ctx, stringId)
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

func (obj *Pointer) GetPoint() int {
	return obj.gaeObj.Point
}

func (obj *Pointer) SetPoint(v int) {
	obj.gaeObj.Point = v
}

func (obj *Pointer) GetInfo() string {
	return obj.gaeObj.Info
}

func (obj *Pointer) GetUpdate() time.Time {
	return obj.gaeObj.Update
}

func (obj *PointerManager) Save(ctx context.Context, pointer *Pointer) error {
	if obj.memcachedOnly == true {
		pointer.UpdateMemcache(ctx)
		return nil
	} else {
		_, err := datastore.Put(ctx, pointer.gaeKey, pointer.gaeObj)
		if err == nil {
			pointer.UpdateMemcache(ctx)
		}
		return err
	}
}
