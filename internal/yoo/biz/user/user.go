package user

import (
	"context"
	"phos.cc/yoo/pkg/auth"
	"phos.cc/yoo/pkg/token"
	"regexp"

	"github.com/jinzhu/copier"

	"phos.cc/yoo/internal/pkg/errno"
	"phos.cc/yoo/internal/pkg/model"
	"phos.cc/yoo/internal/yoo/store"
	v1 "phos.cc/yoo/pkg/api/yoo/v1"
)

type UserBiz interface {
	ChangePassword(ctx context.Context, email string, r *v1.ChangePasswordRequest) error
	Create(ctx context.Context, r *v1.CreateUserRequest) error
	Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error)
}

type userBiz struct {
	ds store.IStore
}

var _UserBiz = (*userBiz)(nil)

func New(ds store.IStore) *userBiz {
	return &userBiz{ds: ds}
}

func (b *userBiz) Create(ctx context.Context, r *v1.CreateUserRequest) error {
	var userM = &model.UserM{}
	_ = copier.Copy(userM, r)

	if err := b.ds.Users().Create(ctx, userM); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'email'", err.Error()); match {
			return errno.ErrUserAlreadyExist
		}
		return err
	}

	return nil
}

func (b *userBiz) Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error) {
	user, err := b.ds.Users().Get(ctx, r.Email)
	if err != nil {
		return nil, errno.ErrUserNotFound
	}

	// compare password
	if err := auth.Compare(user.Password, r.Password); err != nil {
		return nil, errno.ErrPasswordIncorrect
	}

	// generate token
	t, err := token.Sign(user.Email)
	if err != nil {
		return nil, errno.ErrSignToken
	}

	return &v1.LoginResponse{
		Token: t,
	}, nil
}

func (b *userBiz) ChangePassword(ctx context.Context, email string, r *v1.ChangePasswordRequest) error {
	userM, err := b.ds.Users().Get(ctx, email)
	if err != nil {
		return errno.ErrUserNotFound
	}

	if err := auth.Compare(userM.Password, r.OldPassword); err != nil {
		return errno.ErrPasswordIncorrect
	}

	userM.Password, _ = auth.Encrypt(r.NewPassword)
	if err := b.ds.Users().Update(ctx, userM); err != nil {
		return errno.InternalServerError
	}
	return nil
}
