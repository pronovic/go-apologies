// Code generated by mockery v2.40.1. DO NOT EDIT.

package model

import mock "github.com/stretchr/testify/mock"

// MockGame is an autogenerated mock type for the Game type
type MockGame struct {
	mock.Mock
}

// Completed provides a mock function with given fields:
func (_m *MockGame) Completed() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Completed")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Copy provides a mock function with given fields:
func (_m *MockGame) Copy() Game {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Copy")
	}

	var r0 Game
	if rf, ok := ret.Get(0).(func() Game); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Game)
		}
	}

	return r0
}

// CreatePlayerView provides a mock function with given fields: color
func (_m *MockGame) CreatePlayerView(color PlayerColor) (PlayerView, error) {
	ret := _m.Called(color)

	if len(ret) == 0 {
		panic("no return value specified for CreatePlayerView")
	}

	var r0 PlayerView
	var r1 error
	if rf, ok := ret.Get(0).(func(PlayerColor) (PlayerView, error)); ok {
		return rf(color)
	}
	if rf, ok := ret.Get(0).(func(PlayerColor) PlayerView); ok {
		r0 = rf(color)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(PlayerView)
		}
	}

	if rf, ok := ret.Get(1).(func(PlayerColor) error); ok {
		r1 = rf(color)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Deck provides a mock function with given fields:
func (_m *MockGame) Deck() Deck {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Deck")
	}

	var r0 Deck
	if rf, ok := ret.Get(0).(func() Deck); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Deck)
		}
	}

	return r0
}

// History provides a mock function with given fields:
func (_m *MockGame) History() []History {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for History")
	}

	var r0 []History
	if rf, ok := ret.Get(0).(func() []History); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]History)
		}
	}

	return r0
}

// PlayerCount provides a mock function with given fields:
func (_m *MockGame) PlayerCount() int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for PlayerCount")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// Players provides a mock function with given fields:
func (_m *MockGame) Players() map[PlayerColor]Player {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Players")
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

// Started provides a mock function with given fields:
func (_m *MockGame) Started() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Started")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Track provides a mock function with given fields: action, player, card
func (_m *MockGame) Track(action string, player Player, card Card) {
	_m.Called(action, player, card)
}

// Winner provides a mock function with given fields:
func (_m *MockGame) Winner() *Player {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Winner")
	}

	var r0 *Player
	if rf, ok := ret.Get(0).(func() *Player); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Player)
		}
	}

	return r0
}

// NewMockGame creates a new instance of MockGame. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGame(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGame {
	mock := &MockGame{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}