package ports

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type DAOMock struct {
	mock.Mock
}

func (m *DAOMock) Save(ctx context.Context, input ConsumableDAO) (id int64, err error) {
	args := m.Called(ctx, input)
	return args.Get(0).(int64), args.Error(1)
}

func (m *DAOMock) Get(ctx context.Context, uid string) (result ConsumableDAO, err error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(ConsumableDAO), args.Error(1)
}

func (m *DAOMock) List(ctx context.Context, filter Filter) (result []ConsumableDAO, err error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]ConsumableDAO), args.Error(1)
}

func (m *DAOMock) Update(ctx context.Context, input ConsumableDAO) (err error) {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *DAOMock) Delete(ctx context.Context, uid string) (err error) {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *DAOMock) Exist(ctx context.Context, name string) (exist bool, err error) {
	args := m.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func (m *DAOMock) GetByID(ctx context.Context, id int64) (result ConsumableDAO, err error) {
	args := m.Called(ctx, id)
	return args.Get(0).(ConsumableDAO), args.Error(1)
}

func (m *DAOMock) TotalList(ctx context.Context, filter Filter) (total int64, err error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

var _ DAO = (*DAOMock)(nil)
