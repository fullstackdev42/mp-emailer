package suite

import (
	"net/http/httptest"

	"github.com/fullstackdev42/mp-emailer/config"
	"github.com/fullstackdev42/mp-emailer/mocks"
	mocksCampaign "github.com/fullstackdev42/mp-emailer/mocks/campaign"
	mocksShared "github.com/fullstackdev42/mp-emailer/mocks/shared"
	mocksUser "github.com/fullstackdev42/mp-emailer/mocks/user"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type BaseTestSuite struct {
	suite.Suite
	Echo           *echo.Echo
	Context        echo.Context
	Recorder       *httptest.ResponseRecorder
	MockLogger     *mocks.MockLoggerInterface
	MockCampaign   *mocksCampaign.MockServiceInterface
	MockUser       *mocksUser.MockServiceInterface
	MockErrHandler *mocksShared.MockErrorHandlerInterface
	Config         *config.Config
}

func (s *BaseTestSuite) SetupTest() {
	s.Echo = echo.New()
	s.Recorder = httptest.NewRecorder()
	s.MockLogger = mocks.NewMockLoggerInterface(s.T())
	s.MockCampaign = mocksCampaign.NewMockServiceInterface(s.T())
	s.MockUser = mocksUser.NewMockServiceInterface(s.T())
	s.MockErrHandler = mocksShared.NewMockErrorHandlerInterface(s.T())
	s.Config = &config.Config{}
}

func (s *BaseTestSuite) TearDownTest() {
	s.MockLogger.AssertExpectations(s.T())
	s.MockCampaign.AssertExpectations(s.T())
	s.MockUser.AssertExpectations(s.T())
	s.MockErrHandler.AssertExpectations(s.T())
}

func (s *BaseTestSuite) NewContext(method, path string) echo.Context {
	req := httptest.NewRequest(method, path, nil)
	return s.Echo.NewContext(req, s.Recorder)
}
