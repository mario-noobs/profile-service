package business

import (
	"context"
	"demo-service/proto/pb"
	"demo-service/services/user/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserById(ctx context.Context, id int) (*pb.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*pb.User), args.Error(1)
}

func (m *MockUserRepository) GetUsersByIds(ctx context.Context, ids []int) ([]pb.User, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]pb.User), args.Error(1)
}

func (m *MockUserRepository) CreateNewUser(ctx context.Context, data *entity.UserDataCreation) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func TestGetUserDetails_Success(t *testing.T) {
	userRepo := new(MockUserRepository)

	biz := NewBusiness(userRepo)

	// Setup mock expectations
	email := "test@example.com"
	firstName := "password"
	lastName := "random_salt"
	status := "hashed_password"
	authData := &pb.User{Email: email, FirstName: firstName, LastName: lastName, Status: status}

	userId := 1

	userRepo.On("GetUserById", mock.Anything, userId).Return(authData, nil)

	// Call the method
	user, err := biz.GetUserDetails(context.Background(), userId)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)

	// Assert expectations
	userRepo.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	userRepo := new(MockUserRepository)

	biz := NewBusiness(userRepo)

	// Setup mock expectations
	email := "test@example.com"
	firstName := "password"
	lastName := "random_salt"
	authData := &entity.UserDataCreation{Email: email, FirstName: firstName, LastName: lastName, Status: "status"}

	userRepo.On("CreateNewUser", mock.Anything, authData).Return(nil)

	// Call the method
	err := biz.CreateNewUser(context.Background(), authData)

	// Assertions
	assert.NoError(t, err)

	// Assert expectations
	userRepo.AssertExpectations(t)
}

func TestCreateUser_MissingRequiredField(t *testing.T) {
	userRepo := new(MockUserRepository)

	biz := NewBusiness(userRepo)

	// Setup mock expectations
	//email := "test@example.com"
	firstName := "password"
	lastName := "random_salt"
	authData := &entity.UserDataCreation{FirstName: firstName, LastName: lastName, Status: "status"}

	userRepo.On("CreateNewUser", mock.Anything, authData).Return(entity.ErrEmailIsNotValid)

	// Call the method
	err := biz.CreateNewUser(context.Background(), authData)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, err.Error(), entity.ErrCannotCreateUser.Error())

	// Assert expectations
	userRepo.AssertExpectations(t)
}
