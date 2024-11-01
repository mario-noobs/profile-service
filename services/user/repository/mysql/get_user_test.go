package mysql

import (
	"context"
	"demo-service/helpers"
	"demo-service/services/user/entity"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/viettranx/service-context/core"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGetUsersById_Success(t *testing.T) {
	// Arrange
	time := new(helpers.Timer)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db,
		SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm: %s", err)
	}

	repo := &mysqlRepo{
		db:   gormDB,
		time: *time,
	}

	id := 1
	userData := &entity.UserDataCreation{
		Status:     "Unknown",
		Email:      "test@example.com",
		SystemRole: "User",
		FirstName:  "John",
		LastName:   "Doe",
	}
	// Updated expectation for GetAuth
	mock.ExpectQuery("SELECT(.*)").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"email", "status", "first_name", "last_name", "system_role"}).
			AddRow(userData.Email, userData.Status, userData.FirstName, userData.LastName, userData.SystemRole))

	// Act
	result, err := repo.GetUserById(context.Background(), id)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, userData.Email, result.Email)
	assert.Equal(t, userData.FirstName, result.FirstName)
	assert.Equal(t, userData.LastName, result.LastName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUsersById_IdNotExist(t *testing.T) {
	// Arrange
	time := new(helpers.Timer)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db,
		SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm: %s", err)
	}

	repo := &mysqlRepo{
		db:   gormDB,
		time: *time,
	}

	id := 1

	// Updated expectation for GetAuth
	mock.ExpectQuery("SELECT(.*)").
		WithArgs(id).
		WillReturnError(gorm.ErrRecordNotFound)

	// Act
	result, err := repo.GetUserById(context.Background(), id)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, core.ErrNotFound.Error(), err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUsersById_InvalidDataFailed(t *testing.T) {
	// Arrange
	time := new(helpers.Timer)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db,
		SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm: %s", err)
	}

	repo := &mysqlRepo{
		db:   gormDB,
		time: *time,
	}

	id := 1

	// Updated expectation for GetAuth
	mock.ExpectQuery("SELECT(.*)").
		WithArgs(id).
		WillReturnError(entity.ErrStatusIsNotValid)

	// Act
	result, err := repo.GetUserById(context.Background(), id)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, errors.Wrap(entity.ErrStatusIsNotValid, entity.ErrCannotGetUser.Error()).Error(), err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}
