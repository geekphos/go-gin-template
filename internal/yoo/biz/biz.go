package biz

import (
	"phos.cc/yoo/internal/yoo/biz/user"
	"phos.cc/yoo/internal/yoo/store"
)

type Biz interface {
	Users() user.UserBiz
}

type biz struct {
	ds store.IStore
}

var _Biz = (*biz)(nil)

// NewBiz returns a new biz.
func NewBiz(ds store.IStore) *biz {
	return &biz{ds: ds}
}

func (b *biz) Users() user.UserBiz {
	return user.New(b.ds)
}
