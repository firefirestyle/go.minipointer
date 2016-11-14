package minipointer

import (
	"time"

	"errors"

	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

//
//
//
func (obj *PointerManager) GetPointerWithNewForTwitter(ctx context.Context, screenName string, userId string, oauthToken string) *Pointer {
	return obj.GetPointerWithNew(ctx, //screenName,
		userId, TypeTwitter, map[string]string{"token": oauthToken})
}

func (obj *PointerManager) GetPointerWithNewForRelayId(ctx context.Context, value string) *Pointer {
	return obj.GetPointerWithNew(ctx, // value,
		value, TypePointer, map[string]string{})
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
	if obj.memcachedOnly == true {
		return nil, errors.New("Failed to get pointer asis Memcached Only")
	}
	//
	// db
	err := datastore.Get(ctx, gaeKey, &gaeObj)
	if err != nil {
		Debug(ctx, "====> Failed to get pointer:"+identify+":"+identifyType)
		return nil, errors.New(err.Error() + ":" + obj.kind + ":" + obj.MakePointerStringId(identify, identifyType))
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

func (obj *PointerManager) GetPointerWithNew(ctx context.Context, // screenName string,
	identity string, identityType string, infos map[string]string) *Pointer {
	// Debug(ctx, ">>>>>>:userIdType:"+userIdType)
	relayObj, err := obj.GetPointer(ctx, identity, identityType)
	if err != nil {
		relayObj = obj.NewPointer(ctx, //screenName,
			identity, identityType, infos)
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
