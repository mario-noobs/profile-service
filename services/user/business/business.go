package business

import (
	"context"
	"demo-service/helpers"
	"demo-service/proto/pb"
	"demo-service/services/user/entity"
	"github.com/viettranx/service-context/core"
	"log/slog"
	"os"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id int) (*pb.User, error)
	GetUsersByIds(ctx context.Context, ids []int) ([]pb.User, error)
	CreateNewUser(ctx context.Context, data *entity.UserDataCreation) error
}

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type business struct {
	repository UserRepository
	time       helpers.Timer
}

func NewBusiness(repository UserRepository) *business {
	return &business{repository: repository}
}

func (biz *business) GetUserDetails(ctx context.Context, id int) (*pb.User, error) {

	method := "GetUserDetails_Business"
	biz.time.Start()
	logger.Info("request", "method", method)

	user, err := biz.repository.GetUserById(ctx, id)

	if err != nil {
		if err == core.ErrRecordNotFound {
			return nil, core.ErrNotFound.
				WithError(entity.ErrCannotGetUser.Error()).
				WithDebug(err.Error())
		}
		logger.Error("response", "method", method, "err", err, "ms", biz.time.End())
		return nil, core.ErrInternalServerError.
			WithError(entity.ErrCannotGetUser.Error()).
			WithDebug(err.Error())
	}
	logger.Info("response", "method", method, "data", user, "ms", biz.time.End())
	return user, nil
}

func (biz *business) GetUsersByIds(ctx context.Context, ids []int) ([]pb.User, error) {

	method := "GetUsersByIds_Business"
	biz.time.Start()
	logger.Info("request", "method", method)

	users, err := biz.repository.GetUsersByIds(ctx, ids)

	if err != nil {
		logger.Error("response", "method", method, "err", err, "ms", biz.time.End())
		return nil, core.ErrNotFound.
			WithError(entity.ErrCannotGetUsers.Error()).
			WithDebug(err.Error())
	}
	logger.Info("response", "method", method, "data", users, "ms", biz.time.End())

	return users, nil
}

func (biz *business) CreateNewUser(ctx context.Context, data *entity.UserDataCreation) error {
	method := "CreateNewUser_Business"
	biz.time.Start()
	logger.Info("request", "method", method)

	err := biz.repository.CreateNewUser(ctx, data)

	if err != nil {
		logger.Error("response", "method", method, "err", err, "ms", biz.time.End())
		return core.ErrInternalServerError.
			WithError(entity.ErrCannotCreateUser.Error()).
			WithDebug(err.Error())
	}
	logger.Info("response", "method", method, "data", true, "ms", biz.time.End())
	return nil
}
