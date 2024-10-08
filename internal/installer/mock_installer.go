// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/openshift/image-based-install-operator/internal/installer (interfaces: Installer)
//
// Generated by this command:
//
//	mockgen --build_flags=--mod=mod -package=installer -destination=mock_installer.go . Installer
//

// Package installer is a generated GoMock package.
package installer

import (
	context "context"
	reflect "reflect"

	logrus "github.com/sirupsen/logrus"
	gomock "go.uber.org/mock/gomock"
)

// MockInstaller is a mock of Installer interface.
type MockInstaller struct {
	ctrl     *gomock.Controller
	recorder *MockInstallerMockRecorder
}

// MockInstallerMockRecorder is the mock recorder for MockInstaller.
type MockInstallerMockRecorder struct {
	mock *MockInstaller
}

// NewMockInstaller creates a new mock instance.
func NewMockInstaller(ctrl *gomock.Controller) *MockInstaller {
	mock := &MockInstaller{ctrl: ctrl}
	mock.recorder = &MockInstallerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInstaller) EXPECT() *MockInstallerMockRecorder {
	return m.recorder
}

// CreateInstallationIso mocks base method.
func (m *MockInstaller) CreateInstallationIso(arg0 context.Context, arg1 logrus.FieldLogger, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateInstallationIso", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateInstallationIso indicates an expected call of CreateInstallationIso.
func (mr *MockInstallerMockRecorder) CreateInstallationIso(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateInstallationIso", reflect.TypeOf((*MockInstaller)(nil).CreateInstallationIso), arg0, arg1, arg2)
}
