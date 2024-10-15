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

	if err := repo.db.
		Table(entity.User{}.TableName()).
		Where("id in (?)", ids).
		Find(&result).Error; err != nil {
		return nil, errors.Wrap(err, entity.ErrCannotGetUser.Error())
	}

	return result, nil
}

func (repo *mysqlRepo) GetUserById(ctx context.Context, id int) (*pb.User, error) {
	var data pb.User

	if err := repo.db.
		Table("users").
		Where("id = ?", id).
		First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrNotFound
		}
		return nil, errors.Wrap(err, entity.ErrCannotGetUser.Error())
	}

	return &data, nil
}
