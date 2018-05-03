// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ONSdigital/dp-observation-extractor/event (interfaces: Client,VaultClient)

// Package mock_event is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	s3 "github.com/aws/aws-sdk-go/service/s3"
	gomock "github.com/golang/mock/gomock"
)

// MockSDKClient is a mock of SDKClient interface
type MockSDKClient struct {
	ctrl     *gomock.Controller
	recorder *MockSDKClientMockRecorder
}

// MockSDKClientMockRecorder is the mock recorder for MockSDKClient
type MockSDKClientMockRecorder struct {
	mock *MockSDKClient
}

// NewMockSDKClient creates a new mock instance
func NewMockSDKClient(ctrl *gomock.Controller) *MockSDKClient {
	mock := &MockSDKClient{ctrl: ctrl}
	mock.recorder = &MockSDKClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSDKClient) EXPECT() *MockSDKClientMockRecorder {
	return m.recorder
}

// GetObject mocks base method
func (m *MockSDKClient) GetObject(arg0 *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	ret := m.ctrl.Call(m, "GetObject", arg0)
	ret0, _ := ret[0].(*s3.GetObjectOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObject indicates an expected call of GetObject
func (mr *MockSDKClientMockRecorder) GetObject(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObject", reflect.TypeOf((*MockSDKClient)(nil).GetObject), arg0)
}

// MockCryptoClient is a mock of CryptoClient interface
type MockCryptoClient struct {
	ctrl     *gomock.Controller
	recorder *MockCryptoClientMockRecorder
}

// MockCryptoClientMockRecorder is the mock recorder for MockCryptoClient
type MockCryptoClientMockRecorder struct {
	mock *MockCryptoClient
}

// NewMockCryptoClient creates a new mock instance
func NewMockCryptoClient(ctrl *gomock.Controller) *MockCryptoClient {
	mock := &MockCryptoClient{ctrl: ctrl}
	mock.recorder = &MockCryptoClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCryptoClient) EXPECT() *MockCryptoClientMockRecorder {
	return m.recorder
}

// GetObjectWithPSK mocks base method
func (m *MockCryptoClient) GetObjectWithPSK(arg0 *s3.GetObjectInput, arg1 []byte) (*s3.GetObjectOutput, error) {
	ret := m.ctrl.Call(m, "GetObjectWithPSK", arg0, arg1)
	ret0, _ := ret[0].(*s3.GetObjectOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObjectWithPSK indicates an expected call of GetObjectWithPSK
func (mr *MockCryptoClientMockRecorder) GetObjectWithPSK(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObjectWithPSK", reflect.TypeOf((*MockCryptoClient)(nil).GetObjectWithPSK), arg0, arg1)
}

// MockVaultClient is a mock of VaultClient interface
type MockVaultClient struct {
	ctrl     *gomock.Controller
	recorder *MockVaultClientMockRecorder
}

// MockVaultClientMockRecorder is the mock recorder for MockVaultClient
type MockVaultClientMockRecorder struct {
	mock *MockVaultClient
}

// NewMockVaultClient creates a new mock instance
func NewMockVaultClient(ctrl *gomock.Controller) *MockVaultClient {
	mock := &MockVaultClient{ctrl: ctrl}
	mock.recorder = &MockVaultClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockVaultClient) EXPECT() *MockVaultClientMockRecorder {
	return m.recorder
}

// ReadKey mocks base method
func (m *MockVaultClient) ReadKey(arg0, arg1 string) (string, error) {
	ret := m.ctrl.Call(m, "ReadKey", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadKey indicates an expected call of ReadKey
func (mr *MockVaultClientMockRecorder) ReadKey(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadKey", reflect.TypeOf((*MockVaultClient)(nil).ReadKey), arg0, arg1)
}
