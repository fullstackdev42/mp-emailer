package testutil

import (
	"encoding/json"
	"io"
	"net/http/httptest"

	"github.com/gorilla/sessions"
	"github.com/jonesrussell/mp-emailer/config"
	"github.com/jonesrussell/mp-emailer/mocks"
	mocksCampaign "github.com/jonesrussell/mp-emailer/mocks/campaign"
	mocksEmail "github.com/jonesrussell/mp-emailer/mocks/email"
	mocksShared "github.com/jonesrussell/mp-emailer/mocks/shared"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type BaseTestSuite struct {
	suite.Suite
	Echo                        *echo.Echo
	Context                     echo.Context
	Recorder                    *httptest.ResponseRecorder
	Logger                      *mocks.MockLoggerInterface
	CampaignService             *mocksCampaign.MockServiceInterface
	RepresentativeLookupService *mocksCampaign.MockRepresentativeLookupServiceInterface
	EmailService                *mocksEmail.MockService
	CampaignClient              *mocksCampaign.MockClientInterface
	ErrorHandler                *mocksShared.MockErrorHandlerInterface
	TemplateRenderer            *mocksShared.MockTemplateRendererInterface
	Config                      *config.Config
	Store                       sessions.Store
}

func (s *BaseTestSuite) SetupTest() {
	s.Echo = echo.New()
	s.Recorder = httptest.NewRecorder()
	s.Logger = mocks.NewMockLoggerInterface(s.T())
	s.CampaignService = mocksCampaign.NewMockServiceInterface(s.T())
	s.RepresentativeLookupService = mocksCampaign.NewMockRepresentativeLookupServiceInterface(s.T())
	s.EmailService = mocksEmail.NewMockService(s.T())
	s.CampaignClient = mocksCampaign.NewMockClientInterface(s.T())
	s.ErrorHandler = mocksShared.NewMockErrorHandlerInterface(s.T())
	s.TemplateRenderer = mocksShared.NewMockTemplateRendererInterface(s.T())
	s.Config = &config.Config{}
	s.Store = sessions.NewCookieStore([]byte("test-secret"))
}

func (s *BaseTestSuite) TearDownTest() {
	s.Logger.AssertExpectations(s.T())
	s.CampaignService.AssertExpectations(s.T())
	s.RepresentativeLookupService.AssertExpectations(s.T())
	s.EmailService.AssertExpectations(s.T())
	s.CampaignClient.AssertExpectations(s.T())
	s.ErrorHandler.AssertExpectations(s.T())
	s.TemplateRenderer.AssertExpectations(s.T())
}

func (s *BaseTestSuite) NewContext(method, path string, body io.Reader) echo.Context {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return s.Echo.NewContext(req, s.Recorder)
}

// Helper methods for common assertions
func (s *BaseTestSuite) AssertJSONResponse(expectedCode int, expectedBody interface{}) {
	s.Equal(expectedCode, s.Recorder.Code)
	var actualBody interface{}
	s.NoError(json.NewDecoder(s.Recorder.Body).Decode(&actualBody))
	s.Equal(expectedBody, actualBody)
}

func (s *BaseTestSuite) AssertErrorResponse(expectedCode int, expectedMessage string) {
	s.Equal(expectedCode, s.Recorder.Code)
	var errResp struct {
		Message string `json:"message"`
	}
	s.NoError(json.NewDecoder(s.Recorder.Body).Decode(&errResp))
	s.Equal(expectedMessage, errResp.Message)
}

func (s *BaseTestSuite) SetSession(ctx echo.Context, name string, value interface{}) error {
	session, err := s.Store.New(ctx.Request(), "session-name")
	if err != nil {
		return err
	}
	session.Values[name] = value
	return session.Save(ctx.Request(), ctx.Response().Writer)
}

func (s *BaseTestSuite) RunMiddleware(middleware echo.MiddlewareFunc, ctx echo.Context) error {
	handler := middleware(func(_ echo.Context) error {
		return nil
	})
	return handler(ctx)
}
