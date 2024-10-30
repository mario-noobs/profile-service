package mysql

import (
	"context"
	"demo-service/services/user/entity"
	"github.com/pkg/errors"
	"log/slog"
	"os"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func (repo *mysqlRepo) CreateNewUser(ctx context.Context, data *entity.UserDataCreation) error {
	method := "CreateNewUser_SQL"
	repo.time.Start()
	logger.Info("request", "method", method)

	if err := repo.db.Table(data.TableName()).Create(data).Error; err != nil {
		logger.Error("response", "method", method, "err", err, "ms", repo.time.End())
		return errors.WithStack(err)
	}

	logger.Info("response", "method", method, "data", data, "ms", repo.time.End())
	return nil
}
