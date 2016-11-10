package minipointer

import (
	//	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

const (
	TypeTwitter  = "twitter"
	TypeFacebook = "facebook"
	TypePointer  = "pointer"
)

type PointerManagerConfig struct {
	Kind          string
	RootGroup     string
	MemcachedOnly bool
}

type PointerManager struct {
	kind          string
	rootGroup     string
	memcachedOnly bool
}

func NewPointerManager(config PointerManagerConfig) *PointerManager {
	return &PointerManager{
		kind:          config.Kind,
		rootGroup:     config.RootGroup,
		memcachedOnly: config.MemcachedOnly,
	}
}

func (obj *PointerManager) IsMemcachedOnly() bool {
	return obj.memcachedOnly
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
