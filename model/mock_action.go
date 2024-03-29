// Code generated by mockery v2.40.1. DO NOT EDIT.

package model

import mock "github.com/stretchr/testify/mock"

// MockAction is an autogenerated mock type for the Action type
type MockAction struct {
	mock.Mock
}

// Pawn provides a mock function with given fields:
func (_m *MockAction) Pawn() Pawn {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Pawn")
	}

	var r0 Pawn
	if rf, ok := ret.Get(0).(func() Pawn); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Pawn)
		}
	}

	return r0
}

// Position provides a mock function with given fields:
func (_m *MockAction) Position() Position {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Position")
	}

	var r0 Position
	if rf, ok := ret.Get(0).(func() Position); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Position)
		}
	}

	return r0
}

// SetPosition provides a mock function with given fields: position
func (_m *MockAction) SetPosition(position Position) {
	_m.Called(position)
}

// Type provides a mock function with given fields:
func (_m *MockAction) Type() ActionType {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Type")
	}

	var r0 ActionType
	if rf, ok := ret.Get(0).(func() ActionType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(ActionType)
	}

	return r0
}

// NewMockAction creates a new instance of MockAction. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAction(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAction {
	mock := &MockAction{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
