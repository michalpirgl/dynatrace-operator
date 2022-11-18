// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	io "io"

	dtclient "github.com/Dynatrace/dynatrace-operator/src/dtclient"

	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

type Client_Expecter struct {
	mock *mock.Mock
}

func (_m *Client) EXPECT() *Client_Expecter {
	return &Client_Expecter{mock: &_m.Mock}
}

// CreateOrUpdateKubernetesSetting provides a mock function with given fields: name, kubeSystemUUID, scope
func (_m *Client) CreateOrUpdateKubernetesSetting(name string, kubeSystemUUID string, scope string) (string, error) {
	ret := _m.Called(name, kubeSystemUUID, scope)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string) string); ok {
		r0 = rf(name, kubeSystemUUID, scope)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string) error); ok {
		r1 = rf(name, kubeSystemUUID, scope)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_CreateOrUpdateKubernetesSetting_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateOrUpdateKubernetesSetting'
type Client_CreateOrUpdateKubernetesSetting_Call struct {
	*mock.Call
}

// CreateOrUpdateKubernetesSetting is a helper method to define mock.On call
//   - name string
//   - kubeSystemUUID string
//   - scope string
func (_e *Client_Expecter) CreateOrUpdateKubernetesSetting(name interface{}, kubeSystemUUID interface{}, scope interface{}) *Client_CreateOrUpdateKubernetesSetting_Call {
	return &Client_CreateOrUpdateKubernetesSetting_Call{Call: _e.mock.On("CreateOrUpdateKubernetesSetting", name, kubeSystemUUID, scope)}
}

func (_c *Client_CreateOrUpdateKubernetesSetting_Call) Run(run func(name string, kubeSystemUUID string, scope string)) *Client_CreateOrUpdateKubernetesSetting_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *Client_CreateOrUpdateKubernetesSetting_Call) Return(_a0 string, _a1 error) *Client_CreateOrUpdateKubernetesSetting_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetActiveGateAuthToken provides a mock function with given fields: dynakubeName
func (_m *Client) GetActiveGateAuthToken(dynakubeName string) (*dtclient.ActiveGateAuthTokenInfo, error) {
	ret := _m.Called(dynakubeName)

	var r0 *dtclient.ActiveGateAuthTokenInfo
	if rf, ok := ret.Get(0).(func(string) *dtclient.ActiveGateAuthTokenInfo); ok {
		r0 = rf(dynakubeName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dtclient.ActiveGateAuthTokenInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(dynakubeName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_GetActiveGateAuthToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetActiveGateAuthToken'
type Client_GetActiveGateAuthToken_Call struct {
	*mock.Call
}

// GetActiveGateAuthToken is a helper method to define mock.On call
//   - dynakubeName string
func (_e *Client_Expecter) GetActiveGateAuthToken(dynakubeName interface{}) *Client_GetActiveGateAuthToken_Call {
	return &Client_GetActiveGateAuthToken_Call{Call: _e.mock.On("GetActiveGateAuthToken", dynakubeName)}
}

func (_c *Client_GetActiveGateAuthToken_Call) Run(run func(dynakubeName string)) *Client_GetActiveGateAuthToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Client_GetActiveGateAuthToken_Call) Return(_a0 *dtclient.ActiveGateAuthTokenInfo, _a1 error) *Client_GetActiveGateAuthToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetActiveGateConnectionInfo provides a mock function with given fields:
func (_m *Client) GetActiveGateConnectionInfo() (*dtclient.ActiveGateConnectionInfo, error) {
	ret := _m.Called()

	var r0 *dtclient.ActiveGateConnectionInfo
	if rf, ok := ret.Get(0).(func() *dtclient.ActiveGateConnectionInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dtclient.ActiveGateConnectionInfo)
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

// Client_GetActiveGateConnectionInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetActiveGateConnectionInfo'
type Client_GetActiveGateConnectionInfo_Call struct {
	*mock.Call
}

// GetActiveGateConnectionInfo is a helper method to define mock.On call
func (_e *Client_Expecter) GetActiveGateConnectionInfo() *Client_GetActiveGateConnectionInfo_Call {
	return &Client_GetActiveGateConnectionInfo_Call{Call: _e.mock.On("GetActiveGateConnectionInfo")}
}

func (_c *Client_GetActiveGateConnectionInfo_Call) Run(run func()) *Client_GetActiveGateConnectionInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Client_GetActiveGateConnectionInfo_Call) Return(_a0 *dtclient.ActiveGateConnectionInfo, _a1 error) *Client_GetActiveGateConnectionInfo_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetAgent provides a mock function with given fields: os, installerType, flavor, arch, version, technologies, writer
func (_m *Client) GetAgent(os string, installerType string, flavor string, arch string, version string, technologies []string, writer io.Writer) error {
	ret := _m.Called(os, installerType, flavor, arch, version, technologies, writer)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string, string, []string, io.Writer) error); ok {
		r0 = rf(os, installerType, flavor, arch, version, technologies, writer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Client_GetAgent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAgent'
type Client_GetAgent_Call struct {
	*mock.Call
}

// GetAgent is a helper method to define mock.On call
//   - os string
//   - installerType string
//   - flavor string
//   - arch string
//   - version string
//   - technologies []string
//   - writer io.Writer
func (_e *Client_Expecter) GetAgent(os interface{}, installerType interface{}, flavor interface{}, arch interface{}, version interface{}, technologies interface{}, writer interface{}) *Client_GetAgent_Call {
	return &Client_GetAgent_Call{Call: _e.mock.On("GetAgent", os, installerType, flavor, arch, version, technologies, writer)}
}

func (_c *Client_GetAgent_Call) Run(run func(os string, installerType string, flavor string, arch string, version string, technologies []string, writer io.Writer)) *Client_GetAgent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string), args[3].(string), args[4].(string), args[5].([]string), args[6].(io.Writer))
	})
	return _c
}

