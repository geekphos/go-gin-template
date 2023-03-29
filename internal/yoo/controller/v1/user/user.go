package user

import (
	"phos.cc/yoo/internal/yoo/biz"
	"phos.cc/yoo/internal/yoo/store"
)

type UserController struct {
	b biz.Biz
}

// New create a new user controller.
func New(ds store.IStore) *UserController {
	return &UserController{b: biz.NewBiz(ds)}
}
