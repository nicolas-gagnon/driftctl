// Code generated by mockery v2.10.0. DO NOT EDIT.

package repository

import (
	elb "github.com/aws/aws-sdk-go/service/elb"
	mock "github.com/stretchr/testify/mock"
)

// MockELBRepository is an autogenerated mock type for the ELBRepository type
type MockELBRepository struct {
	mock.Mock
}

// ListAllLoadBalancers provides a mock function with given fields:
func (_m *MockELBRepository) ListAllLoadBalancers() ([]*elb.LoadBalancerDescription, error) {
	ret := _m.Called()

	var r0 []*elb.LoadBalancerDescription
	if rf, ok := ret.Get(0).(func() []*elb.LoadBalancerDescription); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*elb.LoadBalancerDescription)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
