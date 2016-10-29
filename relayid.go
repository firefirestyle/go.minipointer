package minipointer

import (
	"time"

	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

const (
	TypeTwitter  = "twitter"
	TypeFacebook = "facebook"
	TypePointer  = "pointer"
)

const (
	TypeProjectId   = "ProjectId"
	TypePointerName = "IdentifyName"
	TypePointerId   = "IdentifyId"
	TypePointerType = "PointerType"
	TypeValue       = "UserName"
	TypeInfo        = "Info"
	TypeUpdate      = "Update"
	TypeSign        = "Sign"
)

type GaePointerItem struct {
	ProjectId   string
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
	ProjectId string
}

type Pointer struct {
	gaeObj *GaePointerItem
	gaeKey *datastore.Key
	kind   string
}

type PointerManager struct {
	kind      string
	projectId string
}

func NewPointerManager(config PointerManagerConfig) *PointerManager {
	return &PointerManager{
		kind:      config.Kind,
		projectId: config.ProjectId,
	}
}

func (obj *Pointer) ToJson() []byte {
	propObj := miniprop.NewMiniProp()
	propObj.SetString(TypeProjectId, obj.gaeObj.ProjectId)
	propObj.SetString(TypePointerName, obj.gaeObj.PointerName)
	propObj.SetString(TypePointerId, obj.gaeObj.PointerId)
	propObj.SetString(TypeValue, obj.gaeObj.Value)
	propObj.SetString(TypeInfo, obj.gaeObj.Info)
	propObj.SetTime(TypeUpdate, obj.gaeObj.Update)
	propObj.SetString(TypeSign, obj.gaeObj.Sign)
	return propObj.ToJson()
}

func (obj *Pointer) SetValueFromJson(data []byte) {
	propObj := miniprop.NewMiniPropFromJson(data)
	obj.gaeObj.ProjectId = propObj.GetString(TypeProjectId, "")
	obj.gaeObj.PointerName = propObj.GetString(TypePointerName, "")
	obj.gaeObj.PointerId = propObj.GetString(TypePointerId, "")
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

func (obj *PointerManager) GetPointerForTwitter(ctx context.Context, screenName string, userId string, oauthToken string) *Pointer {
	return obj.GetPointerWithNew(ctx, screenName, userId, TypeTwitter, map[string]string{"token": oauthToken})
}
func (obj *PointerManager) GetPointerAsPointer(ctx context.Context, userName string) *Pointer {
	return obj.GetPointerWithNew(ctx, userName, userName, TypePointer, map[string]string{})
}

func (obj *PointerManager) NewPointer(ctx context.Context, screenName string, //
	userId string, identifyType string, infos map[string]string) *Pointer {
	gaeKey := obj.NewPointerGaeKey(ctx, userId, identifyType)
	gaeObj := GaePointerItem{
		PointerName: screenName,
		PointerId:   userId,
		PointerType: TypeTwitter,
		ProjectId:   obj.projectId,
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
	return obj.kind + ":" + obj.projectId + ":" + identifyType + ":" + identify
}
