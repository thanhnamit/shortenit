package internal

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/config"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/core"
	"github.com/thanhnamit/shortenit/api-shortenit-v1/internal/platform"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestNewAlias(t *testing.T) {
	mAliasSvc := new(MockAliasService)
	mAliasSvc.On("GetNewAlias", mock.Anything).Return("test", nil)

	mRepo := new(MockUserRepository)
	mRepo.On("GetUserByEmail", mock.Anything, mock.Anything).Return(&core.User{
		ID:        primitive.NewObjectID(),
		Name:      "abc",
		Email:     "abc@test.com",
		CreatedAt: time.Now(),
		LastLogin: time.Now(),
		Aliases:   []core.Alias{},
	})
	mRepo.On("SaveUser", mock.Anything, mock.Anything).Return(nil)

	aRepo := new(MockAliasRepository)
	aRepo.On("SaveAlias", mock.Anything, mock.Anything).Return(nil)

	svc := NewService(mAliasSvc, mRepo, aRepo, &config.Config{})
	ctx := context.WithValue(context.Background(), platform.ContextKey(platform.CtxBasePath), "http://localhost:8085/shortenit")
	res, err := svc.GetNewAlias(ctx, core.ShortenURLRequest{
		OriginalURL: "http://test.decmo",
		CustomAlias: "",
		UserEmail:   "abc@test.com",
	})

	t.Logf("Test 0:\tShould return response object")
	{
		mAliasSvc.AssertExpectations(t)
		mRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, "http://localhost:8085/shortenit/test", res.URL)
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

func (m *MockUserRepository) SaveUser(ctx context.Context, user *core.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *core.User) error {
	return nil
}

func (m *MockUserRepository) Close(ctx context.Context) {
	return
}

func (m *MockUserRepository) GetAllUsers(ctx context.Context) ([]*core.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*core.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*core.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*core.User), nil
}

type MockAliasRepository struct {
	mock.Mock
}

func (r *MockAliasRepository) GetAliasByKey(ctx context.Context, alias string) (*core.Alias, error) {
	args := r.Called(ctx, alias)
	return args.Get(0).(*core.Alias), args.Error(1)
}

func (r *MockAliasRepository) SaveAlias(ctx context.Context, alias *core.Alias) error {
	args := r.Called(ctx, alias)
	return args.Error(0)
}
