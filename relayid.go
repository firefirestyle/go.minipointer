package minipointer

import (
	"time"

	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

const (
	TypeTwitter  = "twitter"
	TypeFacebook = "facebook"
	TypePointer  = "pointer"
)

const (
	TypeRootGroup   = "RootGroup"
	TypePointerName = "IdentifyName"
	TypePointerId   = "IdentifyId"
	TypePointerType = "PointerType"
	TypeValue       = "UserName"
	TypeInfo        = "Info"
	TypeUpdate      = "Update"
	TypeSign        = "Sign"
)

type GaePointerItem struct {
	RootGroup   string
	PointerName string
	PointerId   string
	PointerType string
	Value       string
	Info        string
	Update      time.Time
	Sign        string
}

type PointerManagerConfig struct {
	Kind      string
	RootGroup string
}

type Pointer struct {
	gaeObj *GaePointerItem
	gaeKey *datastore.Key
	kind   string
}

type PointerManager struct {
	kind      string
	rootGroup string
}

func NewPointerManager(config PointerManagerConfig) *PointerManager {
	return &PointerManager{
		kind:      config.Kind,
		rootGroup: config.RootGroup,
	}
}

func (obj *Pointer) ToJson() []byte {
	propObj := miniprop.NewMiniProp()
	propObj.SetString(TypeRootGroup, obj.gaeObj.RootGroup)
	propObj.SetString(TypePointerName, obj.gaeObj.PointerName)
	propObj.SetString(TypePointerId, obj.gaeObj.PointerId)
	propObj.SetString(TypePointerType, obj.gaeObj.PointerType)

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

func (obj *Pointer) GetInfo() string {
	return obj.gaeObj.Info
}

func (obj *Pointer) GetUpdate() time.Time {
	return obj.gaeObj.Update
}

func (obj *Pointer) Save(ctx context.Context) error {
	_, err := datastore.Put(ctx, obj.gaeKey, obj.gaeObj)
	if err == nil {
		obj.UpdateMemcache(ctx)
	}
	return err
}

func (obj *PointerManager) DeleteFromPointer(ctx context.Context, item *Pointer) error {
	return obj.Delete(ctx, item.GetId(), item.GetType())
}

func (obj *PointerManager) Delete(ctx context.Context, userId, identifyType string) error {
	//	Debug(ctx, ">> Pointer >>> : "+userId+" : "+identifyType+"==")
	return datastore.Delete(ctx, obj.NewPointerGaeKey(ctx, userId, identifyType))
}

/*
func (obj *PointerManager) GetPointerForTwitter(ctx context.Context, screenName string, userId string, oauthToken string) *Pointer {
	return obj.GetPointerWithNew(ctx, screenName, userId, TypeTwitter, map[string]string{"token": oauthToken})
}
*/
func (obj *PointerManager) GetPointerForRelayId(ctx context.Context, value string) *Pointer {
	return obj.GetPointerWithNew(ctx, value, value, TypePointer, map[string]string{})
}

func (obj *PointerManager) NewPointer(ctx context.Context, screenName string, //
	userId string, identifyType string, infos map[string]string) *Pointer {
	gaeKey := obj.NewPointerGaeKey(ctx, userId, identifyType)
	gaeObj := GaePointerItem{
		PointerName: screenName,
		PointerId:   userId,
		PointerType: identifyType,
		RootGroup:   obj.rootGroup,
	}
	propObj := miniprop.NewMiniPropFromJson([]byte(gaeObj.Info))
	for k, v := range infos {
		propObj.SetString(k, v)
	}
	gaeObj.Info = string(propObj.ToJson())
	gaeObj.Update = time.Now()
	return &Pointer{
		gaeObj: &gaeObj,
		gaeKey: gaeKey,
		kind:   obj.kind,
	}
}

func (obj *PointerManager) GetPointer(ctx context.Context, identify string, identifyType string) (*Pointer, error) {
	gaeKey := obj.NewPointerGaeKey(ctx, identify, identifyType)
	gaeObj := GaePointerItem{}

	//
	// mem
	memItemObj, errMemObj := memcache.Get(ctx, obj.MakePointerStringId(identify, identifyType))
	if errMemObj == nil {
		ret := &Pointer{
			gaeObj: &gaeObj,
			gaeKey: gaeKey,
			kind:   obj.kind,
		}
		ret.SetValueFromJson(memItemObj.Value)
		return ret, nil
	}
	//
	// db
	err := datastore.Get(ctx, gaeKey, &gaeObj)
	if err != nil {
		return nil, err
	}
	ret := &Pointer{
		gaeObj: &gaeObj,
		gaeKey: gaeKey,
		kind:   obj.kind,
	}
	//
	//
	ret.UpdateMemcache(ctx)
	return ret, nil
}

func (obj *PointerManager) GetPointerWithNew(ctx context.Context, screenName string, userId string, userIdType string, infos map[string]string) *Pointer {
	Debug(ctx, ">>>>>>:userIdType:"+userIdType)
	relayObj, err := obj.GetPointer(ctx, userId, userIdType)
	if err != nil {
		relayObj = obj.NewPointer(ctx, screenName, userId, userIdType, infos)
	}
	//
	propObj := miniprop.NewMiniPropFromJson([]byte(relayObj.gaeObj.Info))
	for k, v := range infos {
		propObj.SetString(k, v)
	}
	relayObj.gaeObj.Info = string(propObj.ToJson())
	relayObj.gaeObj.Update = time.Now()
	return relayObj
}

func (obj *PointerManager) NewPointerGaeKey(ctx context.Context, identify string, identifyType string) *datastore.Key {
	return datastore.NewKey(ctx, obj.kind, obj.MakePointerStringId(identify, identifyType), 0, nil)
}

func (obj *PointerManager) MakePointerStringId(identify string, identifyType string) string {
	return obj.kind + ":" + obj.rootGroup + ":" + identifyType + ":" + identify
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
