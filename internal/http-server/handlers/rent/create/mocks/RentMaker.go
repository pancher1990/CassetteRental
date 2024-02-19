// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// RentMaker is an autogenerated mock type for the RentMaker type
type RentMaker struct {
	mock.Mock
}

// CreateCassetteInOrder provides a mock function with given fields: ctx, cassetteId, orderId, rentCost
func (_m *RentMaker) CreateCassetteInOrder(ctx context.Context, cassetteId string, orderId string, rentCost int) (context.Context, error) {
	ret := _m.Called(ctx, cassetteId, orderId, rentCost)

	var r0 context.Context
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int) (context.Context, error)); ok {
		return rf(ctx, cassetteId, orderId, rentCost)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int) context.Context); ok {
		r0 = rf(ctx, cassetteId, orderId, rentCost)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int) error); ok {
		r1 = rf(ctx, cassetteId, orderId, rentCost)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateOrder provides a mock function with given fields: ctx, customerId
func (_m *RentMaker) CreateOrder(ctx context.Context, customerId string) (context.Context, string, error) {
	ret := _m.Called(ctx, customerId)

	var r0 context.Context
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (context.Context, string, error)); ok {
		return rf(ctx, customerId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) context.Context); ok {
		r0 = rf(ctx, customerId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) string); ok {
		r1 = rf(ctx, customerId)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, customerId)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// CreateRent provides a mock function with given fields: ctx, customerId, cassetteId, rentDays
func (_m *RentMaker) CreateRent(ctx context.Context, customerId string, cassetteId string, rentDays int) (context.Context, string, error) {
	ret := _m.Called(ctx, customerId, cassetteId, rentDays)

	var r0 context.Context
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int) (context.Context, string, error)); ok {
		return rf(ctx, customerId, cassetteId, rentDays)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int) context.Context); ok {
		r0 = rf(ctx, customerId, cassetteId, rentDays)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int) string); ok {
		r1 = rf(ctx, customerId, cassetteId, rentDays)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, string, int) error); ok {
		r2 = rf(ctx, customerId, cassetteId, rentDays)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// FindAvailableCassette provides a mock function with given fields: ctx, filmId
func (_m *RentMaker) FindAvailableCassette(ctx context.Context, filmId string) (context.Context, string, error) {
	ret := _m.Called(ctx, filmId)

	var r0 context.Context
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (context.Context, string, error)); ok {
		return rf(ctx, filmId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) context.Context); ok {
		r0 = rf(ctx, filmId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) string); ok {
		r1 = rf(ctx, filmId)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, filmId)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetCustomerBalance provides a mock function with given fields: ctx, id
func (_m *RentMaker) GetCustomerBalance(ctx context.Context, id string) (context.Context, int, error) {
	ret := _m.Called(ctx, id)

	var r0 context.Context
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (context.Context, int, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) context.Context); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) int); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, id)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetFilm provides a mock function with given fields: ctx, title
func (_m *RentMaker) GetFilm(ctx context.Context, title string) (context.Context, string, int, error) {
	ret := _m.Called(ctx, title)

	var r0 context.Context
	var r1 string
	var r2 int
	var r3 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (context.Context, string, int, error)); ok {
		return rf(ctx, title)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) context.Context); ok {
		r0 = rf(ctx, title)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) string); ok {
		r1 = rf(ctx, title)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) int); ok {
		r2 = rf(ctx, title)
	} else {
		r2 = ret.Get(2).(int)
	}

	if rf, ok := ret.Get(3).(func(context.Context, string) error); ok {
		r3 = rf(ctx, title)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// SetCassetteStatus provides a mock function with given fields: ctx, id, available
func (_m *RentMaker) SetCassetteStatus(ctx context.Context, id string, available bool) (context.Context, error) {
	ret := _m.Called(ctx, id, available)

	var r0 context.Context
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) (context.Context, error)); ok {
		return rf(ctx, id, available)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) context.Context); ok {
		r0 = rf(ctx, id, available)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, bool) error); ok {
		r1 = rf(ctx, id, available)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetCustomerBalance provides a mock function with given fields: ctx, id, balance
func (_m *RentMaker) SetCustomerBalance(ctx context.Context, id string, balance int) (context.Context, error) {
	ret := _m.Called(ctx, id, balance)

	var r0 context.Context
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int) (context.Context, error)); ok {
		return rf(ctx, id, balance)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int) context.Context); ok {
		r0 = rf(ctx, id, balance)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int) error); ok {
		r1 = rf(ctx, id, balance)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRentMaker interface {
	mock.TestingT
	Cleanup(func())
}

// NewRentMaker creates a new instance of RentMaker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRentMaker(t mockConstructorTestingTNewRentMaker) *RentMaker {
	mock := &RentMaker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
