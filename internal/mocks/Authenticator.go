// Code generated by mockery v2.53.2. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// Authenticator is an autogenerated mock type for the Authenticator type
type Authenticator struct {
	mock.Mock
}

type Authenticator_Expecter struct {
	mock *mock.Mock
}

func (_m *Authenticator) EXPECT() *Authenticator_Expecter {
	return &Authenticator_Expecter{mock: &_m.Mock}
}

// CookieGetUserID provides a mock function with given fields: r, secret
func (_m *Authenticator) CookieGetUserID(r *http.Request, secret string) (string, error) {
	ret := _m.Called(r, secret)

	if len(ret) == 0 {
		panic("no return value specified for CookieGetUserID")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request, string) (string, error)); ok {
		return rf(r, secret)
	}
	if rf, ok := ret.Get(0).(func(*http.Request, string) string); ok {
		r0 = rf(r, secret)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*http.Request, string) error); ok {
		r1 = rf(r, secret)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Authenticator_CookieGetUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CookieGetUserID'
type Authenticator_CookieGetUserID_Call struct {
	*mock.Call
}

// CookieGetUserID is a helper method to define mock.On call
//   - r *http.Request
//   - secret string
func (_e *Authenticator_Expecter) CookieGetUserID(r interface{}, secret interface{}) *Authenticator_CookieGetUserID_Call {
	return &Authenticator_CookieGetUserID_Call{Call: _e.mock.On("CookieGetUserID", r, secret)}
}

func (_c *Authenticator_CookieGetUserID_Call) Run(run func(r *http.Request, secret string)) *Authenticator_CookieGetUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*http.Request), args[1].(string))
	})
	return _c
}

func (_c *Authenticator_CookieGetUserID_Call) Return(_a0 string, _a1 error) *Authenticator_CookieGetUserID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Authenticator_CookieGetUserID_Call) RunAndReturn(run func(*http.Request, string) (string, error)) *Authenticator_CookieGetUserID_Call {
	_c.Call.Return(run)
	return _c
}

// CookieSetUserID provides a mock function with given fields: w, secret
func (_m *Authenticator) CookieSetUserID(w http.ResponseWriter, secret string) (string, error) {
	ret := _m.Called(w, secret)

	if len(ret) == 0 {
		panic("no return value specified for CookieSetUserID")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, string) (string, error)); ok {
		return rf(w, secret)
	}
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, string) string); ok {
		r0 = rf(w, secret)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(http.ResponseWriter, string) error); ok {
		r1 = rf(w, secret)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Authenticator_CookieSetUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CookieSetUserID'
type Authenticator_CookieSetUserID_Call struct {
	*mock.Call
}

// CookieSetUserID is a helper method to define mock.On call
//   - w http.ResponseWriter
//   - secret string
func (_e *Authenticator_Expecter) CookieSetUserID(w interface{}, secret interface{}) *Authenticator_CookieSetUserID_Call {
	return &Authenticator_CookieSetUserID_Call{Call: _e.mock.On("CookieSetUserID", w, secret)}
}

func (_c *Authenticator_CookieSetUserID_Call) Run(run func(w http.ResponseWriter, secret string)) *Authenticator_CookieSetUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(http.ResponseWriter), args[1].(string))
	})
	return _c
}

func (_c *Authenticator_CookieSetUserID_Call) Return(userID string, err error) *Authenticator_CookieSetUserID_Call {
	_c.Call.Return(userID, err)
	return _c
}

func (_c *Authenticator_CookieSetUserID_Call) RunAndReturn(run func(http.ResponseWriter, string) (string, error)) *Authenticator_CookieSetUserID_Call {
	_c.Call.Return(run)
	return _c
}

// NewAuthenticator creates a new instance of Authenticator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthenticator(t interface {
	mock.TestingT
	Cleanup(func())
}) *Authenticator {
	mock := &Authenticator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
