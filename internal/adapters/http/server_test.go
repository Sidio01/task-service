package http_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	h "gitlab.com/g6834/team26/task/internal/adapters/http"
	"gitlab.com/g6834/team26/task/internal/domain/models"
	"gitlab.com/g6834/team26/task/internal/domain/task"
	"gitlab.com/g6834/team26/task/pkg/api"
	"gitlab.com/g6834/team26/task/pkg/logger"
	"gitlab.com/g6834/team26/task/pkg/mocks"
)

type authTestSuite struct {
	suite.Suite

	srv *h.Server
	db  *mocks.DbMock
	g   *mocks.GrpcMock
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, &authTestSuite{})
}

func (s *authTestSuite) SetupSuite() {
	// l := zerolog.Nop()
	l := logger.New()

	s.db = new(mocks.DbMock)
	s.g = new(mocks.GrpcMock)

	taskS := task.New(s.db, s.g)

	var err error
	s.srv, err = h.New(l, taskS)
	if err != nil {
		s.Suite.T().Errorf("db init failed: %s", err)
		s.Suite.T().FailNow()
	}

	go s.srv.Start()
}

func (s *authTestSuite) TearDownSuite() {
	_ = s.srv.Stop(context.Background())
}

func (s *authTestSuite) TestListHandlerSuccess() {
	s.db.On("List", "test123").Return([]*models.Task{&models.Task{UUID: "66f5b904-3f54-4da4-ba74-6dfdf8d72efb",
		Name:           "test",
		Text:           "this is test task",
		InitiatorLogin: "test123",
		Status:         "created",
		Approvals: []*models.Approval{&models.Approval{ApprovalLogin: "test626",
			N:        2,
			Sent:     sql.NullBool{Valid: true, Bool: false},
			Approved: sql.NullBool{Valid: false, Bool: false}}}}}, nil)
	s.g.On("Validate", mock.Anything, mock.Anything).Return(&api.AuthResponse{Result: true, Login: "test123", AccessToken: "AccessToken", RefreshToken: "RefreshToken"}, nil)

	bodyReq := strings.NewReader("")

	req, err := http.NewRequest("GET", "http://localhost:3000/task/v1/tasks/", bodyReq)
	s.NoError(err)

	client := http.Client{}
	response, err := client.Do(req)

	// log.Println(err)
	// data, err := ioutil.ReadAll(response.Body)
	// log.Println(string(data))
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)
	response.Body.Close()
}

func (s *authTestSuite) TestRunHandlerSuccess() {
	s.db.On("Run", mock.Anything).Return(nil)
	s.g.On("Validate", mock.Anything, mock.Anything).Return(&api.AuthResponse{Result: true, Login: "test123", AccessToken: "AccessToken", RefreshToken: "RefreshToken"}, nil)

	bodyReq := strings.NewReader("{\"approvalLogins\": [\"test626\",\"zxcvb\"],\"initiatorLogin\": \"test123\"}")

	req, err := http.NewRequest("POST", "http://localhost:3000/task/v1/tasks/run", bodyReq)
	s.NoError(err)

	client := http.Client{}
	response, err := client.Do(req)

	// log.Println(err)
	// data, err := ioutil.ReadAll(response.Body)
	// log.Println(string(data))
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)
	response.Body.Close()
}

func (s *authTestSuite) TestApproveHandlerSuccess() {
	s.db.On("Approve", "test123", "66f5b904-3f54-4da4-ba74-6dfdf8d72efb", "test626").Return(nil)
	s.g.On("Validate", mock.Anything, mock.Anything).Return(&api.AuthResponse{Result: true, Login: "test123", AccessToken: "AccessToken", RefreshToken: "RefreshToken"}, nil)

	approvalLogin := "test626"
	uuid := "66f5b904-3f54-4da4-ba74-6dfdf8d72efb"
	bodyReq := strings.NewReader("")

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:3000/task/v1/tasks/%s/approve/%s", uuid, approvalLogin), bodyReq)
	s.NoError(err)

	client := http.Client{}
	response, err := client.Do(req)

	// log.Println(err)
	// data, err := ioutil.ReadAll(response.Body)
	// log.Println(string(data))
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)
	response.Body.Close()
}

func (s *authTestSuite) TestDeclineHandlerSuccess() {
	s.db.On("Decline", "test123", "66f5b904-3f54-4da4-ba74-6dfdf8d72efb", "test626").Return(nil)
	s.g.On("Validate", mock.Anything, mock.Anything).Return(&api.AuthResponse{Result: true, Login: "test123", AccessToken: "AccessToken", RefreshToken: "RefreshToken"}, nil)

	approvalLogin := "test626"
	uuid := "66f5b904-3f54-4da4-ba74-6dfdf8d72efb"
	bodyReq := strings.NewReader("")

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:3000/task/v1/tasks/%s/decline/%s", uuid, approvalLogin), bodyReq)
	s.NoError(err)

	client := http.Client{}
	response, err := client.Do(req)

	// log.Println(err)
	// data, err := ioutil.ReadAll(response.Body)
	// log.Println(string(data))
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)
	response.Body.Close()
}

func (s *authTestSuite) TestDeleteHandlerSuccess() {
	s.db.On("Delete", "test123", "66f5b904-3f54-4da4-ba74-6dfdf8d72efb").Return(nil)
	s.g.On("Validate", mock.Anything, mock.Anything).Return(&api.AuthResponse{Result: true, Login: "test123", AccessToken: "AccessToken", RefreshToken: "RefreshToken"}, nil)

	uuid := "66f5b904-3f54-4da4-ba74-6dfdf8d72efb"
	bodyReq := strings.NewReader("")

	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:3000/task/v1/tasks/%s", uuid), bodyReq)
	s.NoError(err)

	client := http.Client{}
	response, err := client.Do(req)

	// log.Println(err)
	// data, err := ioutil.ReadAll(response.Body)
	// log.Println(string(data))
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)
	response.Body.Close()
}