func (_c *Client_GetAgent_Call) Return(_a0 error) *Client_GetAgent_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetAgentVersions provides a mock function with given fields: os, installerType, flavor, arch
func (_m *Client) GetAgentVersions(os string, installerType string, flavor string, arch string) ([]string, error) {
	ret := _m.Called(os, installerType, flavor, arch)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string, string, string, string) []string); ok {
		r0 = rf(os, installerType, flavor, arch)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, string) error); ok {
		r1 = rf(os, installerType, flavor, arch)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_GetAgentVersions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAgentVersions'
type Client_GetAgentVersions_Call struct {
	*mock.Call
}

// GetAgentVersions is a helper method to define mock.On call
//   - os string
//   - installerType string
//   - flavor string
//   - arch string
func (_e *Client_Expecter) GetAgentVersions(os interface{}, installerType interface{}, flavor interface{}, arch interface{}) *Client_GetAgentVersions_Call {
	return &Client_GetAgentVersions_Call{Call: _e.mock.On("GetAgentVersions", os, installerType, flavor, arch)}
}

func (_c *Client_GetAgentVersions_Call) Run(run func(os string, installerType string, flavor string, arch string)) *Client_GetAgentVersions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string), args[3].(string))
	})
	return _c
}

func (_c *Client_GetAgentVersions_Call) Return(_a0 []string, _a1 error) *Client_GetAgentVersions_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetAgentViaInstallerUrl provides a mock function with given fields: url, writer
func (_m *Client) GetAgentViaInstallerUrl(url string, writer io.Writer) error {
	ret := _m.Called(url, writer)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, io.Writer) error); ok {
		r0 = rf(url, writer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Client_GetAgentViaInstallerUrl_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAgentViaInstallerUrl'
type Client_GetAgentViaInstallerUrl_Call struct {
	*mock.Call
}

// GetAgentViaInstallerUrl is a helper method to define mock.On call
//   - url string
//   - writer io.Writer
func (_e *Client_Expecter) GetAgentViaInstallerUrl(url interface{}, writer interface{}) *Client_GetAgentViaInstallerUrl_Call {
	return &Client_GetAgentViaInstallerUrl_Call{Call: _e.mock.On("GetAgentViaInstallerUrl", url, writer)}
}

func (_c *Client_GetAgentViaInstallerUrl_Call) Run(run func(url string, writer io.Writer)) *Client_GetAgentViaInstallerUrl_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(io.Writer))
	})
	return _c
}

