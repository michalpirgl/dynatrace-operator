// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	csivolumes "github.com/Dynatrace/dynatrace-operator/src/controllers/csi/driver/volumes"
	csi "github.com/container-storage-interface/spec/lib/go/csi"
	mock "github.com/stretchr/testify/mock"
)

// Publisher is an autogenerated mock type for the Publisher type
type Publisher struct {
	mock.Mock
}

type Publisher_Expecter struct {
	mock *mock.Mock
}

func (_m *Publisher) EXPECT() *Publisher_Expecter {
	return &Publisher_Expecter{mock: &_m.Mock}
}

// CanUnpublishVolume provides a mock function with given fields: ctx, volumeInfo
func (_m *Publisher) CanUnpublishVolume(ctx context.Context, volumeInfo *csivolumes.VolumeInfo) (bool, error) {
	ret := _m.Called(ctx, volumeInfo)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, *csivolumes.VolumeInfo) bool); ok {
		r0 = rf(ctx, volumeInfo)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *csivolumes.VolumeInfo) error); ok {
		r1 = rf(ctx, volumeInfo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Publisher_CanUnpublishVolume_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CanUnpublishVolume'
type Publisher_CanUnpublishVolume_Call struct {
	*mock.Call
}

// CanUnpublishVolume is a helper method to define mock.On call
//   - ctx context.Context
//   - volumeInfo *csivolumes.VolumeInfo
func (_e *Publisher_Expecter) CanUnpublishVolume(ctx interface{}, volumeInfo interface{}) *Publisher_CanUnpublishVolume_Call {
	return &Publisher_CanUnpublishVolume_Call{Call: _e.mock.On("CanUnpublishVolume", ctx, volumeInfo)}
}

func (_c *Publisher_CanUnpublishVolume_Call) Run(run func(ctx context.Context, volumeInfo *csivolumes.VolumeInfo)) *Publisher_CanUnpublishVolume_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*csivolumes.VolumeInfo))
	})
	return _c
}

func (_c *Publisher_CanUnpublishVolume_Call) Return(_a0 bool, _a1 error) *Publisher_CanUnpublishVolume_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// PublishVolume provides a mock function with given fields: ctx, volumeCfg
func (_m *Publisher) PublishVolume(ctx context.Context, volumeCfg *csivolumes.VolumeConfig) (*csi.NodePublishVolumeResponse, error) {
	ret := _m.Called(ctx, volumeCfg)

	var r0 *csi.NodePublishVolumeResponse
	if rf, ok := ret.Get(0).(func(context.Context, *csivolumes.VolumeConfig) *csi.NodePublishVolumeResponse); ok {
		r0 = rf(ctx, volumeCfg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*csi.NodePublishVolumeResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *csivolumes.VolumeConfig) error); ok {
		r1 = rf(ctx, volumeCfg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Publisher_PublishVolume_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishVolume'
type Publisher_PublishVolume_Call struct {
	*mock.Call
}

// PublishVolume is a helper method to define mock.On call
//   - ctx context.Context
//   - volumeCfg *csivolumes.VolumeConfig
func (_e *Publisher_Expecter) PublishVolume(ctx interface{}, volumeCfg interface{}) *Publisher_PublishVolume_Call {
	return &Publisher_PublishVolume_Call{Call: _e.mock.On("PublishVolume", ctx, volumeCfg)}
}

func (_c *Publisher_PublishVolume_Call) Run(run func(ctx context.Context, volumeCfg *csivolumes.VolumeConfig)) *Publisher_PublishVolume_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*csivolumes.VolumeConfig))
	})
	return _c
}

func (_c *Publisher_PublishVolume_Call) Return(_a0 *csi.NodePublishVolumeResponse, _a1 error) *Publisher_PublishVolume_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// UnpublishVolume provides a mock function with given fields: ctx, volumeInfo
func (_m *Publisher) UnpublishVolume(ctx context.Context, volumeInfo *csivolumes.VolumeInfo) (*csi.NodeUnpublishVolumeResponse, error) {
	ret := _m.Called(ctx, volumeInfo)

	var r0 *csi.NodeUnpublishVolumeResponse
	if rf, ok := ret.Get(0).(func(context.Context, *csivolumes.VolumeInfo) *csi.NodeUnpublishVolumeResponse); ok {
		r0 = rf(ctx, volumeInfo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*csi.NodeUnpublishVolumeResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *csivolumes.VolumeInfo) error); ok {
		r1 = rf(ctx, volumeInfo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Publisher_UnpublishVolume_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UnpublishVolume'
type Publisher_UnpublishVolume_Call struct {
	*mock.Call
}

// UnpublishVolume is a helper method to define mock.On call
//   - ctx context.Context
//   - volumeInfo *csivolumes.VolumeInfo
func (_e *Publisher_Expecter) UnpublishVolume(ctx interface{}, volumeInfo interface{}) *Publisher_UnpublishVolume_Call {
	return &Publisher_UnpublishVolume_Call{Call: _e.mock.On("UnpublishVolume", ctx, volumeInfo)}
}

func (_c *Publisher_UnpublishVolume_Call) Run(run func(ctx context.Context, volumeInfo *csivolumes.VolumeInfo)) *Publisher_UnpublishVolume_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*csivolumes.VolumeInfo))
	})
	return _c
}

func (_c *Publisher_UnpublishVolume_Call) Return(_a0 *csi.NodeUnpublishVolumeResponse, _a1 error) *Publisher_UnpublishVolume_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewPublisher interface {
	mock.TestingT
	Cleanup(func())
}

// NewPublisher creates a new instance of Publisher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPublisher(t mockConstructorTestingTNewPublisher) *Publisher {
	mock := &Publisher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
