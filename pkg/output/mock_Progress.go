// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package output

import mock "github.com/stretchr/testify/mock"

// MockProgress is an autogenerated mock type for the Progress type
type MockProgress struct {
	mock.Mock
}

// Inc provides a mock function with given fields:
func (_m *MockProgress) Inc() {
	_m.Called()
}

// Start provides a mock function with given fields:
func (_m *MockProgress) Start() {
	_m.Called()
}

// Stop provides a mock function with given fields:
func (_m *MockProgress) Stop() {
	_m.Called()
}

// Val provides a mock function with given fields:
func (_m *MockProgress) Val() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}
