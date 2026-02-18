package ports

import (
	"context"

	modelconsumable "github.com/Logstay/tech-challange-execution-service/internal/domain/consumable"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) Save(ctx context.Context, in modelconsumable.Consumable) (err error) {
	args := m.Called(ctx, in)
	return args.Error(0)
}

func (m *RepositoryMock) Get(ctx context.Context, uid string) (result GetConsumable, err error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(GetConsumable), args.Error(1)
}

func (m *RepositoryMock) List(ctx context.Context, filter Filter) (result []GetConsumable, total int64, err error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]GetConsumable), args.Get(1).(int64), args.Error(2)
}

func (m *RepositoryMock) Update(ctx context.Context, uid string, in modelconsumable.Consumable) (err error) {
	args := m.Called(ctx, uid, in)
	return args.Error(0)
}

func (m *RepositoryMock) Delete(ctx context.Context, uid string) (err error) {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *RepositoryMock) Exist(ctx context.Context, uid string) (result bool, err error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(bool), args.Error(1)
}

var _ Repository = (*RepositoryMock)(nil)
