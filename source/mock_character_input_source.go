// Code generated by mockery v2.40.1. DO NOT EDIT.

package source

import (
	model "github.com/pronovic/go-apologies/model"
	mock "github.com/stretchr/testify/mock"
)

// MockCharacterInputSource is an autogenerated mock type for the CharacterInputSource type
type MockCharacterInputSource struct {
	mock.Mock
}

// ChooseMove provides a mock function with given fields: mode, view, legalMoves
func (_m *MockCharacterInputSource) ChooseMove(mode model.GameMode, view model.PlayerView, legalMoves []model.Move) (model.Move, error) {
	ret := _m.Called(mode, view, legalMoves)

	if len(ret) == 0 {
		panic("no return value specified for ChooseMove")
	}

	var r0 model.Move
	var r1 error
	if rf, ok := ret.Get(0).(func(model.GameMode, model.PlayerView, []model.Move) (model.Move, error)); ok {
		return rf(mode, view, legalMoves)
	}
	if rf, ok := ret.Get(0).(func(model.GameMode, model.PlayerView, []model.Move) model.Move); ok {
		r0 = rf(mode, view, legalMoves)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.Move)
		}
	}

	if rf, ok := ret.Get(1).(func(model.GameMode, model.PlayerView, []model.Move) error); ok {
		r1 = rf(mode, view, legalMoves)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Name provides a mock function with given fields:
func (_m *MockCharacterInputSource) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewMockCharacterInputSource creates a new instance of MockCharacterInputSource. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCharacterInputSource(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCharacterInputSource {
	mock := &MockCharacterInputSource{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