func (_c *Client_GetAgentViaInstallerUrl_Call) Return(_a0 error) *Client_GetAgentViaInstallerUrl_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetCommunicationHostForClient provides a mock function with given fields:
func (_m *Client) GetCommunicationHostForClient() (dtclient.CommunicationHost, error) {
	ret := _m.Called()

	var r0 dtclient.CommunicationHost
	if rf, ok := ret.Get(0).(func() dtclient.CommunicationHost); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(dtclient.CommunicationHost)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_GetCommunicationHostForClient_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCommunicationHostForClient'
type Client_GetCommunicationHostForClient_Call struct {
	*mock.Call
}

// GetCommunicationHostForClient is a helper method to define mock.On call
func (_e *Client_Expecter) GetCommunicationHostForClient() *Client_GetCommunicationHostForClient_Call {
	return &Client_GetCommunicationHostForClient_Call{Call: _e.mock.On("GetCommunicationHostForClient")}
}

func (_c *Client_GetCommunicationHostForClient_Call) Run(run func()) *Client_GetCommunicationHostForClient_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Client_GetCommunicationHostForClient_Call) Return(_a0 dtclient.CommunicationHost, _a1 error) *Client_GetCommunicationHostForClient_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetEntityIDForIP provides a mock function with given fields: ip
func (_m *Client) GetEntityIDForIP(ip string) (string, error) {
	ret := _m.Called(ip)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(ip)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(ip)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_GetEntityIDForIP_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetEntityIDForIP'
type Client_GetEntityIDForIP_Call struct {
	*mock.Call
}

// GetEntityIDForIP is a helper method to define mock.On call
//   - ip string
func (_e *Client_Expecter) GetEntityIDForIP(ip interface{}) *Client_GetEntityIDForIP_Call {
	return &Client_GetEntityIDForIP_Call{Call: _e.mock.On("GetEntityIDForIP", ip)}
}

func (_c *Client_GetEntityIDForIP_Call) Run(run func(ip string)) *Client_GetEntityIDForIP_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Client_GetEntityIDForIP_Call) Return(_a0 string, _a1 error) *Client_GetEntityIDForIP_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetLatestAgent provides a mock function with given fields: os, installerType, flavor, arch, technologies, writer
func (_m *Client) GetLatestAgent(os string, installerType string, flavor string, arch string, technologies []string, writer io.Writer) error {
	ret := _m.Called(os, installerType, flavor, arch, technologies, writer)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string, []string, io.Writer) error); ok {
		r0 = rf(os, installerType, flavor, arch, technologies, writer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Client_GetLatestAgent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLatestAgent'
type Client_GetLatestAgent_Call struct {
	*mock.Call
}

// GetLatestAgent is a helper method to define mock.On call
//   - os string
//   - installerType string
//   - flavor string
//   - arch string
//   - technologies []string
//   - writer io.Writer
func (_e *Client_Expecter) GetLatestAgent(os interface{}, installerType interface{}, flavor interface{}, arch interface{}, technologies interface{}, writer interface{}) *Client_GetLatestAgent_Call {
	return &Client_GetLatestAgent_Call{Call: _e.mock.On("GetLatestAgent", os, installerType, flavor, arch, technologies, writer)}
}

func (_c *Client_GetLatestAgent_Call) Run(run func(os string, installerType string, flavor string, arch string, technologies []string, writer io.Writer)) *Client_GetLatestAgent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(string), args[3].(string), args[4].([]string), args[5].(io.Writer))
	})
	return _c
}

