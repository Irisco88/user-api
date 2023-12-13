package mock_db

import (
	context "context"
	reflect "reflect"

	squirrel "github.com/Masterminds/squirrel"
	gomock "github.com/golang/mock/gomock"
	userv1 "github.com/irisco88/protos/gen/user/v1"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

type MockUserDBConn struct {
	ctrl     *gomock.Controller
	recorder *MockUserDBConnMockRecorder
}

type MockUserDBConnMockRecorder struct {
	mock *MockUserDBConn
}

func NewMockUserDBConn(ctrl *gomock.Controller) *MockUserDBConn {
	mock := &MockUserDBConn{ctrl: ctrl}
	mock.recorder = &MockUserDBConnMockRecorder{mock}
	return mock
}

func (m *MockUserDBConn) EXPECT() *MockUserDBConnMockRecorder {
	return m.recorder
}

func (m *MockUserDBConn) CreateUser(ctx context.Context, ownerID uint32, user *userv1.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, ownerID, user)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockUserDBConnMockRecorder) CreateUser(ctx, ownerID, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserDBConn)(nil).CreateUser), ctx, ownerID, user)
}

func (m *MockUserDBConn) DeleteUser(ctx context.Context, userID uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}
func (mr *MockUserDBConnMockRecorder) DeleteUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserDBConn)(nil).DeleteUser), ctx, userID)
}

func (m *MockUserDBConn) GetPgConn() *pgxpool.Pool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPgConn")
	ret0, _ := ret[0].(*pgxpool.Pool)
	return ret0
}

// GetPgConn indicates an expected call of GetPgConn.
func (mr *MockUserDBConnMockRecorder) GetPgConn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPgConn", reflect.TypeOf((*MockUserDBConn)(nil).GetPgConn))
}
func (m *MockUserDBConn) GetSQLBuilder() squirrel.StatementBuilderType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSQLBuilder")
	ret0, _ := ret[0].(squirrel.StatementBuilderType)
	return ret0
}
func (mr *MockUserDBConnMockRecorder) GetSQLBuilder() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSQLBuilder", reflect.TypeOf((*MockUserDBConn)(nil).GetSQLBuilder))
}

func (m *MockUserDBConn) GetUser(ctx context.Context, userID uint32) (*userv1.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, userID)
	ret0, _ := ret[0].(*userv1.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockUserDBConnMockRecorder) GetUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUserDBConn)(nil).GetUser), ctx, userID)
}

func (m *MockUserDBConn) GetUsers(ctx context.Context) ([]*userv1.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", ctx)
	ret0, _ := ret[0].([]*userv1.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockUserDBConnMockRecorder) GetUsers(ctx context.Context) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockUserDBConn)(nil).GetUsers), ctx)
}

// GetUserByEmailUserName mocks base method.
func (m *MockUserDBConn) GetUserByEmailUserName(ctx context.Context, userNameEmail string) (*userv1.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmailUserName", ctx, userNameEmail)
	ret0, _ := ret[0].(*userv1.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmailUserName indicates an expected call of GetUserByEmailUserName.
func (mr *MockUserDBConnMockRecorder) GetUserByEmailUserName(ctx, userNameEmail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmailUserName", reflect.TypeOf((*MockUserDBConn)(nil).GetUserByEmailUserName), ctx, userNameEmail)
}

// UpdateUser mocks base method.
func (m *MockUserDBConn) UpdateUser(ctx context.Context, user *userv1.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserDBConnMockRecorder) UpdateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserDBConn)(nil).UpdateUser), ctx, user)
}
