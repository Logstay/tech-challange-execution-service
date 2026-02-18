package ports

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type DAOMock struct {
	mock.Mock
}

func (m *DAOMock) Save(ctx context.Context, input InventoryDAO) (err error) {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *DAOMock) Get(ctx context.Context, consumableID int64) (result InventoryDAO, err error) {
	args := m.Called(ctx, consumableID)
	return args.Get(0).(InventoryDAO), args.Error(1)
}

func (m *DAOMock) List(ctx context.Context, filter Filter) (result []InventoryDAO, err error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]InventoryDAO), args.Error(1)
}

func (m *DAOMock) Update(ctx context.Context, input InventoryDAO) (err error) {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *DAOMock) Delete(ctx context.Context, consumableID int64) (err error) {
	args := m.Called(ctx, consumableID)
	return args.Error(0)
}

func (m *DAOMock) Decrease(ctx context.Context, input InventoryDAO) (err error) {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *DAOMock) Increase(ctx context.Context, input InventoryDAO) (err error) {
	args := m.Called(ctx, input)
	return args.Error(0)
}

var _ DAO = (*DAOMock)(nil)
