package mysql

import (
	"context"
	"demo-service/helpers"
	"demo-service/services/user/entity"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestCreateNewUser_Success(t *testing.T) {
	// Arrange
	time2 := new(helpers.Timer)
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
		time: *time2,
	}

	userData := &entity.UserDataCreation{
		Status:     "Unknown",
		Email:      "test@example.com",
		SystemRole: "User",
		FirstName:  "John",
		LastName:   "Doe",
	}

	// Expectation for AddNewAuth
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO(.*)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), userData.FirstName, userData.LastName, userData.Email, userData.SystemRole, userData.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Act
	err = repo.CreateNewUser(context.Background(), userData)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateNewUser_Failed(t *testing.T) {
	// Arrange
	time2 := new(helpers.Timer)
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
		time: *time2,
	}

	userData := &entity.UserDataCreation{
		Status:     "Unknown",
		Email:      "test@example.com",
		SystemRole: "User",
		FirstName:  "John",
		LastName:   "Doe",
	}

	expectedError := "Cannot insert"

	// Expectation for AddNewAuth
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO(.*)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), userData.FirstName, userData.LastName, userData.Email, userData.SystemRole, userData.Status).
		WillReturnError(fmt.Errorf(expectedError)) // Mocking an insert error
	mock.ExpectRollback() // Expect a rollback since the transaction fails

	// Act
	err = repo.CreateNewUser(context.Background(), userData)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
}