func (_c *Client_GetLatestAgent_Call) Return(_a0 error) *Client_GetLatestAgent_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetLatestAgentVersion provides a mock function with given fields: os, installerType
func (_m *Client) GetLatestAgentVersion(os string, installerType string) (string, error) {
	ret := _m.Called(os, installerType)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(os, installerType)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(os, installerType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_GetLatestAgentVersion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLatestAgentVersion'
type Client_GetLatestAgentVersion_Call struct {
	*mock.Call
}

// GetLatestAgentVersion is a helper method to define mock.On call
//   - os string
//   - installerType string
func (_e *Client_Expecter) GetLatestAgentVersion(os interface{}, installerType interface{}) *Client_GetLatestAgentVersion_Call {
	return &Client_GetLatestAgentVersion_Call{Call: _e.mock.On("GetLatestAgentVersion", os, installerType)}
}

func (_c *Client_GetLatestAgentVersion_Call) Run(run func(os string, installerType string)) *Client_GetLatestAgentVersion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *Client_GetLatestAgentVersion_Call) Return(_a0 string, _a1 error) *Client_GetLatestAgentVersion_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetMonitoredEntitiesForKubeSystemUUID provides a mock function with given fields: kubeSystemUUID
func (_m *Client) GetMonitoredEntitiesForKubeSystemUUID(kubeSystemUUID string) ([]dtclient.MonitoredEntity, error) {
	ret := _m.Called(kubeSystemUUID)

	var r0 []dtclient.MonitoredEntity
	if rf, ok := ret.Get(0).(func(string) []dtclient.MonitoredEntity); ok {
		r0 = rf(kubeSystemUUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dtclient.MonitoredEntity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(kubeSystemUUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_GetMonitoredEntitiesForKubeSystemUUID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMonitoredEntitiesForKubeSystemUUID'
type Client_GetMonitoredEntitiesForKubeSystemUUID_Call struct {
	*mock.Call
}

// GetMonitoredEntitiesForKubeSystemUUID is a helper method to define mock.On call
//   - kubeSystemUUID string
func (_e *Client_Expecter) GetMonitoredEntitiesForKubeSystemUUID(kubeSystemUUID interface{}) *Client_GetMonitoredEntitiesForKubeSystemUUID_Call {
	return &Client_GetMonitoredEntitiesForKubeSystemUUID_Call{Call: _e.mock.On("GetMonitoredEntitiesForKubeSystemUUID", kubeSystemUUID)}
}

func (_c *Client_GetMonitoredEntitiesForKubeSystemUUID_Call) Run(run func(kubeSystemUUID string)) *Client_GetMonitoredEntitiesForKubeSystemUUID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Client_GetMonitoredEntitiesForKubeSystemUUID_Call) Return(_a0 []dtclient.MonitoredEntity, _a1 error) *Client_GetMonitoredEntitiesForKubeSystemUUID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetOneAgentConnectionInfo provides a mock function with given fields:
func (_m *Client) GetOneAgentConnectionInfo() (dtclient.OneAgentConnectionInfo, error) {
	ret := _m.Called()

	var r0 dtclient.OneAgentConnectionInfo
	if rf, ok := ret.Get(0).(func() dtclient.OneAgentConnectionInfo); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(dtclient.OneAgentConnectionInfo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_GetOneAgentConnectionInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOneAgentConnectionInfo'
type Client_GetOneAgentConnectionInfo_Call struct {
	*mock.Call
}

// GetOneAgentConnectionInfo is a helper method to define mock.On call
func (_e *Client_Expecter) GetOneAgentConnectionInfo() *Client_GetOneAgentConnectionInfo_Call {
	return &Client_GetOneAgentConnectionInfo_Call{Call: _e.mock.On("GetOneAgentConnectionInfo")}
}

func (_c *Client_GetOneAgentConnectionInfo_Call) Run(run func()) *Client_GetOneAgentConnectionInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Client_GetOneAgentConnectionInfo_Call) Return(_a0 dtclient.OneAgentConnectionInfo, _a1 error) *Client_GetOneAgentConnectionInfo_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetProcessModuleConfig provides a mock function with given fields: prevRevision
func (_m *Client) GetProcessModuleConfig(prevRevision uint) (*dtclient.ProcessModuleConfig, error) {
	ret := _m.Called(prevRevision)

	var r0 *dtclient.ProcessModuleConfig
	if rf, ok := ret.Get(0).(func(uint) *dtclient.ProcessModuleConfig); ok {
		r0 = rf(prevRevision)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dtclient.ProcessModuleConfig)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(prevRevision)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_GetProcessModuleConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProcessModuleConfig'
type Client_GetProcessModuleConfig_Call struct {
	*mock.Call
}

// GetProcessModuleConfig is a helper method to define mock.On call
//   - prevRevision uint
func (_e *Client_Expecter) GetProcessModuleConfig(prevRevision interface{}) *Client_GetProcessModuleConfig_Call {
	return &Client_GetProcessModuleConfig_Call{Call: _e.mock.On("GetProcessModuleConfig", prevRevision)}
}

func (_c *Client_GetProcessModuleConfig_Call) Run(run func(prevRevision uint)) *Client_GetProcessModuleConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint))
	})
	return _c
}

func (_c *Client_GetProcessModuleConfig_Call) Return(_a0 *dtclient.ProcessModuleConfig, _a1 error) *Client_GetProcessModuleConfig_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetSettingsForMonitoredEntities provides a mock function with given fields: monitoredEntities
func (_m *Client) GetSettingsForMonitoredEntities(monitoredEntities []dtclient.MonitoredEntity) (dtclient.GetSettingsResponse, error) {
	ret := _m.Called(monitoredEntities)

	var r0 dtclient.GetSettingsResponse
	if rf, ok := ret.Get(0).(func([]dtclient.MonitoredEntity) dtclient.GetSettingsResponse); ok {
		r0 = rf(monitoredEntities)
	} else {
		r0 = ret.Get(0).(dtclient.GetSettingsResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]dtclient.MonitoredEntity) error); ok {
		r1 = rf(monitoredEntities)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_GetSettingsForMonitoredEntities_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSettingsForMonitoredEntities'
type Client_GetSettingsForMonitoredEntities_Call struct {
	*mock.Call
}

// GetSettingsForMonitoredEntities is a helper method to define mock.On call
//   - monitoredEntities []dtclient.MonitoredEntity
func (_e *Client_Expecter) GetSettingsForMonitoredEntities(monitoredEntities interface{}) *Client_GetSettingsForMonitoredEntities_Call {
	return &Client_GetSettingsForMonitoredEntities_Call{Call: _e.mock.On("GetSettingsForMonitoredEntities", monitoredEntities)}
}

func (_c *Client_GetSettingsForMonitoredEntities_Call) Run(run func(monitoredEntities []dtclient.MonitoredEntity)) *Client_GetSettingsForMonitoredEntities_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]dtclient.MonitoredEntity))
	})
	return _c
}

