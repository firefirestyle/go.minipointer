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

func (obj *PointerManager) FindFromOwner(ctx context.Context, cursorSrc string, owner string) *FoundPointers {
	q := datastore.NewQuery(obj.kind)
	q = q.Filter("RootGroup =", obj.rootGroup)
	q = q.Filter("Owner = ", owner)
	return obj.FindPointerFromQuery(ctx, q, cursorSrc)
}

func (obj *PointerManager) FindPointerFromQuery(ctx context.Context, q *datastore.Query, cursorSrc string) *FoundPointers {
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
	return &FoundPointers{
		Keys:       pointerKeys,
		CursorNext: cursorNext,
		CursorOne:  cursorOne,
	}
}
