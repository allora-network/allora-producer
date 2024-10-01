// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	proto "github.com/cosmos/gogoproto/proto"
	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/codec/types"
)

// ProtoCodec is an autogenerated mock type for the ProtoCodec type
type ProtoCodec struct {
	mock.Mock
}

type ProtoCodec_Expecter struct {
	mock *mock.Mock
}

func (_m *ProtoCodec) EXPECT() *ProtoCodec_Expecter {
	return &ProtoCodec_Expecter{mock: &_m.Mock}
}

// MarshalJSON provides a mock function with given fields: msg
func (_m *ProtoCodec) MarshalJSON(msg proto.Message) ([]byte, error) {
	ret := _m.Called(msg)

	if len(ret) == 0 {
		panic("no return value specified for MarshalJSON")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(proto.Message) ([]byte, error)); ok {
		return rf(msg)
	}
	if rf, ok := ret.Get(0).(func(proto.Message) []byte); ok {
		r0 = rf(msg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(proto.Message) error); ok {
		r1 = rf(msg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProtoCodec_MarshalJSON_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MarshalJSON'
type ProtoCodec_MarshalJSON_Call struct {
	*mock.Call
}

// MarshalJSON is a helper method to define mock.On call
//   - msg proto.Message
func (_e *ProtoCodec_Expecter) MarshalJSON(msg interface{}) *ProtoCodec_MarshalJSON_Call {
	return &ProtoCodec_MarshalJSON_Call{Call: _e.mock.On("MarshalJSON", msg)}
}

func (_c *ProtoCodec_MarshalJSON_Call) Run(run func(msg proto.Message)) *ProtoCodec_MarshalJSON_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(proto.Message))
	})
	return _c
}

func (_c *ProtoCodec_MarshalJSON_Call) Return(_a0 []byte, _a1 error) *ProtoCodec_MarshalJSON_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProtoCodec_MarshalJSON_Call) RunAndReturn(run func(proto.Message) ([]byte, error)) *ProtoCodec_MarshalJSON_Call {
	_c.Call.Return(run)
	return _c
}

// Unmarshal provides a mock function with given fields: data, msg
func (_m *ProtoCodec) Unmarshal(data []byte, msg proto.Message) error {
	ret := _m.Called(data, msg)

	if len(ret) == 0 {
		panic("no return value specified for Unmarshal")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, proto.Message) error); ok {
		r0 = rf(data, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ProtoCodec_Unmarshal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Unmarshal'
type ProtoCodec_Unmarshal_Call struct {
	*mock.Call
}

// Unmarshal is a helper method to define mock.On call
//   - data []byte
//   - msg proto.Message
func (_e *ProtoCodec_Expecter) Unmarshal(data interface{}, msg interface{}) *ProtoCodec_Unmarshal_Call {
	return &ProtoCodec_Unmarshal_Call{Call: _e.mock.On("Unmarshal", data, msg)}
}

func (_c *ProtoCodec_Unmarshal_Call) Run(run func(data []byte, msg proto.Message)) *ProtoCodec_Unmarshal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte), args[1].(proto.Message))
	})
	return _c
}

func (_c *ProtoCodec_Unmarshal_Call) Return(_a0 error) *ProtoCodec_Unmarshal_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ProtoCodec_Unmarshal_Call) RunAndReturn(run func([]byte, proto.Message) error) *ProtoCodec_Unmarshal_Call {
	_c.Call.Return(run)
	return _c
}

// UnpackAny provides a mock function with given fields: any, iface
func (_m *ProtoCodec) UnpackAny(any *types.Any, iface interface{}) error {
	ret := _m.Called(any, iface)

	if len(ret) == 0 {
		panic("no return value specified for UnpackAny")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*types.Any, interface{}) error); ok {
		r0 = rf(any, iface)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ProtoCodec_UnpackAny_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UnpackAny'
type ProtoCodec_UnpackAny_Call struct {
	*mock.Call
}

// UnpackAny is a helper method to define mock.On call
//   - any *types.Any
//   - iface interface{}
func (_e *ProtoCodec_Expecter) UnpackAny(any interface{}, iface interface{}) *ProtoCodec_UnpackAny_Call {
	return &ProtoCodec_UnpackAny_Call{Call: _e.mock.On("UnpackAny", any, iface)}
}

func (_c *ProtoCodec_UnpackAny_Call) Run(run func(any *types.Any, iface interface{})) *ProtoCodec_UnpackAny_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*types.Any), args[1].(interface{}))
	})
	return _c
}

func (_c *ProtoCodec_UnpackAny_Call) Return(_a0 error) *ProtoCodec_UnpackAny_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ProtoCodec_UnpackAny_Call) RunAndReturn(run func(*types.Any, interface{}) error) *ProtoCodec_UnpackAny_Call {
	_c.Call.Return(run)
	return _c
}

// NewProtoCodec creates a new instance of ProtoCodec. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProtoCodec(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProtoCodec {
	mock := &ProtoCodec{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}