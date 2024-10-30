package mysql

import (
	"context"
	"demo-service/proto/pb"
	"demo-service/services/user/entity"
	"github.com/pkg/errors"
	"github.com/viettranx/service-context/core"
	"gorm.io/gorm"
)

func (repo *mysqlRepo) GetUsersByIds(ctx context.Context, ids []int) ([]pb.User, error) {
	var result []pb.User
	method := "GetUsersByIds_SQL"
	repo.time.Start()
	logger.Info("request", "method", method)
	if err := repo.db.
		Table(entity.User{}.TableName()).
		Where("id in (?)", ids).
		Find(&result).Error; err != nil {
		logger.Error("response", "method", method, "err", err, "ms", repo.time.End())
		return nil, errors.Wrap(err, entity.ErrCannotGetUser.Error())
	}
	logger.Info("response", "method", method, "data", result, "ms", repo.time.End())
	return result, nil
}

func (repo *mysqlRepo) GetUserById(ctx context.Context, id int) (*pb.User, error) {
	var data pb.User
	method := "GetUserById_SQL"
	repo.time.Start()
	logger.Info("request", "method", method)
	if err := repo.db.
		Table("users").
		Where("id = ?", id).
		First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error("response", "method", method, "err", err, "ms", repo.time.End())
			return nil, core.ErrNotFound
		}

		return nil, errors.Wrap(err, entity.ErrCannotGetUser.Error())
	}
	logger.Info("response", "method", method, "data", data, "ms", repo.time.End())
	return &data, nil
}
