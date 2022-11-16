// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	builder "github.com/Dynatrace/dynatrace-operator/src/builder"
	mock "github.com/stretchr/testify/mock"
)

// Builder is an autogenerated mock type for the Builder type
type Builder[T interface{}] struct {
	mock.Mock
}

type Builder_Expecter[T interface{}] struct {
	mock *mock.Mock
}

func (_m *Builder[T]) EXPECT() *Builder_Expecter[T] {
	return &Builder_Expecter[T]{mock: &_m.Mock}
}

// AddModifier provides a mock function with given fields: _a0
func (_m *Builder[T]) AddModifier(_a0 ...builder.Modifier[T]) builder.Builder[T] {
	_va := make([]interface{}, len(_a0))
	for _i := range _a0 {
		_va[_i] = _a0[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 builder.Builder[T]
	if rf, ok := ret.Get(0).(func(...builder.Modifier[T]) builder.Builder[T]); ok {
		r0 = rf(_a0...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(builder.Builder[T])
		}
	}

	return r0
}

// Builder_AddModifier_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddModifier'
type Builder_AddModifier_Call[T interface{}] struct {
	*mock.Call
}

// AddModifier is a helper method to define mock.On call
//   - _a0 ...builder.Modifier[T]
func (_e *Builder_Expecter[T]) AddModifier(_a0 ...interface{}) *Builder_AddModifier_Call[T] {
	return &Builder_AddModifier_Call[T]{Call: _e.mock.On("AddModifier",
		append([]interface{}{}, _a0...)...)}
}

func (_c *Builder_AddModifier_Call[T]) Run(run func(_a0 ...builder.Modifier[T])) *Builder_AddModifier_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]builder.Modifier[T], len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(builder.Modifier[T])
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *Builder_AddModifier_Call[T]) Return(_a0 builder.Builder[T]) *Builder_AddModifier_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

// Build provides a mock function with given fields:
func (_m *Builder[T]) Build() T {
	ret := _m.Called()

	var r0 T
	if rf, ok := ret.Get(0).(func() T); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(T)
	}

	return r0
}

// Builder_Build_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Build'
type Builder_Build_Call[T interface{}] struct {
	*mock.Call
}

// Build is a helper method to define mock.On call
func (_e *Builder_Expecter[T]) Build() *Builder_Build_Call[T] {
	return &Builder_Build_Call[T]{Call: _e.mock.On("Build")}
}

func (_c *Builder_Build_Call[T]) Run(run func()) *Builder_Build_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Builder_Build_Call[T]) Return(_a0 T) *Builder_Build_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewBuilder interface {
	mock.TestingT
	Cleanup(func())
}

// NewBuilder creates a new instance of Builder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBuilder[T interface{}](t mockConstructorTestingTNewBuilder) *Builder[T] {
	mock := &Builder[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
