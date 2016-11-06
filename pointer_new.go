package minipointer

import (
	"time"

	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

func (obj *PointerManager) DeleteFromPointer(ctx context.Context, item *Pointer) error {
	return obj.Delete(ctx, item.GetId(), item.GetType())
}

func (obj *PointerManager) Delete(ctx context.Context, userId, identifyType string) error {
	//	Debug(ctx, ">> Pointer >>> : "+userId+" : "+identifyType+"==")
	return datastore.Delete(ctx, obj.NewPointerGaeKey(ctx, userId, identifyType))
}

func (obj *PointerManager) GetPointerForTwitter(ctx context.Context, screenName string, userId string, oauthToken string) *Pointer {
	return obj.GetPointerWithNew(ctx, screenName, userId, TypeTwitter, map[string]string{"token": oauthToken})
}

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
		Debug(ctx, "====> Failed to get pointer:"+identify+":"+identifyType)
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
	prop := miniprop.NewMiniProp()
	prop.SetString("k", obj.kind)
	prop.SetString("g", obj.rootGroup)
	prop.SetString("i", identify)
	prop.SetString("t", identifyType)
	return string(prop.ToJson())
}

type PointerKeyInfo struct {
	IdentifyType string
	Identify     string
	Kind         string
	RootGroup    string
}

func (obj *PointerManager) GetKeyInfoFromStringId(stringId string) PointerKeyInfo {
	prop := miniprop.NewMiniPropFromJson([]byte(stringId))
	return PointerKeyInfo{
		Kind:         prop.GetString("k", ""),
		RootGroup:    prop.GetString("g", ""),
		Identify:     prop.GetString("i", ""),
		IdentifyType: prop.GetString("t", ""),
	}
}
