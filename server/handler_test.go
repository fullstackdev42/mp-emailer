package server_test

import (
	"testing"

	"github.com/fullstackdev42/mp-emailer/internal/testutil"
	"github.com/fullstackdev42/mp-emailer/server"
	"github.com/fullstackdev42/mp-emailer/shared"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	testutil.BaseTestSuite
	handler server.HandlerInterface
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.BaseTestSuite.SetupTest()

	suite.handler = server.NewHandler(server.HandlerParams{
		BaseHandlerParams: shared.BaseHandlerParams{
			Store:            suite.Store,
			Logger:           suite.Logger,
			ErrorHandler:     suite.ErrorHandler,
			TemplateRenderer: suite.TemplateRenderer,
			Config:           suite.Config,
		},
		CampaignService: suite.CampaignService,
		EmailService:    suite.EmailService,
	})
}
