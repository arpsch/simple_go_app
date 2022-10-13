package mocks

import (
	"context"

	"github.com/arpsch/ha/model"
	"github.com/stretchr/testify/mock"
)

type MockDnsApp struct {
	mock.Mock
}

func (mda *MockDnsApp) GetLocation(ctx context.Context, pos model.Position) (float64, error) {
	args := mda.Called(ctx, pos)
	return args.Get(0).(float64), args.Error(1)
}
