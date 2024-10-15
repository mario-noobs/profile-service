package business

import (
	"context"
	"demo-service/proto/pb"
	"demo-service/services/user/entity"
	"github.com/viettranx/service-context/core"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id int) (*pb.User, error)
	GetUsersByIds(ctx context.Context, ids []int) ([]pb.User, error)
	CreateNewUser(ctx context.Context, data *entity.UserDataCreation) error
}

type business struct {
	repository UserRepository
}

func NewBusiness(repository UserRepository) *business {
	return &business{repository: repository}
}

func (biz *business) GetUserDetails(ctx context.Context, id int) (*pb.User, error) {
	user, err := biz.repository.GetUserById(ctx, id)

	if err != nil {
		if err == core.ErrRecordNotFound {
			return nil, core.ErrNotFound.
				WithError(entity.ErrCannotGetUser.Error()).
				WithDebug(err.Error())
		}

		return nil, core.ErrInternalServerError.
			WithError(entity.ErrCannotGetUser.Error()).
			WithDebug(err.Error())
	}

	return user, nil
}

func (biz *business) GetUsersByIds(ctx context.Context, ids []int) ([]pb.User, error) {
	users, err := biz.repository.GetUsersByIds(ctx, ids)

	if err != nil {
		return nil, core.ErrNotFound.
			WithError(entity.ErrCannotGetUsers.Error()).
			WithDebug(err.Error())
	}

	return users, nil
}

func (biz *business) CreateNewUser(ctx context.Context, data *entity.UserDataCreation) error {
	err := biz.repository.CreateNewUser(ctx, data)

	if err != nil {
		return core.ErrInternalServerError.
			WithError(entity.ErrCannotCreateUser.Error()).
			WithDebug(err.Error())
	}

	return nil
}
