// Code generated by MockGen. DO NOT EDIT.
// Source: logger.go

// Package mocks is a generated GoMock package.
package mocks

import (
	pb "github.com/GoogleCloudPlatform/compute-image-tools/proto/go/pb"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockLogWriter is a mock of LogWriter interface
type MockLogWriter struct {
	ctrl     *gomock.Controller
	recorder *MockLogWriterMockRecorder
}

// MockLogWriterMockRecorder is the mock recorder for MockLogWriter
type MockLogWriterMockRecorder struct {
	mock *MockLogWriter
}

// NewMockLogWriter creates a new mock instance
func NewMockLogWriter(ctrl *gomock.Controller) *MockLogWriter {
	mock := &MockLogWriter{ctrl: ctrl}
	mock.recorder = &MockLogWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogWriter) EXPECT() *MockLogWriterMockRecorder {
	return m.recorder
}

// WriteUser mocks base method
func (m *MockLogWriter) WriteUser(message string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteUser", message)
}

// WriteUser indicates an expected call of WriteUser
func (mr *MockLogWriterMockRecorder) WriteUser(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteUser", reflect.TypeOf((*MockLogWriter)(nil).WriteUser), message)
}

// WriteDebug mocks base method
func (m *MockLogWriter) WriteDebug(message string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteDebug", message)
}

// WriteDebug indicates an expected call of WriteDebug
func (mr *MockLogWriterMockRecorder) WriteDebug(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteDebug", reflect.TypeOf((*MockLogWriter)(nil).WriteDebug), message)
}

// WriteTrace mocks base method
func (m *MockLogWriter) WriteTrace(message string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteTrace", message)
}

// WriteTrace indicates an expected call of WriteTrace
func (mr *MockLogWriterMockRecorder) WriteTrace(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteTrace", reflect.TypeOf((*MockLogWriter)(nil).WriteTrace), message)
}

// WriteMetric mocks base method
func (m *MockLogWriter) WriteMetric(metric *pb.OutputInfo) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteMetric", metric)
}

// WriteMetric indicates an expected call of WriteMetric
func (mr *MockLogWriterMockRecorder) WriteMetric(metric interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteMetric", reflect.TypeOf((*MockLogWriter)(nil).WriteMetric), metric)
}

// MockLogReader is a mock of LogReader interface
type MockLogReader struct {
	ctrl     *gomock.Controller
	recorder *MockLogReaderMockRecorder
}

// MockLogReaderMockRecorder is the mock recorder for MockLogReader
type MockLogReaderMockRecorder struct {
	mock *MockLogReader
}

// NewMockLogReader creates a new mock instance
func NewMockLogReader(ctrl *gomock.Controller) *MockLogReader {
	mock := &MockLogReader{ctrl: ctrl}
	mock.recorder = &MockLogReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogReader) EXPECT() *MockLogReaderMockRecorder {
	return m.recorder
}

// ReadOutputInfo mocks base method
func (m *MockLogReader) ReadOutputInfo() *pb.OutputInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadOutputInfo")
	ret0, _ := ret[0].(*pb.OutputInfo)
	return ret0
}

// ReadOutputInfo indicates an expected call of ReadOutputInfo
func (mr *MockLogReaderMockRecorder) ReadOutputInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadOutputInfo", reflect.TypeOf((*MockLogReader)(nil).ReadOutputInfo))
}

// MockToolLogger is a mock of ToolLogger interface
type MockToolLogger struct {
	ctrl     *gomock.Controller
	recorder *MockToolLoggerMockRecorder
}

// MockToolLoggerMockRecorder is the mock recorder for MockToolLogger
type MockToolLoggerMockRecorder struct {
	mock *MockToolLogger
}

// NewMockToolLogger creates a new mock instance
func NewMockToolLogger(ctrl *gomock.Controller) *MockToolLogger {
	mock := &MockToolLogger{ctrl: ctrl}
	mock.recorder = &MockToolLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockToolLogger) EXPECT() *MockToolLoggerMockRecorder {
	return m.recorder
}

// WriteUser mocks base method
func (m *MockToolLogger) WriteUser(message string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteUser", message)
}

// WriteUser indicates an expected call of WriteUser
func (mr *MockToolLoggerMockRecorder) WriteUser(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteUser", reflect.TypeOf((*MockToolLogger)(nil).WriteUser), message)
}

// WriteDebug mocks base method
func (m *MockToolLogger) WriteDebug(message string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteDebug", message)
}

// WriteDebug indicates an expected call of WriteDebug
func (mr *MockToolLoggerMockRecorder) WriteDebug(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteDebug", reflect.TypeOf((*MockToolLogger)(nil).WriteDebug), message)
}

// WriteTrace mocks base method
func (m *MockToolLogger) WriteTrace(message string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteTrace", message)
}

// WriteTrace indicates an expected call of WriteTrace
func (mr *MockToolLoggerMockRecorder) WriteTrace(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteTrace", reflect.TypeOf((*MockToolLogger)(nil).WriteTrace), message)
}

// WriteMetric mocks base method
func (m *MockToolLogger) WriteMetric(metric *pb.OutputInfo) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteMetric", metric)
}

// WriteMetric indicates an expected call of WriteMetric
func (mr *MockToolLoggerMockRecorder) WriteMetric(metric interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteMetric", reflect.TypeOf((*MockToolLogger)(nil).WriteMetric), metric)
}

// ReadOutputInfo mocks base method
func (m *MockToolLogger) ReadOutputInfo() *pb.OutputInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadOutputInfo")
	ret0, _ := ret[0].(*pb.OutputInfo)
	return ret0
}

// ReadOutputInfo indicates an expected call of ReadOutputInfo
func (mr *MockToolLoggerMockRecorder) ReadOutputInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadOutputInfo", reflect.TypeOf((*MockToolLogger)(nil).ReadOutputInfo))
}
