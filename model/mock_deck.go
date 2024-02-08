// Code generated by mockery v2.40.1. DO NOT EDIT.

package model

import mock "github.com/stretchr/testify/mock"

// MockDeck is an autogenerated mock type for the Deck type
type MockDeck struct {
	mock.Mock
}

// Copy provides a mock function with given fields:
func (_m *MockDeck) Copy() Deck {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Copy")
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

// Discard provides a mock function with given fields: card
func (_m *MockDeck) Discard(card Card) error {
	ret := _m.Called(card)

	if len(ret) == 0 {
		panic("no return value specified for Discard")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(Card) error); ok {
		r0 = rf(card)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Draw provides a mock function with given fields:
func (_m *MockDeck) Draw() (Card, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Draw")
	}

	var r0 Card
	var r1 error
	if rf, ok := ret.Get(0).(func() (Card, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() Card); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Card)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockDeck creates a new instance of MockDeck. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDeck(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDeck {
	mock := &MockDeck{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}