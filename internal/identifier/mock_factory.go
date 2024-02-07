// Code generated by mockery v2.40.1. DO NOT EDIT.

package identifier

import mock "github.com/stretchr/testify/mock"

// MockFactory is an autogenerated mock type for the Factory type
type MockFactory struct {
	mock.Mock
}

// RandomId provides a mock function with given fields:
func (_m *MockFactory) RandomId() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RandomId")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewMockFactory creates a new instance of MockFactory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockFactory(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockFactory {
	mock := &MockFactory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
