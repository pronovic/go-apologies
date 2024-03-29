// Code generated by mockery v2.40.1. DO NOT EDIT.

package model

import mock "github.com/stretchr/testify/mock"

// MockPlayerView is an autogenerated mock type for the PlayerView type
type MockPlayerView struct {
	mock.Mock
}

// AllPawns provides a mock function with given fields:
func (_m *MockPlayerView) AllPawns() []Pawn {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for AllPawns")
	}

	var r0 []Pawn
	if rf, ok := ret.Get(0).(func() []Pawn); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Pawn)
		}
	}

	return r0
}

// Copy provides a mock function with given fields:
func (_m *MockPlayerView) Copy() PlayerView {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Copy")
	}

	var r0 PlayerView
	if rf, ok := ret.Get(0).(func() PlayerView); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(PlayerView)
		}
	}

	return r0
}

// GetPawn provides a mock function with given fields: prototype
func (_m *MockPlayerView) GetPawn(prototype Pawn) Pawn {
	ret := _m.Called(prototype)

	if len(ret) == 0 {
		panic("no return value specified for GetPawn")
	}

	var r0 Pawn
	if rf, ok := ret.Get(0).(func(Pawn) Pawn); ok {
		r0 = rf(prototype)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Pawn)
		}
	}

	return r0
}

// Opponents provides a mock function with given fields:
func (_m *MockPlayerView) Opponents() map[PlayerColor]Player {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Opponents")
	}

	var r0 map[PlayerColor]Player
	if rf, ok := ret.Get(0).(func() map[PlayerColor]Player); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[PlayerColor]Player)
		}
	}

	return r0
}

// Player provides a mock function with given fields:
func (_m *MockPlayerView) Player() Player {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Player")
	}

	var r0 Player
	if rf, ok := ret.Get(0).(func() Player); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Player)
		}
	}

	return r0
}

// NewMockPlayerView creates a new instance of MockPlayerView. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPlayerView(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPlayerView {
	mock := &MockPlayerView{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
