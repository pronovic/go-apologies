// Code generated by mockery v2.40.1. DO NOT EDIT.

package engine

import (
	model "github.com/pronovic/go-apologies/model"
	mock "github.com/stretchr/testify/mock"

	source "github.com/pronovic/go-apologies/source"
)

// MockCharacter is an autogenerated mock type for the Character type
type MockCharacter struct {
	mock.Mock
}

// ChooseMove provides a mock function with given fields: mode, view, legalMoves
func (_m *MockCharacter) ChooseMove(mode model.GameMode, view model.PlayerView, legalMoves []model.Move) (model.Move, error) {
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

// Color provides a mock function with given fields:
func (_m *MockCharacter) Color() model.PlayerColor {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Color")
	}

	var r0 model.PlayerColor
	if rf, ok := ret.Get(0).(func() model.PlayerColor); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.PlayerColor)
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *MockCharacter) Name() string {
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

// SetColor provides a mock function with given fields: color
func (_m *MockCharacter) SetColor(color model.PlayerColor) {
	_m.Called(color)
}

// Source provides a mock function with given fields:
func (_m *MockCharacter) Source() source.CharacterInputSource {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Source")
	}

	var r0 source.CharacterInputSource
	if rf, ok := ret.Get(0).(func() source.CharacterInputSource); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(source.CharacterInputSource)
		}
	}

	return r0
}

// NewMockCharacter creates a new instance of MockCharacter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCharacter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCharacter {
	mock := &MockCharacter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
