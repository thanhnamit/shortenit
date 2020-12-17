package internal

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/config"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestNewAlias(t *testing.T) {
	mAliasSvc := new(MockAliasService)
	mAliasSvc.On("GetNewAlias", mock.Anything).Return("test", nil)

	mRepo := new(MockUserRepository)
	mRepo.On("GetAllUsers", mock.Anything).Return([]*core.User{
		{
			ID:        primitive.NewObjectID(),
			Name:      "Nam",
			Email:     "t",
			CreatedAt: time.Time{},
			LastLogin: time.Time{},
		},
	}, nil)

	svc := NewService(mAliasSvc, mRepo, &config.Config{})
	res, err := svc.NewAlias(context.TODO(), core.ShortenURLRequest{
		OriginalURL: "http://test.decmo",
		CustomAlias: "",
		UserEmail:   "",
	})

	t.Logf("Test 0:\tShould return response object")
	{
		mAliasSvc.AssertExpectations(t)
		mRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, "test", res.URL)
	}
}

type MockAliasService struct {
	mock.Mock
}

func (m *MockAliasService) GetNewAlias(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *core.User) error {
	return nil
}

func (m MockUserRepository) Close(ctx context.Context) {
	return
}

func (m MockUserRepository) GetAllUsers(ctx context.Context) ([]*core.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*core.User), args.Error(1)
}

func (m MockUserRepository) GetUserByEmail(ctx context.Context, email string) *core.User {
	args := m.Called(ctx, email)
	return args.Get(0).(*core.User)
}


