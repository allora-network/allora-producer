// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	context "context"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"

	mock "github.com/stretchr/testify/mock"
)

// AlloraClientInterface is an autogenerated mock type for the AlloraClientInterface type
type AlloraClientInterface struct {
	mock.Mock
}

type AlloraClientInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *AlloraClientInterface) EXPECT() *AlloraClientInterface_Expecter {
	return &AlloraClientInterface_Expecter{mock: &_m.Mock}
}

// GetBlockByHeight provides a mock function with given fields: ctx, height
func (_m *AlloraClientInterface) GetBlockByHeight(ctx context.Context, height int64) (*coretypes.ResultBlock, error) {
	ret := _m.Called(ctx, height)

	if len(ret) == 0 {
		panic("no return value specified for GetBlockByHeight")
	}

	var r0 *coretypes.ResultBlock
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*coretypes.ResultBlock, error)); ok {
		return rf(ctx, height)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *coretypes.ResultBlock); ok {
		r0 = rf(ctx, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.ResultBlock)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AlloraClientInterface_GetBlockByHeight_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBlockByHeight'
type AlloraClientInterface_GetBlockByHeight_Call struct {
	*mock.Call
}

// GetBlockByHeight is a helper method to define mock.On call
//   - ctx context.Context
//   - height int64
func (_e *AlloraClientInterface_Expecter) GetBlockByHeight(ctx interface{}, height interface{}) *AlloraClientInterface_GetBlockByHeight_Call {
	return &AlloraClientInterface_GetBlockByHeight_Call{Call: _e.mock.On("GetBlockByHeight", ctx, height)}
}

func (_c *AlloraClientInterface_GetBlockByHeight_Call) Run(run func(ctx context.Context, height int64)) *AlloraClientInterface_GetBlockByHeight_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *AlloraClientInterface_GetBlockByHeight_Call) Return(_a0 *coretypes.ResultBlock, _a1 error) *AlloraClientInterface_GetBlockByHeight_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AlloraClientInterface_GetBlockByHeight_Call) RunAndReturn(run func(context.Context, int64) (*coretypes.ResultBlock, error)) *AlloraClientInterface_GetBlockByHeight_Call {
	_c.Call.Return(run)
	return _c
}

// GetBlockResults provides a mock function with given fields: ctx, height
func (_m *AlloraClientInterface) GetBlockResults(ctx context.Context, height int64) (*coretypes.ResultBlockResults, error) {
	ret := _m.Called(ctx, height)

	if len(ret) == 0 {
		panic("no return value specified for GetBlockResults")
	}

	var r0 *coretypes.ResultBlockResults
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*coretypes.ResultBlockResults, error)); ok {
		return rf(ctx, height)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *coretypes.ResultBlockResults); ok {
		r0 = rf(ctx, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.ResultBlockResults)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AlloraClientInterface_GetBlockResults_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBlockResults'
type AlloraClientInterface_GetBlockResults_Call struct {
	*mock.Call
}

// GetBlockResults is a helper method to define mock.On call
//   - ctx context.Context
//   - height int64
func (_e *AlloraClientInterface_Expecter) GetBlockResults(ctx interface{}, height interface{}) *AlloraClientInterface_GetBlockResults_Call {
	return &AlloraClientInterface_GetBlockResults_Call{Call: _e.mock.On("GetBlockResults", ctx, height)}
}

func (_c *AlloraClientInterface_GetBlockResults_Call) Run(run func(ctx context.Context, height int64)) *AlloraClientInterface_GetBlockResults_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *AlloraClientInterface_GetBlockResults_Call) Return(_a0 *coretypes.ResultBlockResults, _a1 error) *AlloraClientInterface_GetBlockResults_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AlloraClientInterface_GetBlockResults_Call) RunAndReturn(run func(context.Context, int64) (*coretypes.ResultBlockResults, error)) *AlloraClientInterface_GetBlockResults_Call {
	_c.Call.Return(run)
	return _c
}

// GetHeader provides a mock function with given fields: ctx, height
func (_m *AlloraClientInterface) GetHeader(ctx context.Context, height int64) (*coretypes.ResultHeader, error) {
	ret := _m.Called(ctx, height)

	if len(ret) == 0 {
		panic("no return value specified for GetHeader")
	}

	var r0 *coretypes.ResultHeader
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*coretypes.ResultHeader, error)); ok {
		return rf(ctx, height)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *coretypes.ResultHeader); ok {
		r0 = rf(ctx, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*coretypes.ResultHeader)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AlloraClientInterface_GetHeader_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetHeader'
type AlloraClientInterface_GetHeader_Call struct {
	*mock.Call
}

// GetHeader is a helper method to define mock.On call
//   - ctx context.Context
//   - height int64
func (_e *AlloraClientInterface_Expecter) GetHeader(ctx interface{}, height interface{}) *AlloraClientInterface_GetHeader_Call {
	return &AlloraClientInterface_GetHeader_Call{Call: _e.mock.On("GetHeader", ctx, height)}
}

func (_c *AlloraClientInterface_GetHeader_Call) Run(run func(ctx context.Context, height int64)) *AlloraClientInterface_GetHeader_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *AlloraClientInterface_GetHeader_Call) Return(_a0 *coretypes.ResultHeader, _a1 error) *AlloraClientInterface_GetHeader_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AlloraClientInterface_GetHeader_Call) RunAndReturn(run func(context.Context, int64) (*coretypes.ResultHeader, error)) *AlloraClientInterface_GetHeader_Call {
	_c.Call.Return(run)
	return _c
}

// GetLatestBlockHeight provides a mock function with given fields: ctx
func (_m *AlloraClientInterface) GetLatestBlockHeight(ctx context.Context) (int64, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetLatestBlockHeight")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (int64, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) int64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AlloraClientInterface_GetLatestBlockHeight_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLatestBlockHeight'
type AlloraClientInterface_GetLatestBlockHeight_Call struct {
	*mock.Call
}

// GetLatestBlockHeight is a helper method to define mock.On call
//   - ctx context.Context
func (_e *AlloraClientInterface_Expecter) GetLatestBlockHeight(ctx interface{}) *AlloraClientInterface_GetLatestBlockHeight_Call {
	return &AlloraClientInterface_GetLatestBlockHeight_Call{Call: _e.mock.On("GetLatestBlockHeight", ctx)}
}

func (_c *AlloraClientInterface_GetLatestBlockHeight_Call) Run(run func(ctx context.Context)) *AlloraClientInterface_GetLatestBlockHeight_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *AlloraClientInterface_GetLatestBlockHeight_Call) Return(_a0 int64, _a1 error) *AlloraClientInterface_GetLatestBlockHeight_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AlloraClientInterface_GetLatestBlockHeight_Call) RunAndReturn(run func(context.Context) (int64, error)) *AlloraClientInterface_GetLatestBlockHeight_Call {
	_c.Call.Return(run)
	return _c
}

// NewAlloraClientInterface creates a new instance of AlloraClientInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAlloraClientInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *AlloraClientInterface {
	mock := &AlloraClientInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}