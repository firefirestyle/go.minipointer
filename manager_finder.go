package minipointer

import (
	//	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type FoundPointers struct {
	Keys       []string
	CursorNext string
	CursorOne  string
}

func (obj *PointerManager) makeCursorSrc(founds *datastore.Iterator) string {
	c, e := founds.Cursor()
	if e == nil {
		return c.String()
	} else {
		return ""
	}
}

func (obj *PointerManager) newCursorFromSrc(cursorSrc string) *datastore.Cursor {
	c1, e := datastore.DecodeCursor(cursorSrc)
	if e != nil {
		return nil
	} else {
		return &c1
	}
}

func (obj *PointerManager) FindFromOwner(ctx context.Context, cursorSrc string, owner string) FoundPointers {
	return obj.FindPointerFromQuery(ctx, obj.NewQueryFromOwner(owner), cursorSrc)
}

func (obj *PointerManager) NewQueryFromOwner(owner string) *datastore.Query {
	q := datastore.NewQuery(obj.kind)
	q = q.Filter("RootGroup =", obj.rootGroup)
	q = q.Filter("Owner = ", owner)
	return q
}

func (obj *PointerManager) NewQueryFromOwnerAndRootGroup(owner, rootGroup string) *datastore.Query {
	q := datastore.NewQuery(obj.kind)
	if rootGroup != "" {
		q = q.Filter("RootGroup =", rootGroup)
	}
	q = q.Filter("Owner = ", owner)
	return q
}

func (obj *PointerManager) NewQueryFromPointerId(v string) *datastore.Query {
	q := datastore.NewQuery(obj.kind)
	q = q.Filter("RootGroup =", obj.rootGroup)
	q = q.Filter("PointerId = ", v)
	return q
}

func (obj *PointerManager) FindPointerFromQueryAll(ctx context.Context, q *datastore.Query) FoundPointers {
	founded := obj.FindPointerFromQuery(ctx, q, "")
	oneCursor := founded.CursorOne
	nextCursor := founded.CursorNext
	keys := make([]string, 0)
	for {
		if len(founded.Keys) <= 0 {
			break
		}
		for _, v := range founded.Keys {
			keys = append(keys, v)
		}
		prevFounded := founded
		founded = obj.FindPointerFromQuery(ctx, q, nextCursor)
		nextCursor = founded.CursorNext
		if prevFounded.CursorOne == founded.CursorOne {
			break
		}
	}
	return FoundPointers{
		Keys:       keys,
		CursorNext: nextCursor,
		CursorOne:  oneCursor,
	}
}

func (obj *PointerManager) FindPointerFromQuery(ctx context.Context, q *datastore.Query, cursorSrc string) FoundPointers {
	cursor := obj.newCursorFromSrc(cursorSrc)
	if cursor != nil {
		q = q.Start(*cursor)
	}
	q = q.KeysOnly()
	founds := q.Run(ctx)

	var pointerKeys []string = make([]string, 0)

	var cursorNext string = ""
	var cursorOne string = ""
	for i := 0; ; i++ {
		key, err := founds.Next(nil)

		if err != nil || err == datastore.Done {
			break
		} else {
			pointerKeys = append(pointerKeys, key.StringID())
		}
		if i == 0 {
			cursorOne = obj.makeCursorSrc(founds)
		}
	}
	cursorNext = obj.makeCursorSrc(founds)
	return FoundPointers{
		Keys:       pointerKeys,
		CursorNext: cursorNext,
		CursorOne:  cursorOne,
	}
}
