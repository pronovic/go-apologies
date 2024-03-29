// Code generated by mockery v2.40.1. DO NOT EDIT.

package enum

import mock "github.com/stretchr/testify/mock"

// MockEnum is an autogenerated mock type for the Enum type
type MockEnum struct {
	mock.Mock
}

// Value provides a mock function with given fields:
func (_m *MockEnum) Value() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Value")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewMockEnum creates a new instance of MockEnum. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEnum(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEnum {
	mock := &MockEnum{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