func (_c *Client_GetSettingsForMonitoredEntities_Call) Return(_a0 dtclient.GetSettingsResponse, _a1 error) *Client_GetSettingsForMonitoredEntities_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetTokenScopes provides a mock function with given fields: token
func (_m *Client) GetTokenScopes(token string) (dtclient.TokenScopes, error) {
	ret := _m.Called(token)

	var r0 dtclient.TokenScopes
	if rf, ok := ret.Get(0).(func(string) dtclient.TokenScopes); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(dtclient.TokenScopes)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_GetTokenScopes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTokenScopes'
type Client_GetTokenScopes_Call struct {
	*mock.Call
}

// GetTokenScopes is a helper method to define mock.On call
//   - token string
func (_e *Client_Expecter) GetTokenScopes(token interface{}) *Client_GetTokenScopes_Call {
	return &Client_GetTokenScopes_Call{Call: _e.mock.On("GetTokenScopes", token)}
}

func (_c *Client_GetTokenScopes_Call) Run(run func(token string)) *Client_GetTokenScopes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Client_GetTokenScopes_Call) Return(_a0 dtclient.TokenScopes, _a1 error) *Client_GetTokenScopes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// SendEvent provides a mock function with given fields: eventData
func (_m *Client) SendEvent(eventData *dtclient.EventData) error {
	ret := _m.Called(eventData)

	var r0 error
	if rf, ok := ret.Get(0).(func(*dtclient.EventData) error); ok {
		r0 = rf(eventData)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Client_SendEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendEvent'
type Client_SendEvent_Call struct {
	*mock.Call
}

// SendEvent is a helper method to define mock.On call
//   - eventData *dtclient.EventData
func (_e *Client_Expecter) SendEvent(eventData interface{}) *Client_SendEvent_Call {
	return &Client_SendEvent_Call{Call: _e.mock.On("SendEvent", eventData)}
}

func (_c *Client_SendEvent_Call) Run(run func(eventData *dtclient.EventData)) *Client_SendEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*dtclient.EventData))
	})
	return _c
}

func (_c *Client_SendEvent_Call) Return(_a0 error) *Client_SendEvent_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewClient creates a new instance of Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewClient(t mockConstructorTestingTNewClient) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
